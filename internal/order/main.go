package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/phrara/mallive/common/config"
	"github.com/phrara/mallive/common/server"
	"github.com/phrara/mallive/common/genproto/orderpb"
	"github.com/phrara/mallive/order/ports"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)


func init() {
	if err := config.NewViperConfig(); err != nil {
		log.Fatal(err)
	}
}


func main() {
	serviceName := viper.GetString("order.serviceName")
	// GRPC
	
	go server.RunGRPCServer(serviceName, func(server *grpc.Server) {
		orderpb.RegisterOrderServiceServer(server, ports.NewOrderGRPCServer())
	})

	// HTTP
	server.RunHTTPServer(serviceName, func(router *gin.Engine) {
		ports.RegisterHandlersWithOptions(router, NewOrderHTTPServer(), ports.GinServerOptions{
			BaseURL: "/api",
			Middlewares: nil,
			ErrorHandler: nil,
		})
	})

}	