package testcase

import (
	"context"
	"fmt"

	"github.com/kunalkushwaha/ctr-powertest/libruntime/libcri"

	"github.com/containerd/containerd"
	"github.com/kunalkushwaha/ctr-powertest/libruntime"
	"github.com/kunalkushwaha/ctr-powertest/libruntime/libcontainerd"
	log "github.com/sirupsen/logrus"
)

const (
	testImage         = "docker.io/library/alpine:latest"
	testContainerName = "test-powertest"
)

// Testcases interface to implement testcases
type Testcases interface {
	RunTestCases(context context.Context, testcases, args []string) error
}

//SetupTestEnvironment setups server and client for container runtime
func SetupTestEnvironment(ctx context.Context, proto string, config libruntime.RuntimeConfig, clean bool) (libruntime.Runtime, error) {

	//TODO:
	// if clean {
	//	cleanup(runtime root folder)
	//}

	// Get the runtime.
	ctrRuntime, err := getRuntime(ctx, proto, config)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return ctrRuntime, nil
}

func getRuntime(ctx context.Context, proto string, config libruntime.RuntimeConfig) (libruntime.Runtime, error) {
	//Get available runtime.
	switch proto {
	case "containerd":
		return libcontainerd.GetNewContainerdRuntime(ctx, config, config.RunDefaultServer)
	case "cri":
		return libcri.GetNewCRIRuntime(config, false)

	}
	return nil, fmt.Errorf("Proto not supported : %s ", proto)
}

func waitForContainerEvent(statusC <-chan interface{}) error {
	if statusC == nil {
		return nil
	}
	status := <-statusC
	switch p := status.(type) {
	case containerd.ExitStatus:
		if p.ExitCode() != 0 {
			log.Info(p.Result())
			err := p.Error()
			return err
		}
	default:
		return fmt.Errorf("Unknow Event")
	}
	return nil
}
