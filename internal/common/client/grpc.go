package client

import (
	"context"
	"errors"
	"net"
	"time"

	"github.com/phrara/mallive/common/discovery"
	"github.com/phrara/mallive/common/genproto/orderpb"
	"github.com/phrara/mallive/common/genproto/inventorypb"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// getGRPCAddr 获取 gRPC 服务地址
// K8s 环境使用 Service DNS，本地开发使用配置文件
func getGRPCAddr(serviceName string) string {
	// K8s 环境使用 Service DNS
	if discovery.IsK8sEnvironment() {
		return discovery.GetServiceDNS(serviceName, 5003)
	}
	// 本地开发使用 Consul
	addr, err := discovery.GetServiceAddr(context.Background(), serviceName)
	if err != nil {
		logrus.Warnf("failed to get addr from consul for %s: %v", serviceName, err)
		return viper.Sub(serviceName).GetString("grpcServer.address")
	}
	return addr
}

func NewInventoryGRPCClient(ctx context.Context) (client inventorypb.InventoryServiceClient, close func() error, err error) {
	// 使用 K8s Service DNS 或 Consul
	grpcAddr := "inventory.default.svc.cluster.local:5003"

	if !waitForGRPC(grpcAddr, viper.GetDuration("dial-grpc-timeout")*time.Second) {
		return nil, nil, errors.New("stock grpc not available")
	}

	opts := grpcDialOpts(grpcAddr)
	conn, err := grpc.NewClient(grpcAddr, opts...)
	if err != nil {
		return nil, func() error { return nil }, err
	}
	return inventorypb.NewInventoryServiceClient(conn), conn.Close, nil
}

func NewOrderGRPCClient(ctx context.Context) (client orderpb.OrderServiceClient, close func() error, err error) {
	// 使用 K8s Service DNS 或 Consul
	grpcAddr := "order.default.svc.cluster.local:5002"

	if !waitForGRPC(grpcAddr, viper.GetDuration("dial-grpc-timeout")*time.Second) {
		return nil, nil, errors.New("order grpc not available")
	}
	opts := grpcDialOpts(grpcAddr)

	conn, err := grpc.NewClient(grpcAddr, opts...)
	if err != nil {
		return nil, func() error { return nil }, err
	}
	return orderpb.NewOrderServiceClient(conn), conn.Close, nil
}

func grpcDialOpts(_ string) []grpc.DialOption {
	return []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	}
}

func waitForGRPC(addr string, timeout time.Duration) bool {
	logrus.Infof("waiting for grpc client: %s, timeout: %v seconds", addr, timeout.Seconds())
	portAvailable := make(chan struct{})
	timeoutCh := time.After(timeout)

	go func() {
		for {
			select {
			case <-timeoutCh:
				return
			default:
				// continue
			}
			conn, err := net.Dial("tcp", addr)
			if err == nil {
				conn.Close()
				close(portAvailable)
				return
			}
			time.Sleep(200 * time.Millisecond)
		}
	}()

	select {
	case <-portAvailable:
		return true
	case <-timeoutCh:
		return false
	}
}
