package grpc

import (
	"context"
	"errors"

	"github.com/phrara/mallive/common/genproto/inventorypb"
	"github.com/phrara/mallive/common/genproto/orderpb"
	"github.com/sirupsen/logrus"
)

type InventoryGRPC struct {
	client inventorypb.InventoryServiceClient
}


func NewInventoryGRPC(client inventorypb.InventoryServiceClient) *InventoryGRPC {
	return &InventoryGRPC{client: client}
}

func (s InventoryGRPC) CheckItemsInventory(ctx context.Context, items []*orderpb.ItemWithQuantity) (*inventorypb.CheckItemsInventoryResponse, error) {
	if items == nil {
		return nil, errors.New("grpc items cannot be nil")
	}
	resp, err := s.client.CheckItemsInventory(ctx, &inventorypb.CheckItemsInventoryRequest{Items: items})
	logrus.Info("stock_grpc response", resp)
	return resp, err
}

func (s InventoryGRPC) GetItems(ctx context.Context, itemIDs []string) ([]*orderpb.Item, error) {
	resp, err := s.client.GetItems(ctx, &inventorypb.GetItemsRequest{ItemIDs: itemIDs})
	if err != nil {
		return nil, err
	}
	return resp.Items, nil
}
