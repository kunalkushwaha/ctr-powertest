package libcontainerd

import (
	"github.com/containerd/containerd"
	"github.com/containerd/containerd/server"
	"github.com/kunalkushwaha/ctr-powertest/libruntime"
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
	//TODO: build containerd-config and have opten to start server too
	//localConfig := runtime2containerd(config)

	client, err = GetNewClient(config.RuntimeEndpoint, "powertest")
	if err != nil {
		return nil, err
	}
	return &ContainerdRuntime{serverInstance, client}, nil
}

func runtime2containerd(config libruntime.RuntimeConfig) server.Config {
	return server.Config{}
}
