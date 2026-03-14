# mallive
A live-streaming e-commerce platform implemented in Go


adpters -> queryHandler/commandHandler -> ports

# Middleware

## Consul

### k8s service
+ consul-server
  + 类型: ClusterIP(headless) 
  + 供consul集群之间访问
+ consul-ui
  + 类型可自定义
  + 共外部程序(client)访问api 或 访问webui

```sh
# install
helm repo add hashicorp https://helm.releases.hashicorp.com
helm repo update

helm install consul hashicorp/consul -f manifest/consul-helm-value.yaml

# uninstall
helm uninstall consul
kubectl delete pvc -l chart=consul-helm
```

## RabbitMQ

```sh
kubectl apply -f manifest/rabbitmq.yaml
```



## Mongo
```sh
kubectl apply -f manifest/mongo.yaml
```

## Mysql
```sh
kubectl apply -f manifest/mysql/mysql.yaml
```

## Redis
```sh
kubectl apply -f manifest/redis/redis.yaml
```


## Stripe
- [https://dashboard.stripe.com/acct_1T8IaV4JZtdIYvS2/test/dashboard](https://dashboard.stripe.com/acct_1T8IaV4JZtdIYvS2/test/dashboard)

```sh
# 90 days expire
stripe login

# 监听 webhook，获取 endpoint-stripe-secret
stripe listen --forward-to localhost:8084/api/webhook
```


## Jaeger
链路追踪
```sh
kubectl apply -f manifest/jaeger.yaml
```

## Prometheus
+ 使用的 k8s Prometheus Stack 套件
+ 集成Grafana和多种 CRD/Operator
+ 使用 `ServiceMonitor` CRD 实现对指定服务的监控

```sh
kubectl apply -f manifest/prometheus/monitor.yaml
```


## OpenTelemetry


# K8s 高可用部署

项目提供了完整的高可用 K8s 部署清单，位于 `manifest/apps/` 目录。

## 目录结构

```
manifest/apps/
├── configmap.yaml      # 全局配置
├── services.yaml       # 服务部署 (Deployment + Service + HPA + PDB)
├── ingress.yaml        # 入口路由
├── kustomization.yaml # Kustomize 编排
└── Makefile           # 镜像构建
```

## 高可用特性

- **多副本**: 每个服务 3 个 Pod
- **反亲和性**: Pod 分散到不同节点
- **拓扑分布**: 跨可用区分布
- **健康检查**: Liveness + Readiness Probe
- **资源限制**: CPU/Memory requests/limits
- **滚动更新**: maxSurge: 1, maxUnavailable: 0
- **自动扩缩容**: HPA (CPU 70%, Memory 80%)
- **最小可用**: PDB 保证最少 2 个副本

## 快速开始

### 1. 构建 Docker 镜像

```bash
# 构建所有服务
cd manifest/apps && make build-images

# 或单独构建
make build-order
make build-inventory
make build-payment
make build-kitchen
```

### 2. 部署到 K8s

```bash
# 使用 Kustomize 部署 (推荐)
kubectl apply -k manifest/apps

# 或直接部署
kubectl apply -f manifest/apps/
```

### 3. 验证部署

```bash
# 查看 Pod 状态
kubectl get pods -l app=order
kubectl get pods -l app=inventory
kubectl get pods -l app=payment
kubectl get pods -l app=kitchen

# 查看 HPA
kubectl get hpa

# 查看 Service
kubectl get svc | grep mallive
```

## 服务端口

| 服务 | 内部端口 | NodePort | gRPC |
|------|----------|----------|------|
| order | 8082 | 30082 | 5002/30052 |
| inventory | 5003 | 30053 | - |
| payment | 8084 | 30084 | - |
| kitchen | - | - | - |

## 配置

修改 `manifest/apps/configmap.yaml` 中的 `global.yaml` 配置来调整服务参数。

## 清理

```bash
kubectl delete -k manifest/apps
# 或
kubectl delete -f manifest/apps/
```