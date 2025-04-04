package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/phrara/mallive/order/ports"
)

// 强制检查接口实现
var _ ports.ServerInterface = (*OrderHTTPServer)(nil)

type OrderHTTPServer struct {

}


func (o *OrderHTTPServer) PostCustomerCustomerIDOrder(c *gin.Context, customerID string) {

}



func (o *OrderHTTPServer) GetCustomerCustomerIDOrderOrderID(c *gin.Context, customerID string, orderID string) {
	c.String(http.StatusOK, "customerID: %v, orderID: %v", customerID, orderID)
}

func NewOrderHTTPServer() *OrderHTTPServer {
	return &OrderHTTPServer{}
}