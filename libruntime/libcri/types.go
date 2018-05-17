package libcri

import pb "k8s.io/kubernetes/pkg/kubelet/apis/cri/runtime/v1alpha2"

//Pod keeps the pod info
//FIXME: Think of better name
type Pod struct {
	ContainerID string
	PodID       string
}

/*
	apiVersion: v1
kind: Pod
metadata:
  name: dns-frontend
  labels:
    name: dns-frontend
spec:
  containers:
    - name: dns-frontend
      image: k8s.gcr.io/example-dns-frontend:v1
      command:
        - python
        - client.py
        - http://dns-backend.development.svc.cluster.local:8000
      imagePullPolicy: Always
restartPolicy: Never
*/

//PodConfig stores pod declaration
type PodConfig struct {
	APIVersion string   `yaml:"apiVersion"`
	Kind       string   `yaml:"kind"`
	Metadata   Metadata `yaml:"metadata"`
	Spec       Spec     `yaml:"spec"`
}

//Metadata is pod configuration
type Metadata struct {
	Name        string            `yaml:"name"`
	Image       string            `yaml:"image"`
	Hostname    string            `yaml:"hostname"`
	DNSConfig   DNSConfig         `yaml:"dns_config"`
	Annotations map[string]string `yaml:"annotations"`
}

//DNSConfig defines dns
type DNSConfig struct {
	// List of DNS servers of the cluster.
	Servers []string `yaml:"servers,omitempty"`
	// List of DNS search domains of the cluster.
	Searches []string `yaml:"searches,omitempty"`
	// List of DNS options. See https://linux.die.net/man/5/resolv.conf
	// for all available options.
	Options []string `yaml:"options,omitempty"`
}

type Spec struct {
	Containers []Container `yaml:"containers"`
}

//Container configuration
type Container struct {
	Name         string         `yaml:"name"`
	Image        string         `yaml:"image"`
	Command      []string       `yaml:"command"`
	Ports        []Port         `yaml:"ports"`
	VolumeMounts []VolumeMount  `yaml:"volumeMount"`
	Env          []EnvVariables `yaml:"env"`
	Args         []string       `yaml:"args"`
}

type VolumeMount struct {
	MountPath string `yaml:"mountPath"`
	Name      string `yaml:"name"`
}

type Port struct {
	ContainerPort string `yaml:"containerPort"`
	Protocol      string `yaml:"protocol"`
}

//EnvVariables for storing environment variables
type EnvVariables struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

// PodConfig2PodSandboxConfig builds PodSandboxConfig from podconfig
func PodConfig2PodSandboxConfig(podConfig *PodConfig) (*pb.PodSandboxConfig, error) {
	sandboxConfig := new(pb.PodSandboxConfig)
	sandboxConfig.Metadata = new(pb.PodSandboxMetadata)
	sandboxConfig.Linux = new(pb.LinuxPodSandboxConfig)
	sandboxConfig.Linux.SecurityContext = new(pb.LinuxSandboxSecurityContext)
	sandboxConfig.Linux.SecurityContext.NamespaceOptions = new(pb.NamespaceOption)

	sandboxConfig.Metadata.Name = podConfig.Metadata.Name
	sandboxConfig.Metadata.Uid = podConfig.Metadata.Name
	sandboxConfig.Metadata.Namespace = "ctr-powertest"
	sandboxConfig.Hostname = podConfig.Metadata.Name
	//sandboxConfig.Labels
	//sandboxConfig.Linux.SecurityContext.Privileged
	//sandboxConfig.Linux.SecurityContext.RunAsGroup
	//sandboxConfig.Linux.SecurityContext.RunAsUser
	//sandboxConfig.Linux.SecurityContext.ReadonlyRootfs
	sandboxConfig.Linux.SecurityContext.NamespaceOptions.Network = pb.NamespaceMode_NODE
	//sandboxConfig.LogDirectory
	return sandboxConfig, nil
}

// PodConfig2ContainerConfig builds ContainerConfig from podconfig
func PodConfig2ContainerConfig(podConfig *PodConfig) (*pb.ContainerConfig, error) {
	var containerConfig pb.ContainerConfig
	containerConfig.Image = &pb.ImageSpec{}
	containerConfig.Metadata = &pb.ContainerMetadata{}
	containerConfig.Annotations = map[string]string{}

	containerConfig.Image.Image = podConfig.Spec.Containers[0].Image
	containerConfig.Metadata.Name = podConfig.Spec.Containers[0].Name
	containerConfig.Stdin = true
	containerConfig.Tty = true
	containerConfig.Args = podConfig.Spec.Containers[0].Args
	//	containerConfig.Envs = podConfig.Spec.Containers[0].Env
	containerConfig.Command = podConfig.Spec.Containers[0].Command
	containerConfig.Annotations = podConfig.Metadata.Annotations

	return &containerConfig, nil
}
