package libcontainerd

import (
	"log"

	"github.com/containerd/containerd"
	"github.com/containerd/containerd/server"
	"github.com/kunalkushwaha/ctr-powertest/libruntime"
)

const (
	defaultServerGRPCAddress = "/run/containerd/containerd.sock"
	defaultRoot              = "/var/lib/powertest"
	testImage                = "docker.io/library/alpine:latest"
)

//ContainerdRuntime implements all containerd funtions
type ContainerdRuntime struct {
	cserver *server.Server
	cclient *containerd.Client
}

//GetNewContainerdRuntime creates new instance of containerd test setup
func GetNewContainerdRuntime(config libruntime.RuntimeConfig, startServer bool) (libruntime.Runtime, error) {
	var (
		serverInstance *server.Server
		client         *containerd.Client
		err            error
	)
	//localConfig := runtime2containerd(config)
	localConfig := server.Config{
		Root:  "/var/lib/powertest",
		State: "/run/powertest",
		GRPC: server.GRPCConfig{
			Address: "/run/powertest/containerd.sock",
		},
		Debug: server.Debug{
			Level:   "info",
			Address: "/run/powertest/debug.sock",
		},
	}
	if startServer {

		serverInstance, err = SetupNewServer(localConfig)
		if err != nil {
			log.Fatal("Unable setup server!!", err)
			return nil, err
		}

		client, err = GetNewClient(localConfig.GRPC.Address, "powertest")
		if err != nil {
			return nil, err
		}
	} else {
		client, err = GetNewClient(defaultServerGRPCAddress, "powertest")
		if err != nil {
			return nil, err
		}
	}
	//log := logrus.New()
	return &ContainerdRuntime{serverInstance, client}, nil
}

func runtime2containerd(config libruntime.RuntimeConfig) server.Config {
	return server.Config{}
}
