package ports

import (
	"context"

	"github.com/phrara/mallive/common/genproto/orderpb"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

/*
type OrderServiceServer interface {
	CreateOrder(context.Context, *CreateOrderRequest) (*emptypb.Empty, error)
	GetOrder(context.Context, *GetOrderRequest) (*Order, error)
	UpdateOrder(context.Context, *Order) (*emptypb.Empty, error)
}
**/

// 强制检查接口实现
var _ orderpb.OrderServiceServer = (*OrderGRPCServer)(nil)

type OrderGRPCServer struct {

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

func NewOrderGRPCServer() *OrderGRPCServer {
	return &OrderGRPCServer{} 
}