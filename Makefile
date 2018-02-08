.PHONY: build
build:
	go build

.PHONY: crio
crio:
	cd ./hack/crio-build && make 
	
.PHONY: cri-containerd
cri-containerd:
	cd ./hack/cri-ctrd-build && make

.PHONY: clean
clean:
	go clean
	cd ./hack/crio-build && make clean
	cd ./hack/cri-ctrd-build && make clean