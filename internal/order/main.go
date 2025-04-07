package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/phrara/mallive/common/config"
	"github.com/phrara/mallive/common/genproto/orderpb"
	"github.com/phrara/mallive/common/server"
	"github.com/phrara/mallive/order/ports"
	"github.com/phrara/mallive/order/service"
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
	
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	app := service.NewApplication(ctx)

	go server.RunGRPCServer(serviceName, func(server *grpc.Server) {
		srv := ports.NewOrderGRPCServer(app)
		orderpb.RegisterOrderServiceServer(server, srv)
	})

	// HTTP
	server.RunHTTPServer(serviceName, func(router *gin.Engine) {
		srv := NewOrderHTTPServer(app)
		ports.RegisterHandlersWithOptions(router, srv, ports.GinServerOptions{
			BaseURL: "/api",
			Middlewares: nil,
			ErrorHandler: nil,
		})
	})

}	