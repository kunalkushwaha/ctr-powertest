package libcri

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/kunalkushwaha/ctr-powertest/libruntime"
	runtimespecs "github.com/opencontainers/runtime-spec/specs-go"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	pb "k8s.io/kubernetes/pkg/kubelet/apis/cri/v1alpha1/runtime"
)

func GetNewRuntimeClient(socket string, timeout time.Duration) (*pb.RuntimeServiceClient, error) {
	conn, err := grpc.Dial(socket, grpc.WithInsecure(), grpc.WithTimeout(timeout),
		grpc.WithDialer(func(addr string, timeout time.Duration) (net.Conn, error) {
			return net.DialTimeout("unix", addr, timeout)
		}))
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %v", err)
	}

	client := pb.NewRuntimeServiceClient(conn)
	return &client, nil
}

func GetNewImageClient(socket string, timeout time.Duration) (*pb.ImageServiceClient, error) {
	conn, err := grpc.Dial(socket, grpc.WithInsecure(), grpc.WithTimeout(timeout),
		grpc.WithDialer(func(addr string, timeout time.Duration) (net.Conn, error) {
			return net.DialTimeout("unix", addr, timeout)
		}))
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %v", err)
	}

	client := pb.NewImageServiceClient(conn)
	return &client, nil
}

func (cr *CRIRuntime) Version(ctx context.Context) string {

	r, err := (*cr.RuntimeClient).Version(ctx, &pb.VersionRequest{})
	if err != nil {
		return err.Error()
	}
	log.Debug(r.String())
	return r.String()
}

func (cr *CRIRuntime) Pull(ctx context.Context, imageName string) (libruntime.Image, error) {

	img, err := (*cr.ImageClient).PullImage(ctx, &pb.PullImageRequest{Image: &pb.ImageSpec{Image: imageName}})
	if err != nil {
		return libruntime.Image{}, fmt.Errorf("pulling image failed: %v", err)
	}

	return libruntime.Image{Name: img.ImageRef}, nil
}

func (cr *CRIRuntime) Create(ctx context.Context, containerName string, imageName string, OCISpecs *runtimespecs.Spec) (*libruntime.Container, error) {

	//TODO : Check if pod exist
	podID, err := cr.CreateSandbox(ctx, "pod"+containerName, defaultPodID, defaultSanboxConfig)
	if err != nil {
		return nil, err
	}

	log.Debug(podID)
	//TODO: Instead of reading from config file, build ContainerConfig from specs.

	config, err := loadContainerConfig(defaultContainerConfig)
	if err != nil {
		return nil, err
	}
	config.Metadata.Name = containerName
	sandboxConfig, err := loadPodSandboxConfig(defaultSanboxConfig)
	if err != nil {
		return nil, err
	}
	//	startTime := time.Now()
	r, err := (*cr.RuntimeClient).CreateContainer(ctx, &pb.CreateContainerRequest{
		PodSandboxId:  podID,
		Config:        config,
		SandboxConfig: sandboxConfig,
	})
	if err != nil {
		return nil, err
	}
	//	totalTime := time.Now().Sub(startTime)
	//	log.Infof("Container Create time %s ", totalTime.String())
	log.Debug(r.ContainerId)

	return &libruntime.Container{ID: r.ContainerId, PodID: podID}, nil
}

func (cr *CRIRuntime) Run(ctx context.Context, containerName string, imageName string, OCISpecs *runtimespecs.Spec) (<-chan interface{}, *libruntime.Container, error) {
	ctr, err := cr.Create(ctx, containerName, imageName, OCISpecs)
	if err != nil {
		return nil, nil, err
	}
	//	startTime := time.Now()
	err = cr.Start(ctx, ctr)
	if err != nil {
		return nil, nil, err
	}
	//	totalTime := time.Now().Sub(startTime)
	//	log.Infof("Container Run time %s ", totalTime.String())
	return nil, ctr, nil
}
func (cr *CRIRuntime) Stop(ctx context.Context, ctr *libruntime.Container) error {
	if ctr.ID == "" {
		return fmt.Errorf("Container ID cannot be empty")
	}
	//	startTime := time.Now()
	_, err := (*cr.RuntimeClient).StopContainer(ctx, &pb.StopContainerRequest{
		ContainerId: ctr.ID,
		Timeout:     10,
	})
	if err != nil {
		return err
	}
	//	totalTime := time.Now().Sub(startTime)
	//	log.Infof("Container Stop time %s ", totalTime.String())
	return err
}
func (cr *CRIRuntime) Delete(ctx context.Context, ctr *libruntime.Container) error {
	if ctr.ID == "" {
		return fmt.Errorf("Container ID cannot be empty")
	}
	//	startTime := time.Now()
	_, err := (*cr.RuntimeClient).RemoveContainer(ctx, &pb.RemoveContainerRequest{
		ContainerId: ctr.ID,
	})
	if err != nil {
		return nil
	}
	//	totalTime := time.Now().Sub(startTime)
	//	log.Infof("Container Delete time %s ", totalTime.String())
	err = cr.StopPodSandbox(ctx, ctr.PodID)
	if err != nil {
		return err
	}
	err = cr.RemovePodSandbox(ctx, ctr.PodID)

	return err
}
func (cr *CRIRuntime) Runnable(context.Context, *libruntime.Container) error {
	return nil
}
func (cr *CRIRuntime) Start(ctx context.Context, ctr *libruntime.Container) error {
	if ctr.ID == "" {
		return fmt.Errorf("Container ID cannot be empty")
	}
	_, err := (*cr.RuntimeClient).StartContainer(ctx, &pb.StartContainerRequest{
		ContainerId: ctr.ID,
	})
	return err
}
func (cr *CRIRuntime) Exec(context.Context, libruntime.Container, []string) error {
	return nil
}
func (cr *CRIRuntime) Wait(context.Context, *libruntime.Container) (<-chan interface{}, error) {
	return nil, nil
}
func (cr *CRIRuntime) GetContainer(context.Context, string) (*libruntime.Container, error) {
	return nil, nil
}

func (cr *CRIRuntime) CreateSandbox(ctx context.Context, podName, podID, configFilePath string) (string, error) {
	config, err := loadPodSandboxConfig(configFilePath)
	if err != nil {
		return "", err
	}
	config.Metadata.Name = podName

	r, err := (*cr.RuntimeClient).RunPodSandbox(ctx, &pb.RunPodSandboxRequest{Config: config})
	if err != nil {
		return "", err
	}
	log.Debug("Pod Created: ", r.PodSandboxId)

	return r.PodSandboxId, nil
}

func (cr *CRIRuntime) StopPodSandbox(ctx context.Context, ID string) error {
	if ID == "" {
		return fmt.Errorf("ID cannot be empty")
	}
	_, err := (*cr.RuntimeClient).StopPodSandbox(ctx, &pb.StopPodSandboxRequest{PodSandboxId: ID})

	return err
}

func (cr *CRIRuntime) RemovePodSandbox(ctx context.Context, ID string) error {
	if ID == "" {
		return fmt.Errorf("ID cannot be empty")
	}
	_, err := (*cr.RuntimeClient).RemovePodSandbox(ctx, &pb.RemovePodSandboxRequest{PodSandboxId: ID})
	return err
}
