package main

import (
	"context"
	"log"

	"github.com/phrara/mallive/common/config"
	"github.com/phrara/mallive/common/genproto/inventorypb"
	"github.com/phrara/mallive/common/server"
	"github.com/phrara/mallive/inventory/ports"
	"github.com/phrara/mallive/inventory/service"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func init() {
	if err := config.NewViperConfig(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	serviceName := viper.GetString("inventory.serviceName")
	serverToRun := viper.GetString("inventory.serverToRun")

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	app := service.NewApplication(ctx)

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