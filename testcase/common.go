package testcase

import (
	"context"
	"fmt"

	"github.com/kunalkushwaha/ctr-powertest/libruntime/libcrio"

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
	RunAllTests(context context.Context, args []string) error
}

//SetupTestEnvironment setups server and client for container runtime
func SetupTestEnvironment(runtime string, config libruntime.RuntimeConfig, clean bool) (libruntime.Runtime, error) {

	//TODO:
	// if clean {
	//	cleanup(runtime root folder)
	//}

	// Get the runtime.
	ctrRuntime, err := getRuntime(config)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return ctrRuntime, nil
}

func getRuntime(config libruntime.RuntimeConfig) (libruntime.Runtime, error) {
	//Get available runtime.
	switch config.RuntimeName {
	case "containerd":
		return libcontainerd.GetNewContainerdRuntime(config, config.RunDefaultServer)
	case "crio":
		return libcrio.GetNewCRIORuntime(config, false)

	}
	return nil, fmt.Errorf("Runtime not supported : %s ", config.RuntimeName)
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
