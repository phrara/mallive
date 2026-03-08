package ports

import (
	"context"

	"github.com/phrara/mallive/common/genproto/inventorypb"
	"github.com/phrara/mallive/common/tracing"
	"github.com/phrara/mallive/inventory/app"
	"github.com/phrara/mallive/inventory/app/query"
	"github.com/phrara/mallive/inventory/convertor"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	// app 注入
	app app.Application 

	// 匿名字段注入，实现接口
	inventorypb.InventoryServiceServer
}

func (i *InventoryGRPCServer) GetItems(ctx context.Context, req *inventorypb.GetItemsRequest) (*inventorypb.GetItemsResponse, error) {
	_, span := tracing.Start(ctx, "GetItems")
	defer span.End()
	
	if items, err := i.app.Queries.GetItems.Handle(ctx, query.GetItems{
		ItemIDs: req.ItemIDs,
	}); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	} else {
		return &inventorypb.GetItemsResponse{
			Items: convertor.NewItemConvertor().EntitiesToProtos(items),
		}, nil
	}
}

func (i *InventoryGRPCServer) CheckItemsInventory(ctx context.Context,  req *inventorypb.CheckItemsInventoryRequest) (*inventorypb.CheckItemsInventoryResponse, error) {
	_, span := tracing.Start(ctx, "CheckIfItemsInStock")
	defer span.End()
	
	if items, err := i.app.Queries.CheckIfItemsInInventory.Handle(ctx, query.CheckIfItemsInInventory{
		Items: convertor.NewItemWithQuantityConvertor().ProtosToEntities(req.Items),
	}); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	} else {
		return &inventorypb.CheckItemsInventoryResponse{
			InStock: 1,
			Items: convertor.NewItemConvertor().EntitiesToProtos(items),
		}, nil
	}

}


func NewInventoryGRPCServer(app app.Application) *InventoryGRPCServer {
	return &InventoryGRPCServer{
		app: app,
	}
}