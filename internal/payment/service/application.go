package service

import (
	"context"

	grpcClient "github.com/phrara/mallive/common/client"
	"github.com/phrara/mallive/common/metrics"
	"github.com/phrara/mallive/payment/adapters"
	"github.com/phrara/mallive/payment/app"
	"github.com/phrara/mallive/payment/app/command"
	"github.com/phrara/mallive/payment/domain"
	"github.com/phrara/mallive/payment/infrastructure/processor"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func NewApplication(ctx context.Context) (app.Application, func()) {
	orderClient, closeOrderClient, err := grpcClient.NewOrderGRPCClient(ctx)
	if err != nil {
		panic(err)
	}
	orderGRPC := adapters.NewOrderGRPC(orderClient)
	//memoryProcessor := processor.NewInmemProcessor()
	stripeProcessor := processor.NewStripeProcessor(viper.GetString("stripe-key"))
	return newApplication(ctx, orderGRPC, stripeProcessor), func() {
		_ = closeOrderClient()
	}
}

func newApplication(_ context.Context, orderGRPC command.OrderService, processor domain.Processor) app.Application {
	logger := logrus.NewEntry(logrus.StandardLogger())
	metricClient := metrics.TodoMetrics{}
	return app.Application{
		Commands: app.Commands{
			CreatePayment: command.NewCreatePaymentHandler(
				processor, orderGRPC, logger, metricClient,
			),
		},
	}
}
