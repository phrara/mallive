package adapters

import (
	"context"

	"github.com/phrara/mallive/common/genproto/orderpb"
	"github.com/phrara/mallive/common/logging"
	"github.com/phrara/mallive/common/tracing"
	"google.golang.org/grpc/status"
)

type OrderGRPC struct {
	GrpcServiceName string

	client orderpb.OrderServiceClient
}

func NewOrderGRPC(client orderpb.OrderServiceClient) *OrderGRPC {
	return &OrderGRPC{GrpcServiceName: "OrderGRPC", client: client}
}

func (o OrderGRPC) UpdateOrder(ctx context.Context, order *orderpb.Order) (err error) {
	ctx, span := tracing.Start(ctx, "order_grpc.update_order")
	defer span.End()

	grpcLog := logging.WhenGRPC(ctx, o.GrpcServiceName, order)
	resp, err := o.client.UpdateOrder(ctx, order)
	grpcLog(resp, err)

	return status.Convert(err).Err()
}
