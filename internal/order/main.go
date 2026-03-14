package main

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/phrara/mallive/common/broker"
	_ "github.com/phrara/mallive/common/config"
	"github.com/phrara/mallive/common/discovery"
	"github.com/phrara/mallive/common/genproto/orderpb"
	"github.com/phrara/mallive/common/logging"
	"github.com/phrara/mallive/common/metrics"
	"github.com/phrara/mallive/common/server"
	"github.com/phrara/mallive/common/tracing"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/phrara/mallive/order/infrastructure/consumer"
	"github.com/phrara/mallive/order/ports"
	"github.com/phrara/mallive/order/service"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)


func init() {
	logging.Init()
}


func main() {
	serviceName := viper.GetString("order.serviceName")
	// GRPC

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	// Tracing: 链路追踪
	tracingShutdown, err := tracing.InitJaegerGrpcProvider(ctx, viper.GetString("jaeger.otlp-grpc"), serviceName)
	if err != nil {
		logrus.Fatal(err)
	}
	defer tracingShutdown(ctx)
	
	metricClient := metrics.NewPrometheusMetricsClient(
		viper.GetString("order.serviceName"))

	// App
    app, closeF := service.NewApplication(ctx, metricClient)
	defer closeF()
	
	// Service Register
	deregisterFunc, err := discovery.RegisterToConsul(ctx, serviceName)
	if err != nil {
		logrus.Fatal(err)
	}
	defer func ()  {
		_ = deregisterFunc()
	}()
		
	// Message Que
	ch, closeMQ := broker.Connect(
		viper.GetString("rabbitmq.user"),
		viper.GetString("rabbitmq.password"),
		viper.GetString("rabbitmq.host"),
		viper.GetString("rabbitmq.port"),
	)
	defer func() {
		_ = ch.Close()
		_ = closeMQ()
	}()

	// 消费订单支付完成事件
	go consumer.NewConsumer(app).Listen(ch)

	go server.RunGRPCServer(serviceName, func(server *grpc.Server) {
		srv := ports.NewOrderGRPCServer(app)
		orderpb.RegisterOrderServiceServer(server, srv)
	})

	// HTTP
	server.RunHTTPServer(serviceName, func(router *gin.Engine) {
		router.StaticFile("/success", "../../public/success.html")
		srv := ports.NewOrderHTTPServer(app)
		ports.RegisterHandlersWithOptions(router, srv, ports.GinServerOptions{
			BaseURL: "/api",
			Middlewares: nil,
			ErrorHandler: nil,
		})

		// Prometheus
		router.GET("/metrics", gin.WrapH(promhttp.HandlerFor(
			metricClient.GetPromRegistry(),
			promhttp.HandlerOpts{},
		)))
	})

}