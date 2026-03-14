package discovery

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/phrara/mallive/common/discovery/consul"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func RegisterToConsul(ctx context.Context, serviceName string) (func() error, error) {
	// K8s 环境下跳过 Consul 注册，使用 K8s Service 发现
	if IsK8sEnvironment() {
		logrus.Infof("running in k8s, skip consul registration for %s", serviceName)
		return func() error { return nil }, nil
	}

	// 本地开发环境使用 Consul
	registry, err := consul.New(viper.GetString("consul.addr"))
	if err != nil {
		return func() error { return nil }, err
	}
	instanceID := GenerateInstanceID(serviceName)

	// 获取 gRPC 地址
	grpcAddr := getGrpcAddr(serviceName)

	// 服务注册
	if err := registry.Register(ctx, instanceID, serviceName, grpcAddr); err != nil {
		return func() error { return nil }, err
	}

	// 定时器: 心跳上报
	ticker := time.NewTicker(time.Second * 2)
	go func() {
		for {
			if err := registry.HealthCheck(instanceID, serviceName); err != nil {
				logrus.Panicf("no heartbeat from %s to registry, err=%v", serviceName, err)
			}
			<- ticker.C
		}
	}()
	logrus.WithFields(logrus.Fields{
		"serviceName": serviceName,
		"addr":        grpcAddr,
	}).Info("registered to consul")

	return func() error {
		ticker.Stop()
		return registry.Deregister(ctx, instanceID, serviceName)
	}, nil
}

// getGrpcAddr 获取 gRPC 服务地址
// K8s 环境优先使用 POD_IP 环境变量
func getGrpcAddr(serviceName string) string {
	// 1. 优先使用环境变量 POD_IP（K8s 自动注入）
	if podIP := os.Getenv("POD_IP"); podIP != "" {
		port := viper.Sub(serviceName).GetString("grpcServer.address")
		// 提取端口
		if idx := 0; len(port) > idx {
			for i := len(port) - 1; i >= 0; i-- {
				if port[i] == ':' {
					port = port[i+1:]
					break
				}
			}
		}
		return fmt.Sprintf("%s:%s", podIP, port)
	}

	// 2. 降级使用配置文件中的地址
	return viper.Sub(serviceName).GetString("grpcServer.address")
}

func GetServiceAddr(ctx context.Context, serviceName string) (string, error) {
	registry, err := consul.New(viper.GetString("consul.addr"))
	if err != nil {
		return "", err
	}
	addrs, err := registry.Discover(ctx, serviceName)
	if err != nil {
		return "", err
	}
	if len(addrs) == 0 {
		return "", fmt.Errorf("got empty %s addrs from consul", serviceName)
	}
	i := rand.Intn(len(addrs))
	logrus.Infof("Discovered %d instance of %s, addrs=%v", len(addrs), serviceName, addrs)
	return addrs[i], nil
}
