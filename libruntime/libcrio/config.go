package libcrio

import (
	pb "k8s.io/kubernetes/pkg/kubelet/apis/cri/v1alpha1/runtime"
)

// CRIOOpts sets spec specific information to a newly generated OCI spec
type CRIOOpts func(s *pb.ContainerConfig) error

//GenerateContainerConfig returns crio compaitable ContainerConfig
func GenerateContainerConfig(opts ...CRIOOpts) *pb.ContainerConfig {
	/*
		TODO:
		- GetPODInfo
		- GetImage
		- GetLinuxConfig
	*/
	return nil
}

func createDefaultConfig() (*pb.ContainerConfig, error) {
	/*	config := pb.ContainerConfig{
			Metadata: {
				Name: "container1",
				Attmept: 1,
			},
			Image: {
				Image: "redis:alpine",
			},

			WorkingDir: "/",
			Privilaged: true,
			LogPath: "",
			Stdin: false,
			StdinOnce: false,
			Tty: false,
			Linux: {
				Resources: {
					"cpu_period": 10000,
					"cpu_quota": 20000,
					"cpu_shares": 512,
					"oom_score_adj": 30
				},
				SecurityContext: {
					"readonly_rootfs": false,
					"capabilities": {
						"add_capabilities": [
							"setuid",
							"setgid"
						],
						"drop_capabilities": [
						]
					},
					"selinux_options": {
						"user": "system_u",
						"role": "system_r",
						"type": "container_t",
						"level": "s0:c4,c5"
					}
				}
			}
		}
	*/
	return nil, nil
}
