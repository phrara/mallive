package service

import (
	"context"

	"github.com/phrara/mallive/common/broker"
	grpcClient "github.com/phrara/mallive/common/client"
	"github.com/phrara/mallive/common/metrics"
	"github.com/phrara/mallive/order/adapters"
	_ "github.com/phrara/mallive/order/adapters"
	"github.com/phrara/mallive/order/adapters/grpc"
	"github.com/phrara/mallive/order/app"
	"github.com/phrara/mallive/order/app/command"
	"github.com/phrara/mallive/order/app/query"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)


func NewApplication(ctx context.Context) (app.Application, func()) {
	if inventoryGRPCClient, inventoryGRPCClientClose, err := grpcClient.NewInventoryGRPCClient(ctx); err != nil {
		panic(err)	
	} else {
		// MQ
		ch, closeCh := broker.Connect(
		viper.GetString("rabbitmq.user"),
		viper.GetString("rabbitmq.password"),
		viper.GetString("rabbitmq.host"),
		viper.GetString("rabbitmq.port"),
		)
		inventoryGRPC := grpc.NewInventoryGRPC(inventoryGRPCClient)
		return newApplication(
			ctx, inventoryGRPC, ch,
		), func() {
			_ = inventoryGRPCClientClose()
			_ = closeCh()
			_ = ch.Close()
		}
	}
}

func newApplication(_ context.Context, inventoryGRPC query.InventoryService, ch *amqp.Channel) app.Application {
	//orderRepo := adapters.NewMemoryOrderRepository()
	// mongoClient := newMongoClient()
	// orderRepo := adapters.NewOrderRepositoryMongo(mongoClient)
	orderRepo := adapters.NewMemoryOrderRepository()
	logger := logrus.NewEntry(logrus.StandardLogger())
	metricClient := metrics.TodoMetrics{}
	return app.Application{
		Commands: app.Commands{
			CreateOrder: command.NewCreateOrderHandler(
				orderRepo, inventoryGRPC, ch, logger, metricClient),
			UpdateOrder: command.NewUpdateOrderHandler(
				orderRepo, logger, metricClient),
		},
		Queries: app.Queries{
			GetCustomerOrder: query.NewGetCustomerOrderHandler(
				orderRepo, logger, metricClient),
		},
	}
}