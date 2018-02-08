FROM golang:1.9 as bench
ARG POWERTEST_REPO=https://github.com/kunalkushwaha/ctr-powertest.git
ARG POWERTEST_COMMIT=HEAD
ARG POWERTEST_BRANCH=master
RUN go get github.com/spf13/cobra
RUN mkdir -p /go/src/github.com/kunalkushwaha && cd /go/src/github.com/kunalkushwaha && git clone -b $POWERTEST_BRANCH $POWERTEST_REPO && cd ctr-powertest && git checkout $POWERTEST_COMMIT
WORKDIR /go/src/github.com/kunalkushwaha/ctr-powertest
RUN go build
ENTRYPOINT ["./ctr-powertest","-p", "cri","-r","cri-containerd","profile"]
#ENTRYPOINT ["./ctr-powertest", "profile"]

FROM golang:1.9 as c8d
ARG CONTAINERD_REPO=https://github.com/containerd/containerd.git
ARG CONTAINERD_BRANCH=master
ARG CONTAINERD_COMMIT=HEAD
#RUN eval $(go env); curl -SLf https://github.com/containerd/containerd/releases/download/v${CONTAINERD_VERSION}/containerd-${CONTAINERD_VERSION}.${GOOS}-${GOARCH}.tar.gz | tar -zx -C /usr/local
RUN mkdir -p /go/src/github.com/containerd/ \
	&& cd /go/src/github.com/containerd \
	&& git clone -b $CONTAINERD_BRANCH $CONTAINERD_REPO containerd \
	&& cd containerd \
	&& git checkout $CONTAINERD_COMMIT \
	&& make BUILDTAGS=no_btrfs && make install

FROM golang:1.9 as cri
ARG CNI_VERSION=0.6.0
RUN go get github.com/opencontainers/runc
RUN mkdir -p /opt/cni/bin && curl -sSLf https://github.com/containernetworking/plugins/releases/download/v${CNI_VERSION}/cni-plugins-amd64-v${CNI_VERSION}.tgz | tar -zx -C /opt/cni/bin
RUN apt-get update && apt-get install -y socat iptables
RUN mkdir -p /go/src/github.com/containerd
ARG CRI_CONTAINERD_REPO=https://github.com/containerd/cri-containerd.git
ARG CRI_CONTAINERD_COMMIT=HEAD
ARG CRI_CONTAINERD_BRANCH=master
RUN cd /go/src/github.com/containerd && git clone -b $CRI_CONTAINERD_BRANCH $CRI_CONTAINERD_REPO && cd cri-containerd && git checkout $CRI_CONTAINERD_COMMIT
WORKDIR  /go/src/github.com/containerd/cri-containerd
RUN make BUILD_TAGS="" && make install
RUN mkdir -p /etc/cni/net.d
COPY 10-containerd-net.conflist /etc/cni/net.d/
COPY cri-c8d-run.sh /
COPY --from=c8d /go/src/github.com/containerd/containerd/bin/* /usr/local/bin/
ENTRYPOINT ["/cri-c8d-run.sh", "--profiling", "--profiling-addr=127.0.0.1", "--profiling-port=8080"]
