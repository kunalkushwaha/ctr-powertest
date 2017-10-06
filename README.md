ctr-powertest - _Container runtime test tool_
--------------

Easy and extensible tool for any OCI supported container runtime.

Useful to run same testcases on supported runtimes.

Currently supports 
- [containerd](https://github.com/containerd/containerd) 
    - Status : WIP
- [cri-o](https://github.com/kubernetes-incubator/cri-o)
    - Status : WIP

Usage:

`` ctr-powertest -r <runtime-name> <test-cases>``


Examples
```
$ sudo ./ctr-powertest -r containerd basic
INFO[0000] Running tests on containerd v1.0.0-beta.1-23-g70b353d.m
INFO[0000] TestPullContainerImage..
INFO[0000] OK..
INFO[0000] TestCreateContainers..
INFO[0000] OK..
INFO[0000] TestCreateRunningContainers..
INFO[0000] OK..
INFO[0000] TestCreateRunningNWaitContainers..
INFO[0006] OK..


$ sudo ./ctr-powertest -r crio basic
INFO[0000] Running tests on CRIO v1alpha1 (Runtime: runc 0.1.0)
INFO[0000] TestPullContainerImage..
INFO[0004] OK..
INFO[0004] TestCreateContainers..
INFO[0011] OK..
INFO[0011] TestCreateRunningContainers..
INFO[0019] OK..
INFO[0019] TestCreateRunningNWaitContainers..
INFO[0027] OK..
```

Usage:

```
$ ./ctr-powertest -h
container runtime testing tool

Usage:
  ctr-powertest [flags]
  ctr-powertest [command]

Available Commands:
  basic       runs basic tests
  stress      Run container tests in parallel (Stress Test)

Flags:
  -d, --debug            debug mode (default false)
  -r, --runtime string   runtime [ containerd|crio ] (default "containerd")

Use "ctr-powertest [command] --help" for more information about a command.
```

#### Current Status:

Still under development.
- ``Containerd`` support is 80% completed.
- ``CRIO`` support is 20% completed.

Contribution , Feedback and reviews are welcome :).
