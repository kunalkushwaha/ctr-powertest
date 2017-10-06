package testcase

import (
	"context"
	"strconv"
	"sync"
	"time"

	"github.com/kunalkushwaha/ctr-powertest/libruntime"
	log "github.com/sirupsen/logrus"
)

type StressTest struct {
	Runtime libruntime.Runtime
}

func (t *StressTest) RunAllTests(ctx context.Context, args []string) error {
	log.Info("Running tests on ", t.Runtime.Version(ctx))
	//Create and Delete container in loop.
	if err := t.TestContainerCreateDelete(ctx, 4, 50); err != nil {
		return err
	}

	return nil
}

func (t *StressTest) TestContainerCreateDelete(ctx context.Context, parallelCount, loopCount int) error {
	//Make Sure image exist
	_, err := t.Runtime.Pull(ctx, testImage)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	wg.Add(parallelCount)
	log.Infof("Creating %d container in %d goroutines", loopCount, parallelCount)
	startTime := time.Now()
	for i := 0; i < parallelCount; i++ {
		go t.createDeleteContainers(ctx, i, loopCount, &wg)
	}
	wg.Wait()
	totalTime := time.Now().Sub(startTime)
	log.Infof("%d containers in %s ", loopCount*parallelCount, totalTime.String())

	log.Info("OK")
	return nil
}

func (t *StressTest) createDeleteContainers(ctx context.Context, id, loopCount int, wg *sync.WaitGroup) error {
	defer wg.Done()
	//Genertate Specs
	//Seed random number
	//Image name
	for i := 0; i < loopCount; i++ {

		ctr, err := t.Runtime.Create(ctx, testContainerName+"-"+strconv.Itoa(id+1020)+"-"+strconv.Itoa(i+1010), testImage, nil)
		//ctr, err := t.Runtime.Create(ctx, containerName, imageName, specs)
		if err != nil {
			log.Error(err)
			return err
		}

		err = t.Runtime.Runnable(ctx, ctr)
		if err != nil {
			log.Error(err)
			return err
		}

		statusC, err := t.Runtime.Wait(ctx, ctr)
		if err != nil {
			log.Error(err)
			return err
		}

		err = t.Runtime.Start(ctx, ctr)
		if err != nil {
			log.Error(err)
			return err
		}

		waitForContainerEvent(statusC)

		err = t.Runtime.Stop(ctx, ctr)
		if err != nil {
			log.Error(err)
			return err
		}

		err = t.Runtime.Delete(ctx, ctr)
		if err != nil {
			log.Error(err)
			return err
		}

	}
	return nil
}

/*
	TODO:
	- Parallel pull images
	- Pull &delete images at same time
	- Delete and Exec Containers
*/
