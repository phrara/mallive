package server

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)



func RunHTTPServer(serviceName string, wrapper func(router *gin.Engine)) {
	addr := viper.Sub(serviceName).GetString("httpServer.address")
	RunHTTPServerOnAddress(addr, wrapper)
}

func RunHTTPServerOnAddress(addr string, wrapper func(router *gin.Engine)) {
	app := gin.New()
	wrapper(app)
	app.Group("/api")
	if err := app.Run(addr); err != nil {
		log.Fatal(err)
	}
}