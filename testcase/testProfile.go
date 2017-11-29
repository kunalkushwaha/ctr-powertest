package testcase

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/kunalkushwaha/ctr-powertest/libruntime"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type ProfileContainerTest struct {
	Runtime libruntime.Runtime
}

type profileData struct {
	minCreate, minStart, minStop, minDel time.Duration
	avgCreate, avgStart, avgStop, avgDel time.Duration
	maxCreate, maxStart, maxStop, maxDel time.Duration
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

	data := new(profileData)

	for i := 0; i <= 50; i++ {

		createStartTime := time.Now()
		ctr, err := t.Runtime.Create(ctx, containerName+strconv.Itoa(i+5010), imageName, nil)
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

		waitForContainerEvent(statusC)

		stopStartTime := time.Now()
		err = t.Runtime.Stop(ctx, ctr)
		if err != nil {
			return errors.Wrap(err, "failed to stop container")
		}
		stopTotalTime := time.Now().Sub(stopStartTime)

		deleteStartTime := time.Now()
		err = t.Runtime.Delete(ctx, ctr)
		if err != nil {
			return errors.Wrap(err, "failed to delete container")
		}
		deleteTotalTime := time.Now().Sub(deleteStartTime)

		data = updateProfileData(data, createTotalTime, startTotalTime, stopTotalTime, deleteTotalTime)

	}
	fmt.Println("\nContainer Profile data\n")
	fmt.Printf("\tCreate\tStart\tStop\tDelete\n")
	fmt.Printf("Min\t%.2fs\t%.2fs\t%.2fs\t%.2fs\n", data.minCreate.Seconds(), data.minStart.Seconds(), data.minStop.Seconds(), data.minDel.Seconds())
	fmt.Printf("Avg\t%.2fs\t%.2fs\t%.2fs\t%.2fs\n", data.avgCreate.Seconds()/50, data.avgStart.Seconds()/50, data.avgStop.Seconds()/50, data.avgDel.Seconds()/50)
	fmt.Printf("Max\t%.2fs\t%.2fs\t%.2fs\t%.2fs\n", data.maxCreate.Seconds(), data.maxStart.Seconds(), data.maxStop.Seconds(), data.maxDel.Seconds())
	fmt.Printf("Total\t%.2fs\t%.2fs\t%.2fs\t%.2fs\n", data.avgCreate.Seconds(), data.avgStart.Seconds(), data.avgStop.Seconds(), data.avgDel.Seconds())
	return nil
}

func updateProfileData(base *profileData, create, start, stop, delete time.Duration) *profileData {
	if create < base.minCreate || base.minCreate == 0 {
		base.minCreate = create
	}
	if create > base.maxCreate || base.maxCreate == 0 {
		base.maxCreate = create
	}
	if start < base.minStart || base.minStart == 0 {
		base.minStart = start
	}
	if start > base.maxStart || base.maxStart == 0 {
		base.maxStart = start
	}
	if stop < base.minStop || base.minStop == 0 {
		base.minStop = stop
	}
	if stop > base.maxStop || base.maxStop == 0 {
		base.maxStop = stop
	}
	if delete < base.minDel || base.minDel == 0 {
		base.minDel = delete
	}
	if delete > base.maxDel || base.maxDel == 0 {
		base.maxDel = delete
	}

	base.avgCreate = base.avgCreate + create
	base.avgStart = base.avgStart + start
	base.avgStop = base.avgStop + stop
	base.avgDel = base.avgDel + delete

	return base
}
