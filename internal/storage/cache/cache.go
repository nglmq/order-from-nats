package cache

import (
	"github.com/nglmq/wildberries-0/internal/models"
	"sync"
)

type Cache struct {
	rw     sync.RWMutex
	orders map[string]models.Order
}

func NewCache() *Cache {
	return &Cache{
		orders: make(map[string]models.Order),
	}
}

func (c *Cache) GetFromCache(orderID string) (models.Order, bool) {
	c.rw.RLock()
	defer c.rw.RUnlock()

	order, exists := c.orders[orderID]

	return order, exists
}

func (c *Cache) SaveToCache(orderID string, order models.Order) {
	c.rw.Lock()
	defer c.rw.Unlock()

	c.orders[orderID] = order
}
