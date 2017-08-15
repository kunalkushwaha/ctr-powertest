package libcontainerd

import (
	"context"
	"net"

	"github.com/containerd/containerd/server"
	"github.com/containerd/containerd/sys"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func SetupNewServer(config server.Config) (*server.Server, error) {

	ctx := context.TODO()
	/*	file, err := os.OpenFile("logrus.log", os.O_CREATE|os.O_WRONLY, 0666)
		if err == nil {
			logrus.SetOutput(file)
		} else {
			logrus.Info("Failed to log to file, using default stderr")
		}
	*/

	serverInstance, err := server.New(ctx, &config)
	if err != nil {
		return serverInstance, err
	}

	if config.Debug.Address != "" {
		l, err := sys.GetLocalListener(config.Debug.Address, config.Debug.Uid, config.Debug.Gid)
		if err != nil {
			return serverInstance, errors.Wrapf(err, "failed to get listener for debug endpoint")
		}
		serve(ctx, l, serverInstance.ServeDebug)
	}
	if config.Metrics.Address != "" {
		l, err := net.Listen("tcp", config.Metrics.Address)
		if err != nil {
			return serverInstance, errors.Wrapf(err, "failed to get listener for metrics endpoint")
		}
		serve(ctx, l, serverInstance.ServeMetrics)
	}

	if config.GRPC.Address != "" {
		l, err := sys.GetLocalListener(config.GRPC.Address, config.GRPC.Uid, config.GRPC.Gid)
		if err != nil {
			return serverInstance, errors.Wrapf(err, "failed to get listener for main endpoint")
		}
		serve(ctx, l, serverInstance.ServeGRPC)
	}
	//log.G(ctx).Infof("containerd successfully booted in %fs", time.Since(start).Seconds())
	//	return handleSignals(ctx, signals, server)
	return serverInstance, nil
}

func serve(ctx context.Context, l net.Listener, serveFunc func(net.Listener) error) {
	path := l.Addr().String()
	logrus.Info("serving...", path)
	go func() {
		defer l.Close()
		if err := serveFunc(l); err != nil {
			logrus.Fatal("serve failure", path)
		}
	}()
}
