package main

import (
	"context"

	"github.com/phrara/mallive/common/broker"
	_ "github.com/phrara/mallive/common/config"
	"github.com/phrara/mallive/common/server"
	"github.com/phrara/mallive/common/tracing"
	"github.com/phrara/mallive/payment/infrastructure/consumer"
	"github.com/phrara/mallive/payment/ports"
	"github.com/phrara/mallive/payment/service"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	// "github.com/phrara/mallive/common/discovery"
	"github.com/phrara/mallive/common/logging"
)


func init() {
	logging.Init()
}

func main() {
	serviceName := viper.GetString("payment.serviceName")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	
	serverToRun := viper.GetString("payment.serverToRun")

	// Tracing: 链路追踪
	tracingShutdown, err := tracing.InitJaegerGrpcProvider(ctx, viper.GetString("jaeger.otlp-grpc"), serviceName)
	if err != nil {
		logrus.Fatal(err)
	}
	defer tracingShutdown(ctx)

	// App
	app, closeFunc := service.NewApplication(ctx)
	defer closeFunc()

	// MQ
	ch, closeCh := broker.Connect(
		viper.GetString("rabbitmq.user"),
		viper.GetString("rabbitmq.password"),
		viper.GetString("rabbitmq.host"),
		viper.GetString("rabbitmq.port"),
	)
	defer func() {
		_ = ch.Close()
		_ = closeCh()
	}()
	// 消费订单创建事件
	go consumer.NewConsumer(app).Listen(ch)

	
	paymentHandler := ports.NewPaymentHandler(ch)
	switch serverToRun {
	case "http":
		server.RunHTTPServer(serviceName, paymentHandler.RegisterRoutes)
	case "grpc":
		logrus.Panic("unsupported server type: grpc")
	default:
		logrus.Panic("unreachable code")
	}
}	


