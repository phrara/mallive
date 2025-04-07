package order

import (
	"context"
	"fmt"
)


type Repository interface {
	Create(context.Context, *Order) (*Order, error)
	Get(context.Context, string, string) (*Order, error)
	Update(context.Context, *Order,
		func(context.Context, *Order) (*Order, error),
	) error
}


type NotFoundError struct {
	OrderID string
}

func (n NotFoundError) Error() string {
	return fmt.Sprintf("Order %v is not found", n.OrderID)
}