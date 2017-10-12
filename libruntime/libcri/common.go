package libcri

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/kunalkushwaha/ctr-powertest/libruntime"
	log "github.com/sirupsen/logrus"
	pb "k8s.io/kubernetes/pkg/kubelet/apis/cri/v1alpha1/runtime"
)

type CRIRuntime struct {
	RuntimeClient *pb.RuntimeServiceClient
	ImageClient   *pb.ImageServiceClient
	RuntimeServer *pb.RuntimeServiceServer
}

var (
	defaultSandboxName     = "powertestPod"
	defaultPodID           = "powertestPod"
	defaultSanboxConfig    = "contrib/crio/sandbox_config.json"
	defaultContainerConfig = "contrib/crio/container_config.json"
	defaultTimeout         = time.Duration(time.Second * 10)
)

func GetNewCRIRuntime(config libruntime.RuntimeConfig, startServer bool) (libruntime.Runtime, error) {
	var (
		runtimeClient *pb.RuntimeServiceClient
		imageClient   *pb.ImageServiceClient

		err error
	)
	//localConfig := runtime2containerd(config)

	if startServer {

	} else {
		//cri - containerd
		runtimeClient, err = GetNewRuntimeClient(config.RuntimeEndpoint, time.Duration(100*time.Second))
		if err != nil {
			log.Error("Could not initialize runtimeClient")
			return nil, err
		}
		imageClient, err = GetNewImageClient(config.RuntimeEndpoint, time.Duration(100*time.Second))
		if err != nil {
			log.Error("Could not initialize runtimeClient")
			return nil, err
		}
	}

	return &CRIRuntime{RuntimeClient: runtimeClient, ImageClient: imageClient}, nil
}

func openFile(path string) (*os.File, error) {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("config at %s not found", path)
		}
		return nil, err
	}
	return f, nil
}

func loadPodSandboxConfig(path string) (*pb.PodSandboxConfig, error) {
	f, err := openFile(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var config pb.PodSandboxConfig
	if err := json.NewDecoder(f).Decode(&config); err != nil {
		return nil, err
	}
	return &config, nil
}

func loadContainerConfig(path string) (*pb.ContainerConfig, error) {
	f, err := openFile(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var config pb.ContainerConfig
	if err := json.NewDecoder(f).Decode(&config); err != nil {
		return nil, err
	}
	return &config, nil
}
