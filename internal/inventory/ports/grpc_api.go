package ports

import (
	"context"

	"github.com/phrara/mallive/common/genproto/inventorypb"
)

/*
type InventoryServiceServer interface {
	GetItems(context.Context, *GetItemsRequest) (*GetItemsResponse, error)
	CheckItemsInventory(context.Context, *CheckItemsInventoryRequest) (*CheckItemsInventoryResponse, error)
}
*/

// 强制检查接口实现
var _ inventorypb.InventoryServiceServer = (*InventoryGRPCServer)(nil)

type InventoryGRPCServer struct {

}

func (i *InventoryGRPCServer) GetItems(context.Context, *inventorypb.GetItemsRequest) (*inventorypb.GetItemsResponse, error) {
	return nil, nil
}

func (i *InventoryGRPCServer) CheckItemsInventory(context.Context, *inventorypb.CheckItemsInventoryRequest) (*inventorypb.CheckItemsInventoryResponse, error) {
	return nil, nil
}


func NewInventoryGRPCServer() *InventoryGRPCServer {
	return &InventoryGRPCServer{}
}