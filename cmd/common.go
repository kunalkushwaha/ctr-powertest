package cmd

import (
	"github.com/kunalkushwaha/ctr-powertest/libruntime"
	"github.com/kunalkushwaha/ctr-powertest/testcase"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var stdContainerdConfig = libruntime.RuntimeConfig{
	RuntimeName:      "containerd",
	RunDefaultServer: false,
	Root:             "/var/lib/containerd",
	State:            "/run/containerd",
	RuntimeEndpoint:  "/run/containerd/containerd.sock",
	DebugEndpoint:    "/run/containerd/debug.sock",
	DebugLevel:       "info",
}

var stdCRIOConfig = libruntime.RuntimeConfig{
	RuntimeName:      "crio",
	RunDefaultServer: false,
	RuntimeEndpoint:  "/var/run/crio.sock",
}

var stdCRIContainerdConfig = libruntime.RuntimeConfig{
	RuntimeName:      "cri-containerd",
	RunDefaultServer: false,
	RuntimeEndpoint:  "/var/run/cri-containerd.sock",
}

func initTestSuite(cmd *cobra.Command) {
	var err error

	debug, _ := cmd.Flags().GetBool("debug")
	if debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	proto, _ := cmd.Flags().GetString("proto")
	runtime, _ := cmd.Flags().GetString("runtime")
	switch proto {
	case "containerd":
		ctrRuntime, err = testcase.SetupTestEnvironment(proto, stdContainerdConfig, false)
		if err != nil {
			log.Fatal("Error while setting up environment : ", err)
		}
	case "cri":
		var config libruntime.RuntimeConfig
		if runtime == "crio" {
			config = stdCRIOConfig
		} else {
			config = stdCRIContainerdConfig
		}
		ctrRuntime, err = testcase.SetupTestEnvironment(proto, config, false)
		if err != nil {
			log.Fatal("Error while setting up environment : ", err)
		}
	}
}
