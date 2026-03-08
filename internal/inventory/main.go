package main

import (
	"context"

	_ "github.com/phrara/mallive/common/config"
	"github.com/phrara/mallive/common/discovery"
	"github.com/phrara/mallive/common/genproto/inventorypb"
	"github.com/phrara/mallive/common/logging"
	"github.com/phrara/mallive/common/server"
	"github.com/phrara/mallive/common/tracing"
	"github.com/phrara/mallive/inventory/ports"
	"github.com/phrara/mallive/inventory/service"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func init() {
	logging.Init()
}

func main() {
	serviceName := viper.GetString("inventory.serviceName")
	serverToRun := viper.GetString("inventory.serverToRun")

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	// Tracing: 链路追踪
	tracingShutdown, err := tracing.InitJaegerGrpcProvider(ctx, viper.GetString("jaeger.otlp-grpc"), serviceName)
	if err != nil {
		logrus.Fatal(err)
	}
	defer tracingShutdown(ctx)

	// App
	app := service.NewApplication(ctx)

	// Service Register
	deregisterFunc, err := discovery.RegisterToConsul(ctx, serviceName)
	if err != nil {
		logrus.Fatal(err)
	}
	defer func ()  {
		_ = deregisterFunc()
	}()

	switch serverToRun {
	case "grpc":
		server.RunGRPCServer(serviceName, func(server *grpc.Server) {
			srv := ports.NewInventoryGRPCServer(app)
			inventorypb.RegisterInventoryServiceServer(server, srv)
		})
	case "http":
		// TODO
	default:
		panic("Unexpected Server Type")
	}



}