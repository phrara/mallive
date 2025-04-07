package service

import (
	"context"

	_ "github.com/phrara/mallive/order/adapters"
	"github.com/phrara/mallive/order/app"
)


func NewApplication(ctx context.Context) app.Application {
	// orderRepo := adapters.NewMemoryOrderRepository()
	return app.Application{
		Commands: app.Commands{},
		Queries:  app.Queries{},
	}
}