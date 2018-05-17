package cmd

import (
	"context"

	"github.com/containerd/containerd/namespaces"
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
	RuntimeEndpoint:  "/run/crio.sock",
}

var stdCRIContainerdConfig = libruntime.RuntimeConfig{
	RuntimeName:      "cri-containerd",
	RunDefaultServer: false,
	RuntimeEndpoint:  "/run/containerd/containerd.sock",
}

var stdCRIDockershimConfig = libruntime.RuntimeConfig{
	RuntimeName:      "dockershim",
	RunDefaultServer: false,
	RuntimeEndpoint:  "/run/dockershim.sock",
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
	ctx = context.Background()
	switch proto {
	case "containerd":
		ctx = namespaces.WithNamespace(ctx, "powertest")
		ctrRuntime, err = testcase.SetupTestEnvironment(ctx, proto, stdContainerdConfig, false)
		if err != nil {
			log.Fatal("Error while setting up environment : ", err)
		}
	case "cri":
		var config libruntime.RuntimeConfig
		if runtime == "crio" {
			config = stdCRIOConfig
		} else if runtime == "dockershim" {
			config = stdCRIDockershimConfig
		} else {
			config = stdCRIContainerdConfig
		}

		ctrRuntime, err = testcase.SetupTestEnvironment(ctx, proto, config, false)
		if err != nil {
			log.Fatal("Error while setting up environment : ", err)
		}
	}
}
