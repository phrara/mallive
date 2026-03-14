package discovery

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"time"
)

type Registry interface {
	Register(ctx context.Context, instanceID, serviceName, hostPort string) error
	Deregister(ctx context.Context, instanceID, serviceName string) error
	Discover(ctx context.Context, serviceName string) ([]string, error)
	HealthCheck(instanceID, serviceName string) error
}

func GenerateInstanceID(serviceName string) string {
	x := rand.New(rand.NewSource(time.Now().UnixNano())).Int()
	return fmt.Sprintf("%s-%d", serviceName, x)
}

// IsK8sEnvironment 检测是否在 K8s 环境中运行
func IsK8sEnvironment() bool {
	// 检查 K8s 特有环境变量
	_, inK8s := os.LookupEnv("KUBERNETES_SERVICE_HOST")
	return inK8s
}

// GetServiceDNS 获取服务的 K8s DNS 地址
func GetServiceDNS(serviceName string, port int) string {
	namespace := os.Getenv("SERVICE_NAMESPACE")
	if namespace == "" {
		namespace = "default"
	}
	return fmt.Sprintf("%s.%s.svc.cluster.local:%d", serviceName, namespace, port)
}
