package testcase

import (
	"context"
	"strconv"
	"sync"

	"github.com/kunalkushwaha/ctr-powertest/libruntime"
	log "github.com/sirupsen/logrus"
)

type ParallelClientSetup struct {
	Runtime libruntime.Runtime
}

func SetupParallelTestEnvironment(runtime string, config libruntime.RuntimeConfig, clean bool) (Testcases, error) {
	//Setup server and client.
	//Setup the number of threads and itterations to run.
	containerdRuntime, err := getRuntime(config)
	if err != nil {
		return singleClientTest{}, err
	}
	return ParallelClientSetup{containerdRuntime}, nil
}

func (t ParallelClientSetup) RunAllTests(ctx context.Context) error {
	//Create and Delete container in loop.
	if err := t.TestContainerCreateDelete(ctx, 12, 10); err != nil {
		return err
	}
	//Run and Stop containers in loop.
	return nil
}

func (t ParallelClientSetup) TestContainerCreateDelete(ctx context.Context, parallelCount, loopCount int) error {
	var wg sync.WaitGroup
	wg.Add(parallelCount)
	for i := 0; i < parallelCount; i++ {
		go t.createDeleteContainers(ctx, i, loopCount, &wg)
	}
	wg.Wait()
	return nil
}

func (t ParallelClientSetup) createDeleteContainers(ctx context.Context, id, loopCount int, wg *sync.WaitGroup) {
	defer wg.Done()
	//Genertate Specs
	//Seed random number
	//Image name
	for i := 0; i < loopCount; i++ {
		ctr, err := t.Runtime.Run(ctx, testContainerName+"-"+strconv.Itoa(id+100)+"-"+strconv.Itoa(i+100), testImage, nil)
		if err != nil {
			log.Debug("Container Run: %v", err)
			return
		}
		log.Info("Container ID : ", ctr.ID)
		err = t.Runtime.Stop(ctx, ctr)
		if err != nil {
			log.Debug("Container Stop: %v", err)
			return
		}
		err = t.Runtime.Delete(ctx, ctr)
		if err != nil {
			log.Debug("Container Delete: %v", err)
			return
		}
	}

}

//TestParallelPullImage Pulls image concurently.
func (t ParallelClientSetup) TestParallelPullImage(ctx context.Context, image string, wg *sync.WaitGroup) {
	//
}
