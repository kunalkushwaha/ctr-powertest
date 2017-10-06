package cmd

import (
	"github.com/kunalkushwaha/ctr-powertest/libruntime"
	"github.com/kunalkushwaha/ctr-powertest/testcase"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var containerdConfig = libruntime.RuntimeConfig{
	RuntimeName:      "containerd",
	RunDefaultServer: true,
	Root:             "/var/lib/powertest",
	RuntimeEndpoint:  "/run/powertest/containerd.sock",
	DebugEndpoint:    "/run/powertest/debug.sock",
}

var stdContainerdConfig = libruntime.RuntimeConfig{
	RuntimeName:      "containerd",
	RunDefaultServer: false,
	Root:             "/var/lib/containerd",
	RuntimeEndpoint:  "/run/containerd/containerd.sock",
	DebugEndpoint:    "/run/containerd/debug.sock",
}

var stdCRIOConfig = libruntime.RuntimeConfig{
	RuntimeName:      "crio",
	RunDefaultServer: false,
	RuntimeEndpoint:  "/var/run/crio.sock",
}

func initTestSuite(cmd *cobra.Command) {
	var err error

	debug, _ := cmd.Flags().GetBool("debug")
	if debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	runtime, _ := cmd.Flags().GetString("runtime")
	switch runtime {
	case "containerd":
		ctrRuntime, err = testcase.SetupTestEnvironment(runtime, stdContainerdConfig, false)
		if err != nil {
			log.Fatal("Error while setting up environment : ", err)
		}
	case "crio":
		ctrRuntime, err = testcase.SetupTestEnvironment(runtime, stdCRIOConfig, false)
		if err != nil {
			log.Fatal("Error while setting up environment : ", err)
		}
	}
}
