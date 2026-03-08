package ports

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/phrara/mallive/common/genproto/orderpb"
	"github.com/phrara/mallive/order/app"
	"github.com/phrara/mallive/order/app/command"
	"github.com/phrara/mallive/order/app/query"
	"github.com/phrara/mallive/order/convertor"
	domain "github.com/phrara/mallive/order/domain/order"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// 强制检查接口实现
var _ orderpb.OrderServiceServer = (*OrderGRPCServer)(nil)

type OrderGRPCServer struct {
	app app.Application

	// 匿名字段注入，实现接口
	orderpb.OrderServiceServer
}


func (o *OrderGRPCServer) CreateOrder(ctx context.Context, request *orderpb.CreateOrderRequest) (*emptypb.Empty, error) {
	_, err := o.app.Commands.CreateOrder.Handle(ctx, command.CreateOrder{
		CustomerID: request.CustomerID,
		Items:      convertor.NewItemWithQuantityConvertor().ProtosToEntities(request.Items),
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &empty.Empty{}, nil
}
func (o *OrderGRPCServer) GetOrder(ctx context.Context, request *orderpb.GetOrderRequest) (*orderpb.Order, error) {
	order, err := o.app.Queries.GetCustomerOrder.Handle(ctx, query.GetCustomerOrder{
		CustomerID: request.CustomerID,
		OrderID:    request.OrderID,
	})
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	return convertor.NewOrderConvertor().EntityToProto(order), nil
}
func (o *OrderGRPCServer) UpdateOrder(ctx context.Context, request *orderpb.Order) (*emptypb.Empty, error) {
	logrus.Infof("order_grpc||request_in||request=%+v", request)
	order, err := domain.NewOrder(
		request.ID,
		request.CustomerID,
		request.Status,
		request.PaymentLink,
		convertor.NewItemConvertor().ProtosToEntities(request.Items))
	if err != nil {
		err = status.Error(codes.Internal, err.Error())
		return nil, err
	}
	_, err = o.app.Commands.UpdateOrder.Handle(ctx, command.UpdateOrder{
		Order: order,
		UpdateFn: func(ctx context.Context, order *domain.Order) (*domain.Order, error) {
			return order, nil
		},
	})
	return nil, err
}

func NewOrderGRPCServer(app app.Application) *OrderGRPCServer {
	return &OrderGRPCServer{
		app: app,
	} 
}