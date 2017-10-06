package libcrio

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

func (cr *CRIORuntime) Version(ctx context.Context) string {

	r, err := (*cr.RuntimeClient).Version(ctx, &pb.VersionRequest{})
	if err != nil {
		return err.Error()
	}
	log.Debug(r.String())
	return "CRIO " + r.GetRuntimeApiVersion() + " (Runtime: " + r.GetRuntimeName() + " " + r.GetVersion() + ")"
}

func (cr *CRIORuntime) Pull(ctx context.Context, imageName string) (libruntime.Image, error) {

	img, err := (*cr.ImageClient).PullImage(ctx, &pb.PullImageRequest{Image: &pb.ImageSpec{Image: imageName}})
	if err != nil {
		return libruntime.Image{}, fmt.Errorf("pulling image failed: %v", err)
	}

	return libruntime.Image{Name: img.ImageRef}, nil
}

func (cr *CRIORuntime) Create(ctx context.Context, containerName string, imageName string, OCISpecs *runtimespecs.Spec) (*libruntime.Container, error) {

	//TODO : Check if pod exist
	podID, err := cr.CreateSandbox(ctx, defaultSandboxName, defaultPodID, defaultSanboxConfig)
	if err != nil {
		return nil, err
	}

	log.Debug(podID)
	//TODO: Instead of reading from config file, build ContainerConfig from specs.

	config, err := loadContainerConfig(defaultContainerConfig)
	if err != nil {
		return nil, err
	}

	r, err := (*cr.RuntimeClient).CreateContainer(ctx, &pb.CreateContainerRequest{
		PodSandboxId: podID,
		Config:       config,
	})
	if err != nil {
		return nil, err
	}
	log.Debug(r.ContainerId)

	return &libruntime.Container{ID: r.ContainerId, PodID: podID}, nil
}

func (cr *CRIORuntime) Run(ctx context.Context, containerName string, imageName string, OCISpecs *runtimespecs.Spec) (<-chan interface{}, *libruntime.Container, error) {
	ctr, err := cr.Create(ctx, containerName, imageName, OCISpecs)
	if err != nil {
		return nil, nil, err
	}

	err = cr.Start(ctx, ctr)
	if err != nil {
		return nil, nil, err
	}

	return nil, ctr, nil
}
func (cr *CRIORuntime) Stop(ctx context.Context, ctr *libruntime.Container) error {
	if ctr.ID == "" {
		return fmt.Errorf("Container ID cannot be empty")
	}
	_, err := (*cr.RuntimeClient).StopContainer(ctx, &pb.StopContainerRequest{
		ContainerId: ctr.ID,
		Timeout:     10,
	})
	if err != nil {
		return err
	}

	return err
}
func (cr *CRIORuntime) Delete(ctx context.Context, ctr *libruntime.Container) error {
	if ctr.ID == "" {
		return fmt.Errorf("Container ID cannot be empty")
	}
	_, err := (*cr.RuntimeClient).RemoveContainer(ctx, &pb.RemoveContainerRequest{
		ContainerId: ctr.ID,
	})
	if err != nil {
		return nil
	}
	err = cr.StopPodSandbox(ctx, ctr.PodID)
	if err != nil {
		return err
	}
	err = cr.RemovePodSandbox(ctx, ctr.PodID)

	return err
}
func (cr *CRIORuntime) Runnable(context.Context, *libruntime.Container) error {
	return nil
}
func (cr *CRIORuntime) Start(ctx context.Context, ctr *libruntime.Container) error {
	if ctr.ID == "" {
		return fmt.Errorf("Container ID cannot be empty")
	}
	_, err := (*cr.RuntimeClient).StartContainer(ctx, &pb.StartContainerRequest{
		ContainerId: ctr.ID,
	})
	return err
}
func (cr *CRIORuntime) Exec(context.Context, libruntime.Container, []string) error {
	return nil
}
func (cr *CRIORuntime) Wait(context.Context, *libruntime.Container) (<-chan interface{}, error) {
	return nil, nil
}
func (cr *CRIORuntime) GetContainer(context.Context, string) (*libruntime.Container, error) {
	return nil, nil
}

func (cr *CRIORuntime) CreateSandbox(ctx context.Context, podName, podID, configFilePath string) (string, error) {
	config, err := loadPodSandboxConfig(configFilePath)
	if err != nil {
		return "", err
	}

	r, err := (*cr.RuntimeClient).RunPodSandbox(ctx, &pb.RunPodSandboxRequest{Config: config})
	if err != nil {
		return "", err
	}
	log.Debug("Pod Created: ", r.PodSandboxId)

	return r.PodSandboxId, nil
}

func (cr *CRIORuntime) StopPodSandbox(ctx context.Context, ID string) error {
	if ID == "" {
		return fmt.Errorf("ID cannot be empty")
	}
	_, err := (*cr.RuntimeClient).StopPodSandbox(ctx, &pb.StopPodSandboxRequest{PodSandboxId: ID})

	return err
}

func (cr *CRIORuntime) RemovePodSandbox(ctx context.Context, ID string) error {
	if ID == "" {
		return fmt.Errorf("ID cannot be empty")
	}
	_, err := (*cr.RuntimeClient).RemovePodSandbox(ctx, &pb.RemovePodSandboxRequest{PodSandboxId: ID})
	return err
}
