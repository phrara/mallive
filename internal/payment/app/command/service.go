package command

import (
	"context"

	"github.com/phrara/mallive/common/genproto/orderpb"
)

type OrderService interface {
	UpdateOrder(ctx context.Context, order *orderpb.Order) error
}
