package libcontainerd

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"math/rand"
	"syscall"

	"github.com/containerd/containerd"
	"github.com/kunalkushwaha/ctr-powertest/libruntime"
	runtimespecs "github.com/opencontainers/runtime-spec/specs-go"
	log "github.com/sirupsen/logrus"
)

func GetNewClient(address string, namespace string) (*containerd.Client, error) {
	if namespace == "" {
		return nil, fmt.Errorf("Namespace is required")
	}
	client, err := containerd.New(address, containerd.WithDefaultNamespace(namespace))
	if err != nil {
		return nil, err
	}
	return client, nil
}

//Pull the image from remote registry.
func (cr *ContainerdRuntime) Pull(ctx context.Context, image string) (libruntime.Image, error) {
	img, err := cr.cclient.GetImage(ctx, image)
	if err != nil {
		img, err = cr.cclient.Pull(ctx, image, containerd.WithPullUnpack)
		if err != nil {
			return libruntime.Image{}, err
		}
	} else {
		log.Debugf("Image \" %s \", already present on system ", image)
	}
	return libruntime.Image{Name: img.Name()}, nil
}

//RemoveImage from remote registry.
func (cr *ContainerdRuntime) RemoveImage(ctx context.Context, image string) error {

	err := cr.cclient.ImageService().Delete(ctx, image)
	if err != nil {
		return err
	}
	return nil
}

//Run creates and run the task (container instance)
func (cr *ContainerdRuntime) Run(ctx context.Context, containerName, imageName string, specs *runtimespecs.Spec) (<-chan interface{}, *libruntime.Container, error) {

	image, err := cr.cclient.GetImage(ctx, imageName)
	if err != nil {
		return nil, nil, err
	}
	if specs == nil {
		specs, err = containerd.GenerateSpec(ctx, cr.cclient, nil, containerd.WithImageConfig(image), containerd.WithProcessArgs("true"))
		if err != nil {
			return nil, nil, err
		}
	}

	//Create new container.
	ctr, err := cr.Create(ctx, containerName, imageName, specs)
	if err != nil {
		return nil, nil, err
	}

	err = cr.Runnable(ctx, ctr)
	if err != nil {
		return nil, nil, err
	}

	statusC, err := cr.Wait(ctx, ctr)
	if err != nil {
		return nil, nil, err
	}

	if err := cr.Start(ctx, ctr); err != nil {
		return statusC, nil, err
	}

	return statusC, ctr, nil
}

//Delete  the container from filesystem
func (cr *ContainerdRuntime) Delete(ctx context.Context, ctr *libruntime.Container) error {

	container, err := cr.cclient.LoadContainer(ctx, ctr.ID)
	if err != nil {
		return err
	}

	err = container.Delete(ctx, containerd.WithSnapshotCleanup)
	if err != nil {
		return err
	}
	return nil
}

func (cr *ContainerdRuntime) Wait(ctx context.Context, ctr *libruntime.Container) (<-chan interface{}, error) {
	container, err := cr.cclient.LoadContainer(ctx, ctr.ID)
	if err != nil {
		return nil, err
	}

	task, err := container.Task(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("Task Error : %v", err)
	}

	signal, err := task.Wait(ctx)
	if err != nil {
		return nil, err
	}

	returnSignal := make(chan interface{})
	go func() {
		//FIXME: Can this be better?
		tSignal := <-signal
		returnSignal <- tSignal
	}()
	return returnSignal, nil
}

func (cr *ContainerdRuntime) Start(ctx context.Context, ctr *libruntime.Container) error {
	container, err := cr.cclient.LoadContainer(ctx, ctr.ID)
	if err != nil {
		return err
	}

	task, err := container.Task(ctx, nil)
	if err == nil {
		status, _ := task.Status(ctx)
		if status.Status == containerd.Running {
			return fmt.Errorf("Container already running")
		}
		// Start the task
		if err := task.Start(ctx); err != nil {
			task.Delete(ctx)
			return fmt.Errorf("Task Start : %v ", err)
		}
	}
	return err
}

