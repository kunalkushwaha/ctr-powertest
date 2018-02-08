#!/bin/bash

trap handle_exit EXIT

handle_exit() {
	pkill cri-containerd
	pkill containerd
}

containerd &
cri-containerd $@

wait
