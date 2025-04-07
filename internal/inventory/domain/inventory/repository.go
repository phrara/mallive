package inventory

import (
	"context"
	"fmt"

	"github.com/phrara/mallive/common/genproto/orderpb"
)


type Repository interface {
	GetItems(ctx context.Context, itemIDs []string) ([]*orderpb.Item, error)
}


type NotFoundError struct {
	Lacks []string
}

func (n NotFoundError) Error() string {
	return fmt.Sprintf("Order %v is not found", n.Lacks)
}