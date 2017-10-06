package libruntime

import (
	"context"

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
	ID    string
	PodID string
}

type Image struct {
	Name string
}

type Runtime interface {
	//	NewServer(RuntimeConfig)
	//	GetClient(string, string) (Runtime, error)
	Version(context.Context) string
	Pull(context.Context, string) (Image, error)
	Create(context context.Context, containerName string, imageName string, OCISpecs *specs.Spec) (*Container, error)
	Run(context.Context, string, string, *specs.Spec) (<-chan interface{}, *Container, error)
	Stop(context.Context, *Container) error
	Delete(context.Context, *Container) error
	Runnable(context.Context, *Container) error
	Start(context.Context, *Container) error
	Exec(context.Context, Container, []string) error
	Wait(context.Context, *Container) (<-chan interface{}, error)
	GetContainer(context.Context, string) (*Container, error)
}
