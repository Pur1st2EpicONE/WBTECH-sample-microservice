package cache

import (
	"sync"
	"time"

	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/models"
)

const orderTTL = 24 * time.Hour

type Cache struct {
	mu       sync.RWMutex
	orders   map[string]*CachedOrder
	orderTTL time.Duration
}

func NewCache() *Cache {
	return &Cache{orders: make(map[string]*CachedOrder), orderTTL: orderTTL}
}

func (c *Cache) GetCachedOrder(orderID string) (*models.Order, bool) {
	c.mu.RLock()
	cachedOrder, found := c.orders[orderID]
	c.mu.RUnlock()
	if !found {
		return nil, false
	}
	c.mu.Lock()
	cachedOrder.lastAccess = time.Now()
	c.mu.Unlock()
	return cachedOrder.order, true
}

func (c *Cache) CacheOrder(order *models.Order) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if cachedOrder, found := c.orders[order.OrderUID]; found {
		cachedOrder.order = order // to keep cache updated if the order changes in the database (it shouldn't but just in case)
		cachedOrder.lastAccess = time.Now()
	} else {
		c.orders[order.OrderUID] = &CachedOrder{order: order, lastAccess: time.Now()}
	}

}
