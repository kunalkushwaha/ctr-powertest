package testcase

import (
	"context"
	"fmt"
	"time"

	"github.com/kunalkushwaha/ctr-powertest/libruntime"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type ProfileContainerTest struct {
	Runtime libruntime.Runtime
}

func (t *ProfileContainerTest) RunTestCases(ctx context.Context, testcases, args []string) error {

	log.Info("Running tests on ", t.Runtime.Version(ctx))
	if err := t.TestPullContainerImage(ctx, testImage); err != nil {
		return err
	}
	if err := t.TestCreateRunningContainers(ctx, testContainerName, testImage); err != nil {
		return err
	}

	return nil
}

func (t *ProfileContainerTest) TestPullContainerImage(ctx context.Context, imageName string) error {
	// Pull image from remote repo.
	_, err := t.Runtime.Pull(ctx, imageName)
	if err != nil {
		return errors.Wrap(err, "failed to pull")
	}
	return nil
}

func (t *ProfileContainerTest) TestCreateRunningContainers(ctx context.Context, containerName, imageName string) error {

	createStartTime := time.Now()
	ctr, err := t.Runtime.Create(ctx, containerName, imageName, nil)
	if err != nil {
		return errors.Wrap(err, "failed to create")
	}
	createTotalTime := time.Now().Sub(createStartTime)

	err = t.Runtime.Runnable(ctx, ctr)
	if err != nil {
		return errors.Wrap(err, "failed to check runnable")
	}

	statusC, err := t.Runtime.Wait(ctx, ctr)
	if err != nil {
		return errors.Wrap(err, "failed to get wait")
	}

	startStartTime := time.Now()
	err = t.Runtime.Start(ctx, ctr)
	if err != nil {
		return errors.Wrap(err, "failed to start container")
	}
	startTotalTime := time.Now().Sub(startStartTime)

	stopStartTime := time.Now()
	err = t.Runtime.Stop(ctx, ctr)
	if err != nil {
		return errors.Wrap(err, "failed to stop container")
	}
	stopTotalTime := time.Now().Sub(stopStartTime)

	waitForContainerEvent(statusC)

	deleteStartTime := time.Now()
	err = t.Runtime.Delete(ctx, ctr)
	if err != nil {
		return errors.Wrap(err, "failed to delete container")
	}
	deleteTotalTime := time.Now().Sub(deleteStartTime)

	fmt.Println("\nContainer Profile data\n")
	fmt.Printf("Create\tStart\tStop\tDelete\n")
	fmt.Printf("%.2fs\t%.2fs\t%.2fs\t%.2fs\n", createTotalTime.Seconds(), startTotalTime.Seconds(), stopTotalTime.Seconds(), deleteTotalTime.Seconds())
	return nil
}
