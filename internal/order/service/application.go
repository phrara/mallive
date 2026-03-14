package service

import (
	"context"
	"fmt"
	"time"

	"github.com/phrara/mallive/common/broker"
	grpcClient "github.com/phrara/mallive/common/client"
	"github.com/phrara/mallive/common/decorator"
	"github.com/phrara/mallive/order/adapters"
	_ "github.com/phrara/mallive/order/adapters"
	"github.com/phrara/mallive/order/adapters/grpc"
	"github.com/phrara/mallive/order/app"
	"github.com/phrara/mallive/order/app/command"
	"github.com/phrara/mallive/order/app/query"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)


func NewApplication(ctx context.Context, metricsClient decorator.MetricsClient) (app.Application, func()) {
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
			ctx, inventoryGRPC, ch, metricsClient,
		), func() {
			_ = inventoryGRPCClientClose()
			_ = closeCh()
			_ = ch.Close()
		}
	}
}

func newApplication(_ context.Context, inventoryGRPC query.InventoryService, ch *amqp.Channel, metricClient decorator.MetricsClient) app.Application {
	mongoClient := newMongoClient()
	orderRepo := adapters.NewOrderRepositoryMongo(mongoClient)
	logger := logrus.NewEntry(logrus.StandardLogger())
	
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

func newMongoClient() *mongo.Client {
	uri := fmt.Sprintf(
		"mongodb://%s:%s@%s:%s",
		viper.GetString("mongo.user"),
		viper.GetString("mongo.password"),
		viper.GetString("mongo.host"),
		viper.GetString("mongo.port"),
	)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	c, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	if err = c.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}
	return c
}