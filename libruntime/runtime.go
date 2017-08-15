package libruntime

import (
	"context"

	"github.com/containerd/containerd"
	"github.com/opencontainers/runtime-spec/specs-go"
)

type RuntimeConfig struct {
	RuntimeName      string
	RunDefaultServer bool
	Root             string
	RuntimeEndpoint  string
	DebugEndpoint    string
}

type Container struct {
	ID string
}

type Image struct {
	Name string
}

type Runtime interface {
	//	NewServer(RuntimeConfig)
	//	GetClient(string, string) (Runtime, error)
	Version(context.Context) string
	Pull(context.Context, string) (Image, error)
	Create(context.Context, string, string, *specs.Spec) (Container, error)
	Run(context.Context, string, string, *specs.Spec) (Container, error)
	Stop(context.Context, Container) error
	Delete(context.Context, Container) error
	Runnable(context.Context, Container, containerd.IOCreation) error
	Start(context.Context, Container) error
	Exec(context.Context, Container, []string) error
	Wait(context.Context, Container) (uint32, error)
	GetContainer(context.Context, string) (Container, error)
}
