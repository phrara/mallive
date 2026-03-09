package service

import (
	"context"

	"github.com/phrara/mallive/common/metrics"
	"github.com/phrara/mallive/inventory/adapters"
	"github.com/phrara/mallive/inventory/app"
	"github.com/phrara/mallive/inventory/app/query"
	"github.com/phrara/mallive/inventory/infrastructure/integration"
	"github.com/phrara/mallive/inventory/infrastructure/persistent"
	"github.com/sirupsen/logrus"
)


func NewApplication(ctx context.Context) app.Application {
	logger := logrus.NewEntry(logrus.StandardLogger())
	metricsClient := metrics.TodoMetrics{}
	stripeAPI := integration.NewStripeAPI()
	repo := adapters.NewMySQLInventoryRepository(persistent.NewMySQL())
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
				stripeAPI,
				logger,
				metricsClient,
			),
		},
	}
}