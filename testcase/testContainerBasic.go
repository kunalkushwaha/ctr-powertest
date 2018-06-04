package testcase

import (
	"context"
	"fmt"
	"time"

	"github.com/kunalkushwaha/ctr-powertest/libocispec"
	"github.com/kunalkushwaha/ctr-powertest/libruntime"
	log "github.com/sirupsen/logrus"
)

type BasicContainerTest struct {
	Runtime libruntime.Runtime
}

func (t *BasicContainerTest) RunTestCases(ctx context.Context, testcases, args []string) error {

	log.Info("Running tests on ", t.Runtime.Version(ctx))
	if err := t.TestPullContainerImage(ctx, testImage); err != nil {
		return err
	}
	
	//time.Sleep(8 * time.Second)

	if err := t.TestCreateContainers(ctx, "test", testImage); err != nil {
		return err
	}
	if err := t.TestCreateRunningContainers(ctx, testContainerName, testImage); err != nil {
		return err
	}
	if err := t.TestCreateRunningNWaitContainers(ctx, testContainerName, testImage); err != nil {
		return err
	}
	
	return nil
}

func (t *BasicContainerTest) TestPullContainerImage(ctx context.Context, imageName string) error {
	log.Info("TestPullContainerImage..")
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

func (t *BasicContainerTest) TestCreateContainers(ctx context.Context, containerName, imageName string) error {
	//TODO :
	// Test with tty container,
	// Test without tty container
	// Test background container.
	log.Info("TestCreateContainers..")
	startTime := time.Now()
	ctr, err := t.Runtime.Create(ctx, containerName, imageName, nil)
	if err != nil {
		return err
	}

	err = t.Runtime.Delete(ctx, ctr)
	if err != nil {
		return err
	}
	totalTime := time.Now().Sub(startTime)
	log.Infof("%d containers in %s ", 1, totalTime.String())
	log.Info("OK..")
	return nil
}

func (t *BasicContainerTest) TestCreateRunningContainers(ctx context.Context, containerName, imageName string) error {
	log.Info("TestCreateRunningContainers..")
	//startTime := time.Now()
	statusC, ctr, err := t.Runtime.Run(ctx, containerName, imageName, nil)
	if err != nil {
		return err
	}

	StopstartTime := time.Now()
	err = t.Runtime.Stop(ctx, ctr)
	if err != nil {
		return fmt.Errorf("Container Stop: %v", err)
	}
	stopTotalTime := time.Now().Sub(StopstartTime)
	log.Infof("Stop time: %d containers in %s ", 1, stopTotalTime.String())

	waitForContainerEvent(statusC)

	deleteTime := time.Now()
	err = t.Runtime.Delete(ctx, ctr)
	if err != nil {
		return fmt.Errorf("Container Delete: %v", err)
	}
	totalTime := time.Now().Sub(deleteTime)
	log.Infof("%d containers in %s ", 1, totalTime.String())
	log.Info("OK..")
	return nil
}

func (t *BasicContainerTest) TestCreateRunningNWaitContainers(ctx context.Context, containerName, imageName string) error {
	log.Info("TestCreateRunningNWaitContainers..")
	specs, err := libocispec.GenerateSpec(libocispec.WithProcessArgs("sleep", "5s"))
	if err != nil {
		return err
	}
	// Create -> Runnable -> Wait -> Start -> Stop -> Listen -> Delete

	ctr, err := t.Runtime.Create(ctx, containerName, imageName, specs)
	if err != nil {
		return err
	}

	err = t.Runtime.Runnable(ctx, ctr)
	if err != nil {
		return err
	}

	statusC, err := t.Runtime.Wait(ctx, ctr)
	if err != nil {
		return err
	}

	err = t.Runtime.Start(ctx, ctr)
	if err != nil {
		return err
	}

	waitForContainerEvent(statusC)

	err = t.Runtime.Stop(ctx, ctr)
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
