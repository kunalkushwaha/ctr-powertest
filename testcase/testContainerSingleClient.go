package testcase

import (
	"bytes"
	"context"
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/kunalkushwaha/ctr-powertest/libocispec"
	"github.com/kunalkushwaha/ctr-powertest/libruntime"
	"github.com/kunalkushwaha/ctr-powertest/libruntime/libcontainerd"
	log "github.com/sirupsen/logrus"
)

type singleClientTest struct {
	Runtime libruntime.Runtime
}

//SetupTestEnvironment to run test
func SetupTestEnvironment(runtime string, config libruntime.RuntimeConfig, clean bool) (Testcases, error) {

	//TODO:
	// if clean {
	//	cleanup(runtime root folder)
	//}

	// Get the runtime.
	containerdRuntime, err := getRuntime(config)
	if err != nil {
		log.Error("Error in getRUntime")
		return singleClientTest{}, err
	}
	return singleClientTest{containerdRuntime}, nil
}

//Test cases with 1 client.

//SetupTestEnvironment("containerd", config)

func (t singleClientTest) RunAllTests(ctx context.Context) error {
	//TODO: Instead of calling all tests directly,
	//		Build a machanism to execute, store results on file or DB for later processing.
	log.Info("Running tests on ", t.Runtime.Version(ctx))
	if err := t.TestPullContainerImage(ctx, testImage); err != nil {
		return err
	}
	if err := t.TestCreateContainers(ctx, testContainerName, testImage); err != nil {
		return err
	}

	if err := t.TestCreateRunningContainers(ctx, testContainerName, testImage); err != nil {
		return err
	}

	if err := t.TestContainerOutput(ctx, testContainerName, testImage); err != nil {
		return err
	}

	if err := t.TestExecContainers(ctx, testContainerName, testImage); err != nil {
		return err
	}

	log.Info("All Tests successfuly completed! ")
	return nil
}

func getRuntime(config libruntime.RuntimeConfig) (libruntime.Runtime, error) {
	//Get available runtime.
	if config.RuntimeName == "containerd" {
		return libcontainerd.GetNewContainerdRuntime(config, config.RunDefaultServer)
	}
	return nil, fmt.Errorf("Runtime not supported : %s ", config.RuntimeName)
}

//TestPullContainerImage tests Pull image API.
func (t singleClientTest) TestPullContainerImage(ctx context.Context, imageName string) error {
	log.Info("TestPullContainerImage..")
	//TODO:
	// Pull image from remote repo.
	_, err := t.Runtime.Pull(ctx, imageName)
	if err != nil {
		return err
	}

	// Pull image in already present in locally.
	_, err = t.Runtime.Pull(ctx, imageName)
	if err != nil {
		return err
	}

	//TODO: Cleanup function.
	log.Info("OK..")
	return nil

}

func (t singleClientTest) TestCreateContainers(ctx context.Context, containerName, imageName string) error {
	//TODO :
	// Test with tty container,
	// Test without tty container
	// Test background container.
	log.Info("TestCreateContainers..")
	ctr, err := t.Runtime.Create(ctx, containerName, imageName, nil)
	if err != nil {
		return err
	}
	err = t.Runtime.Delete(ctx, ctr)
	if err != nil {
		return err
	}
	log.Info("OK..")
	return nil
}

func (t singleClientTest) TestCreateRunningContainers(ctx context.Context, containerName, imageName string) error {
	log.Info("TestCreateRunningContainers..")
	spec, err := libocispec.GenerateSpec(libocispec.WithProcessArgs("echo", "Hello-world"))
	if err != nil {
		return err
	}
	ctr, err := t.Runtime.Run(ctx, containerName, imageName, spec)
	if err != nil {
		return err
	}
	log.Info("Container ID : ", ctr.ID)
	err = t.Runtime.Stop(ctx, ctr)
	if err != nil {
		return fmt.Errorf("Container Stop: %v", err)
	}
	err = t.Runtime.Delete(ctx, ctr)
	if err != nil {
		return fmt.Errorf("Container Delete: %v", err)
	}
	log.Info("OK..")
	return nil
}

func (t singleClientTest) TestContainerOutput(ctx context.Context, containerName, imageName string) error {
	log.Info("TestContainerOutput..")

	spec, err := libocispec.GenerateSpec(libocispec.WithProcessArgs("echo", "Hello-world"))
	if err != nil {
		return err
	}

	container, err := t.Runtime.Create(ctx, containerName, imageName, spec)
	if err != nil {
		return err
	}
	defer t.Runtime.Delete(ctx, container)

	stdout := bytes.NewBuffer(nil)
	err = t.Runtime.Runnable(ctx, container, libruntime.NewIO(bytes.NewBuffer(nil), stdout, bytes.NewBuffer(nil)))
	if err != nil {
		return err
	}
	defer t.Runtime.Stop(ctx, container)

	//statusC := make(chan uint32, 1)
	//	go func() {
	statusC, err := t.Runtime.Wait(ctx, container)
	if err != nil {
		logrus.Error(err)
		return err
	}

	//}()

	if err := t.Runtime.Start(ctx, container); err != nil {
		return err
	}

	res := <-statusC
	code, _, _ := res.Result()
	if code != 0 {
		return fmt.Errorf("expected status 0 but received %d", res)
	}
	if err := t.Runtime.Stop(ctx, container); err != nil {
		return err

	}

	actual := stdout.String()
	// echo adds a new line
	if actual != "Hello-world\n" {
		return fmt.Errorf("expected output %q but received %q", "Hello-world\n", actual)
	}
	log.Info("Ok..")
	return nil
}

func (t singleClientTest) TestExecContainers(ctx context.Context, containerName, imageName string) error {
	log.Info("TestExecContainers..")
	//Create & Start a Task.
	spec, err := libocispec.GenerateSpec(libocispec.WithProcessArgs("sleep", "200"))
	if err != nil {
		return err
	}

	container, err := t.Runtime.Run(ctx, containerName, imageName, spec)
	if err != nil {
		return err
	}

	defer t.Runtime.Delete(ctx, container)

	finished := make(chan struct{}, 1)
	go func() {
		if _, err := t.Runtime.Wait(ctx, container); err != nil {
			return
		}
		close(finished)
	}()

	//Prepare another process, which can get exec'ed into container.
	err = t.Runtime.Exec(ctx, container, []string{"ps", "-ax"})
	if err != nil {
		t.Runtime.Stop(ctx, container)
		return err
	}
	//Kill the process.
	t.Runtime.Stop(ctx, container)
	<-finished
	log.Info("Ok..")
	return nil
}
