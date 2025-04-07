package ports

import (
	"context"

	"github.com/phrara/mallive/common/genproto/orderpb"
	"github.com/phrara/mallive/order/app"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// 强制检查接口实现
var _ orderpb.OrderServiceServer = (*OrderGRPCServer)(nil)

type OrderGRPCServer struct {
	app app.Application
}


func (o *OrderGRPCServer) CreateOrder(context.Context, *orderpb.CreateOrderRequest) (*emptypb.Empty, error) {
	return nil, nil
}
func (o *OrderGRPCServer) GetOrder(context.Context, *orderpb.GetOrderRequest) (*orderpb.Order, error) {
	return nil, nil
}
func (o *OrderGRPCServer) UpdateOrder(context.Context, *orderpb.Order) (*emptypb.Empty, error) {
	return nil, nil
}

func NewOrderGRPCServer(app app.Application) *OrderGRPCServer {
	return &OrderGRPCServer{
		app: app,
	} 
}