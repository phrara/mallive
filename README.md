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




## Stripe
- [https://dashboard.stripe.com/acct_1T8IaV4JZtdIYvS2/test/dashboard](https://dashboard.stripe.com/acct_1T8IaV4JZtdIYvS2/test/dashboard)

```sh
stripe login

# 监听 webhook，获取 endpoint-stripe-secret
stripe listen --forward-to localhost:8084/api/webhook
```


## Jaeger
链路追踪
```sh
kubectl apply -f manifest/jaeger.yaml
```


## OpenTelemetry

