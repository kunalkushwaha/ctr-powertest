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

// GetClient return containerd client
func GetClient(address string, namespace string) (*containerd.Client, error) {
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
func (c *ContainerdRuntime) Pull(ctx context.Context, image string) (libruntime.Image, error) {
	img, err := c.cclient.GetImage(ctx, image)
	if err != nil {
		img, err = c.cclient.Pull(ctx, image, containerd.WithPullUnpack)
		if err != nil {
			return libruntime.Image{}, err
		}
	} else {
		log.Debugf("Image \" %s \", already present on system ", image)
	}
	return libruntime.Image{Name: img.Name()}, nil
}

//Create creates the container from given image.
func (c *ContainerdRuntime) Create(ctx context.Context, containerName, imageName string, specs *runtimespecs.Spec) (libruntime.Container, error) {
	var image containerd.Image
	image, err := c.cclient.GetImage(ctx, imageName)
	if err != nil {
		return libruntime.Container{}, err
	}

	//If specs not provided, generate default specs.
	if specs == nil {
		specs, err = containerd.GenerateSpec(ctx, nil, nil)
		if err != nil {
			return libruntime.Container{}, err
		}
	}

	//Create new container.
	//ctr, err := c.cclient.NewContainer(ctx, containerName, containerd.WithSpec(specs), containerd.WithNewRootFS(containerName, image))
	ctr, err := c.cclient.NewContainer(ctx, containerName, containerd.WithSpec(specs), containerd.WithNewSnapshot(containerName, image))
	if err != nil {
		return libruntime.Container{}, err
	}
	return libruntime.Container{ID: ctr.ID()}, nil

}

//Run creates and run the task (container instance)
func (c *ContainerdRuntime) Run(ctx context.Context, containerName, imageName string, specs *runtimespecs.Spec) (libruntime.Container, error) {
	var container containerd.Container
	container, err := c.cclient.LoadContainer(ctx, containerName)
	if err != nil {
		log.Debugf("Not found container %s, so creating ", containerName)
		_, err := c.Create(ctx, containerName, imageName, specs)
		if err != nil {
			return libruntime.Container{}, fmt.Errorf("Create : %v ", err)
		}
		container, err = c.cclient.LoadContainer(ctx, containerName)
		if err != nil {
			return libruntime.Container{}, fmt.Errorf("Load Container : %v ", err)
		}
	}

	//Create task.
	log.Debug("Creating Task now.")
	stdout := bytes.NewBuffer(nil)
	task, err := container.NewTask(ctx, containerd.NewIO(bytes.NewBuffer(nil), stdout, bytes.NewBuffer(nil)))
	//task, err := container.NewTask(ctx, containerd.Stdio)
	if err != nil {
		return libruntime.Container{}, fmt.Errorf("Run : %v ", err)
	}
	status, err := task.Status(ctx)
	log.Debug("Starting Task now.", status, err)
	if err := task.Start(ctx); err != nil {
		log.Debug("Failed in Starting Task now.")
		//task.Delete(ctx)
		return libruntime.Container{}, fmt.Errorf("Task Start : %v ", err)
	}
	log.Debug("Returning after container Run()")
	return libruntime.Container{ID: container.ID()}, nil

}

func (c *ContainerdRuntime) Runnable(ctx context.Context, ctr libruntime.Container, ioCreate containerd.IOCreation) error {

	container, err := c.cclient.LoadContainer(ctx, ctr.ID)
	if err != nil {
		return err
	}

	_, err = container.Task(ctx, nil)
	if err == nil {
		return fmt.Errorf("Container already Runnable")
	}
	log.Debug("Runnable :  No TasK")
	//ioCtr := containerd.IOCreation(ioCreate)
	//Create new task.
	_, err = container.NewTask(ctx, ioCreate)
	if err != nil {
		return err
	}
	log.Debug("Runnable :  New task Created")

	return nil
}

func (c *ContainerdRuntime) Start(ctx context.Context, ctr libruntime.Container) error {
	container, err := c.cclient.LoadContainer(ctx, ctr.ID)
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
			//	task.Kill(ctx, syscall.SIGKILL)
			task.Delete(ctx)
			return fmt.Errorf("Task Start : %v ", err)
		}
	}
	return err
}

// Stop the task running instance of container
func (c *ContainerdRuntime) Stop(ctx context.Context, ctr libruntime.Container) error {

	container, err := c.cclient.LoadContainer(ctx, ctr.ID)
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
		log.Infof("container.%s function invoked for  %d", status, task.Pid())
		statusC, err := task.Wait(ctx)
		if err != nil {
			log.Errorf("container %q: error during wait: %v", task.Pid(), err)
		}

		err = task.Kill(ctx, syscall.SIGKILL)
		if err != nil {
			task.Delete(ctx)
			return err
		}
		log.Debug("Waiting for event")
		status := <-statusC
		code, _, err := status.Result()
		if err != nil {
			return err
		}
		if code != 0 {
			log.Debugf("exited container process: code: %d", status)
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
	log.Debug("Container Stop was called with status : ", status)
	return nil
}

func (c *ContainerdRuntime) Wait(ctx context.Context, ctr libruntime.Container) (<-chan containerd.ExitStatus, error) {
	container, err := c.cclient.LoadContainer(ctx, ctr.ID)
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

	return signal, nil
}

//Delete  the container from filesystem
func (c *ContainerdRuntime) Delete(ctx context.Context, ctr libruntime.Container) error {
	log.Debug("Delete container is invoked.. ", ctr.ID)
	container, err := c.cclient.LoadContainer(ctx, ctr.ID)
	if err != nil {
		return err
	}
	//	err = container.Delete(ctx, containerd.WithRootFSDeletion)
	err = container.Delete(ctx, containerd.WithSnapshotCleanup)
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

	//processStatusC := make(chan containerd.ExitStatus)
	//go func() {
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
func (c *ContainerdRuntime) GetContainer(ctx context.Context, containerName string) (libruntime.Container, error) {
	ctr, err := c.cclient.LoadContainer(ctx, containerName)
	if err != nil {
		return libruntime.Container{}, err
	}
	return libruntime.Container{ID: ctr.ID()}, nil
}

//Version of containerd
func (c *ContainerdRuntime) Version(ctx context.Context) string {
	version, _ := c.cclient.Version(ctx)
	return "containerd " + version.Version
}

func empty() containerd.IOCreation {
	null := ioutil.Discard
	return containerd.NewIO(bytes.NewBuffer(nil), null, null)
}
