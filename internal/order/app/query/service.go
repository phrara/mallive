package query

import (
	"context"

	"github.com/phrara/mallive/common/genproto/inventorypb"
	"github.com/phrara/mallive/common/genproto/orderpb"
)

type InventoryService interface {
	CheckItemsInventory(ctx context.Context, items []*orderpb.ItemWithQuantity) (*inventorypb.CheckItemsInventoryResponse, error)
	GetItems(ctx context.Context, itemIDs []string) ([]*orderpb.Item, error)
}
