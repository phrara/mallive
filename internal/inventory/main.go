package main

import (
	"log"

	"github.com/phrara/mallive/common/config"
	"github.com/phrara/mallive/common/genproto/inventorypb"
	"github.com/phrara/mallive/common/server"
	"github.com/phrara/mallive/inventory/ports"
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

	switch serverToRun {
	case "grpc":
		server.RunGRPCServer(serviceName, func(server *grpc.Server) {
			inventorypb.RegisterInventoryServiceServer(server, ports.NewInventoryGRPCServer())
		})
	case "http":
		// TODO
	default:
		panic("Unexpected Server Type")
	}



}