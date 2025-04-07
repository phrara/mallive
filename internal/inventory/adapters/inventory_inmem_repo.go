package adapters

import (
	"context"
	"sync"

	"github.com/phrara/mallive/common/genproto/orderpb"
	domain "github.com/phrara/mallive/inventory/domain/inventory"
)



var _ domain.Repository = (*MemoryInventoryRepository)(nil)
type MemoryInventoryRepository struct {
	mu *sync.RWMutex
	store map[string]*orderpb.Item
}

func (m *MemoryInventoryRepository) GetItems(ctx context.Context, itemIDs []string) (res []*orderpb.Item, err error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	lacks := make([]string, 0) 
	var lacksF bool 
	for _ , id := range itemIDs {
		if item, b := m.store[id]; b {
			res = append(res, item)
		} else {
			lacks = append(lacks, id)
			lacksF = true
		}
	}
	if lacksF {
		return res, domain.NotFoundError{
			Lacks: lacks,
		}
	} else {
		return res, nil
	}
}



func NewMemoryInventoryRepository() *MemoryInventoryRepository {
	return &MemoryInventoryRepository{
		mu:    &sync.RWMutex{},
		store: map[string]*orderpb.Item{},
	}
}




