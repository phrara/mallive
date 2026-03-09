package query

import (
	"context"
	"strings"
	"time"

	"github.com/phrara/mallive/common/decorator"
	"github.com/phrara/mallive/common/handler/redis"
	domain "github.com/phrara/mallive/inventory/domain/inventory"
	"github.com/phrara/mallive/inventory/entity"
	"github.com/phrara/mallive/inventory/infrastructure/integration"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const (
	redisLockPrefix = "check_inventory_"
)

type CheckIfItemsInInventory struct {
	Items []*entity.ItemWithQuantity
}

type CheckIfItemsInInventoryHandler decorator.QueryHandler[CheckIfItemsInInventory, []*entity.Item]

type checkIfItemsInInventoryHandler struct {
	inventoryRepo domain.Repository
	stripeAPI *integration.StripeAPI
}

func NewCheckIfItemsInInventoryHandler(
	inventoryRepo domain.Repository,
	stripeAPI *integration.StripeAPI,
	logger *logrus.Entry,
	metricClient decorator.MetricsClient,
) CheckIfItemsInInventoryHandler {
	if inventoryRepo == nil {
		panic("nil inventoryRepo")
	}
	if stripeAPI == nil {
		panic("nil stripeAPI")
	}
	return decorator.ApplyQueryDecorators[CheckIfItemsInInventory, []*entity.Item](
		checkIfItemsInInventoryHandler{
			inventoryRepo: inventoryRepo,
			stripeAPI: stripeAPI,
		},
		logger,
		metricClient,
	)
}

// Deprecated
var stub = map[string]string{
	"1": "price_1QBYvXRuyMJmUCSsEyQm2oP7",
	"2": "price_1QBYl4RuyMJmUCSsWt2tgh6d",
}

func (h checkIfItemsInInventoryHandler) Handle(ctx context.Context, query CheckIfItemsInInventory) ([]*entity.Item, error) {
	// 施加分布式锁, 防止超卖和数据不一致
	/*
	如果不加锁：两个线程可能同时读取到“库存剩余 1”，然后各自执行 -1 操作，
		最后库存变成了 -1，但实际上卖出了两件商品。
	加锁后：同一时刻只有一个实例的请求能拿到锁。
		它处理完“检查库存 -> 扣减库存 -> 下单”这一系列操作后，
		才会释放锁，让下一个请求进来。
	*/
	if err := lock(ctx, getLockKey(query)); err != nil {
		return nil, errors.Wrapf(err, "redis lock error: key=%s", getLockKey(query))
	}
	defer func() {
		if err := unlock(ctx, getLockKey(query)); err != nil {
			logrus.Warnf("redis unlock fail, err=%v", err)
		}
	}()

	var res []*entity.Item
	for _, i := range query.Items {
		priceID, err := h.stripeAPI.GetPriceByProductID(ctx, i.ID)
		if err != nil || priceID == "" {
			return nil, err
		}
		res = append(res, &entity.Item{
			ID:       i.ID,
			Quantity: i.Quantity,
			PriceID:  priceID,
		})
	}
	if err := h.checkInventory(ctx, query.Items); err != nil {
		return nil, err
	}
	return res, nil
}

func getLockKey(query CheckIfItemsInInventory) string {
	var ids []string
	for _, i := range query.Items {
		ids = append(ids, i.ID)
	}
	return redisLockPrefix + strings.Join(ids, "_")
}

func unlock(ctx context.Context, key string) error {
	return redis.Del(ctx, redis.LocalClient(), key)
}

func lock(ctx context.Context, key string) error {
	return redis.SetNX(ctx, redis.LocalClient(), key, "1", 5*time.Minute)
}

func (h checkIfItemsInInventoryHandler) checkInventory(ctx context.Context, query []*entity.ItemWithQuantity) error {
	var ids []string
	for _, i := range query {
		ids = append(ids, i.ID)
	}
	records, err := h.inventoryRepo.GetInventory(ctx, ids)
	if err != nil {
		return err
	}
	idQuantityMap := make(map[string]int32)
	for _, r := range records {
		idQuantityMap[r.ID] += r.Quantity
	}
	var (
		ok       = true
		failedOn []struct {
			ID   string
			Want int32
			Have int32
		}
	)
	for _, item := range query {
		if item.Quantity > idQuantityMap[item.ID] {
			ok = false
			failedOn = append(failedOn, struct {
				ID   string
				Want int32
				Have int32
			}{ID: item.ID, Want: item.Quantity, Have: idQuantityMap[item.ID]})
		}
	}
	if ok {
		return h.inventoryRepo.UpdateInventory(ctx, query, func(
			ctx context.Context,
			existing []*entity.ItemWithQuantity,
			query []*entity.ItemWithQuantity,
		) ([]*entity.ItemWithQuantity, error) {
			var newItems []*entity.ItemWithQuantity
			for _, e := range existing {
				for _, q := range query {
					if e.ID == q.ID {
						newItems = append(newItems, &entity.ItemWithQuantity{
							ID:       e.ID,
							Quantity: e.Quantity - q.Quantity,
						})
					}
				}
			}
			return newItems, nil
		})
	}
	return domain.ExceedInventoryError{FailedOn: failedOn}
}

func getStubPriceID(id string) string {
	priceID, ok := stub[id]
	if !ok {
		priceID = stub["1"]
	}
	return priceID
}