// Stop the task running instance of container
func (cr *ContainerdRuntime) Stop(ctx context.Context, ctr *libruntime.Container) error {

	container, err := cr.cclient.LoadContainer(ctx, ctr.ID)
	if err != nil {
		return err
	}

	task, err := container.Task(ctx, nil)
	if err != nil {
		return err
	}

	status, _ := task.Status(ctx)
	switch status.Status {

	case containerd.Stopped:

		_, err := task.Delete(ctx)
		if err != nil {
			return err
		}

	case containerd.Running, containerd.Created: // Created too creates a shim process and need to delete by killing process

		statusC, err := task.Wait(ctx)
		if err != nil {
			log.Errorf("container %s: error during wait: %v", container.ID(), err)
		}
		go func() {
			if err := task.Kill(ctx, syscall.SIGKILL); err != nil {
				task.Delete(ctx)
				return
			}
		}()

		status := <-statusC
		code, _, err := status.Result()
		if err != nil {
			log.Errorf("container %q: error getting task result code: %v", container.ID(), err)
		}
		if code != 0 {
			log.Debugf("%s: exited container process: code: %d", container.ID(), status)
		}
		_, err = task.Delete(ctx)
		if err != nil {
			return err
		}

	case containerd.Paused:
		return fmt.Errorf("Can't stop a paused container; unpause first")

	default:
		return fmt.Errorf("Undefined Task state %v ", status)
	}

	return nil
}

//Create creates the container from given image.
func (cr *ContainerdRuntime) Create(ctx context.Context, containerName, imageName string, specs *runtimespecs.Spec) (*libruntime.Container, error) {
	image, err := cr.cclient.GetImage(ctx, imageName)
	if err != nil {
		return nil, err
	}
	if specs == nil {
		specs, err = containerd.GenerateSpec(ctx, cr.cclient, nil, containerd.WithImageConfig(image), containerd.WithProcessArgs("true"))
		if err != nil {
			return nil, err
		}
	}

	//Create new container.
	_, err = cr.cclient.NewContainer(ctx, containerName, containerd.WithSpec(specs), containerd.WithNewSnapshot(containerName, image))
	if err != nil {
		return nil, fmt.Errorf("Error in Container Creation : %v", err)
	}

	return &libruntime.Container{ID: containerName}, nil
}

func (cr *ContainerdRuntime) Runnable(ctx context.Context, ctr *libruntime.Container) error {
	//FIXME: There will be testcases, which needs to be tested for output and tty.
	// INterface should provide mechanism for same.
	ioCreate := containerd.NullIO

	container, err := cr.cclient.LoadContainer(ctx, ctr.ID)
	if err != nil {
		return err
	}

	_, err = container.Task(ctx, nil)
	if err == nil {
		return fmt.Errorf("Container already Runnable")
	}

	//Create new task.
	_, err = container.NewTask(ctx, ioCreate)
	if err != nil {
		return err
	}
	return nil
}

func (c ContainerdRuntime) Exec(ctx context.Context, ctr libruntime.Container, cmd []string) error {
	// Get the task from container.
	container, err := c.cclient.LoadContainer(ctx, ctr.ID)
	if err != nil {
		return err
	}

	task, err := container.Task(ctx, nil)
	if err != nil {
		return fmt.Errorf("Task Error : %v", err)
	}

	//Extract Process spec from Specs
	specs, _ := container.Spec()
	processSpec := specs.Process
	// Prepare process
	execID := "powertest-" + fmt.Sprintf("%X", rand.Int())
	processSpec.Args = cmd
	// execute as exec.
	process, err := task.Exec(ctx, execID, processSpec, empty())
	if err != nil {
		return err

	}

	processStatusC, err := process.Wait(ctx)
	if err != nil {
		return err
	}

	if err := process.Start(ctx); err != nil {
		return err

	}
	// wait for the exec to return
	status := <-processStatusC
	code, _, err := status.Result()
	if err != nil {

		return err
	}
	if code != 6 {
		log.Errorf("expected exec exit code 6 but received %d", code)
	}
	_, err = process.Delete(ctx)
	if err != nil {
		return err

	}
	// store it or return PID.
	return nil
}

//GetContainer return the container, if it exists.
func (cr *ContainerdRuntime) GetContainer(ctx context.Context, containerName string) (*libruntime.Container, error) {
	ctr, err := cr.cclient.LoadContainer(ctx, containerName)
	if err != nil {
		return nil, err
	}
	return &libruntime.Container{ID: ctr.ID()}, nil
}

//Version of containerd
func (cr *ContainerdRuntime) Version(ctx context.Context) string {
	version, _ := cr.cclient.Version(ctx)
	return "containerd " + version.Version
}

func empty() containerd.IOCreation {
	null := ioutil.Discard
	return containerd.NewIO(bytes.NewBuffer(nil), null, null)
}
