package service

import (
	"context"

	"github.com/phrara/mallive/common/metrics"
	"github.com/phrara/mallive/inventory/adapters"
	"github.com/phrara/mallive/inventory/app"
	"github.com/phrara/mallive/inventory/app/query"
	"github.com/sirupsen/logrus"
)


func NewApplication(ctx context.Context) app.Application {
	logger := logrus.NewEntry(logrus.StandardLogger())
	metricsClient := metrics.TodoMetrics{}
	repo := adapters.NewMemoryInventoryRepository()
	return app.Application{
		Commands: app.Commands{},
		Queries: app.Queries{
			GetItems: query.NewGetItemsHandler(
				repo,
				logger,
				metricsClient,
			),
			CheckIfItemsInInventory: query.NewCheckIfItemsInInventoryHandler(
				repo,
				nil,
				logger,
				metricsClient,
			),
		},
	}
}