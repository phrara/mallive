package adapters

import (
	"context"
	"strconv"
	"sync"
	"time"

	domain "github.com/phrara/mallive/order/domain/order"
	"github.com/sirupsen/logrus"
)

var _ domain.Repository = (*MemoryOrderRepository)(nil)
type MemoryOrderRepository struct {
	mu *sync.RWMutex
	store []*domain.Order
}

func (m *MemoryOrderRepository) Create(_ context.Context, order *domain.Order) (*domain.Order, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	res := &domain.Order{
		ID:     strconv.FormatInt(time.Now().Unix(), 10),
		CustomerID:  order.CustomerID,
		Status:      order.Status,
		PaymentLink: order.PaymentLink,
		Items:       order.Items,
	}
	m.store = append(m.store, res)
	logrus.WithFields(logrus.Fields{
		"input_order": order,
		"store_res": m.store,
	}).Debug("mem_order_repo_create")
	return res, nil
}

func (m *MemoryOrderRepository) Get(_ context.Context, orderID, customerID string) (*domain.Order, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, o := range m.store {
		if o.ID == orderID && o.CustomerID == customerID {
			logrus.Debugf("mem_order_repo_get || found: orderID=%v, customerID=%v, res=%v", orderID, customerID, *o)
			return o, nil
		}
	}
	return nil, domain.NotFoundError{
		OrderID: orderID,
	}
}

func (m *MemoryOrderRepository) Update(ctx context.Context, order *domain.Order, updateFunc func(context.Context, *domain.Order) (*domain.Order, error)) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for i, o := range m.store {
		if o.ID == order.ID && o.CustomerID == order.CustomerID {
			new_o, err := updateFunc(ctx, o)
			if err != nil {
				return err
			}
			m.store[i] = new_o
			return nil
		}
	}

	return domain.NotFoundError{
		OrderID: order.ID,
	}
}


func NewMemoryOrderRepository() *MemoryOrderRepository {
	return &MemoryOrderRepository{
		mu: &sync.RWMutex{},
		store: make([]*domain.Order, 0),
	}
}

