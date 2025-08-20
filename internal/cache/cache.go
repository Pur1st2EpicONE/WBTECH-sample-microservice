package cache

import (
	"context"
	"sync"
	"time"

	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/logger"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/models"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/repository"
)

//go:generate mockgen -source=cache.go -destination=mocks/mock.go

type Cache interface {
	GetCachedOrder(orderID string) (*models.Order, bool)
	CacheOrder(order *models.Order)
	CacheCleaner(ctx context.Context)
}

type defaultCache struct {
	mu       sync.RWMutex
	orders   map[string]*CachedOrder
	orderTTL time.Duration
}

type CachedOrder struct {
	order      *models.Order
	lastAccess time.Time
}

func NewCache(storage repository.Storer, orderTTL time.Duration) Cache {
	cachedOrders := make(map[string]*CachedOrder)

	allOrders, err := storage.GetAllOrders()
	if err != nil {
		logger.LogError("cache — failed to load orders from database", err)
	} else {
		for _, order := range allOrders {
			cachedOrders[order.OrderUID] = &CachedOrder{
				order:      order,
				lastAccess: time.Now(),
			}
		}
		logger.LogInfo("cache — load from database complete")
	}

	return &defaultCache{
		orders:   cachedOrders,
		orderTTL: orderTTL,
	}
}

func (c *defaultCache) GetCachedOrder(orderID string) (*models.Order, bool) {
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

func (c *defaultCache) CacheOrder(order *models.Order) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if cachedOrder, found := c.orders[order.OrderUID]; found {
		cachedOrder.order = order // to keep cache updated if the order changes in the database (it shouldn't but just in case)
		cachedOrder.lastAccess = time.Now()
	} else {
		c.orders[order.OrderUID] = &CachedOrder{order: order, lastAccess: time.Now()}
		logger.LogInfo("cache — saved order", "orderUID", order.OrderUID)
	}
}

func (c *defaultCache) CacheCleaner(ctx context.Context) {
	logger.LogInfo("cache — cleaner started")
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			logger.LogInfo("cache — cleaner stopped")
			return
		case <-ticker.C:
			logger.LogInfo("cache — starting cleanup cycle")
			var expiredOrders []string
			c.mu.RLock()
			for orderUID, cached := range c.orders {
				if time.Since(cached.lastAccess) > c.orderTTL {
					expiredOrders = append(expiredOrders, orderUID)
				}
			}
			c.mu.RUnlock()
			if len(expiredOrders) > 0 {
				c.mu.Lock()
				for _, orderUID := range expiredOrders {
					delete(c.orders, orderUID)
					logger.LogInfo("cache — deleted order", "orderUID", orderUID)
				}
				c.mu.Unlock()
			}
			logger.LogInfo("cache — cache cleaned")
		}
	}
}
