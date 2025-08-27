package memory

import (
	"context"
	"sync"
	"time"

	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/configs"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/models"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/repository"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/pkg/logger"
)

type Cache struct {
	bgCleanup     bool
	mu            sync.RWMutex
	cachedOrders  map[string]*CachedOrder
	orderTTL      time.Duration
	queue         *Queue
	cleanupPeriod time.Duration
}

func NewCache(storage repository.Storage, config configs.Cache, logger logger.Logger) *Cache {
	if !config.SaveInCache || config.CacheSize < 1 {
		return new(Cache) // this is clunky but it's too late for that
	}

	var queue *Queue
	cachedOrders := make(map[string]*CachedOrder, config.CacheSize)
	queue = newQueue(config.CacheSize)

	allOrders, err := storage.GetOrders(config.CacheSize)
	if err != nil {
		logger.LogError("failed to load orders from database -> %v", err, "layer", "cache.memory")
	} else {
		for _, order := range allOrders {
			cachedOrders[order.OrderUID] = newCachedOrder(order)
			queue.enqueue(order.OrderUID)
		}
		logger.LogInfo("load from database completed", "layer", "cache.memory")
	}

	return &Cache{
		bgCleanup:     config.BgCleanup,
		cachedOrders:  cachedOrders,
		orderTTL:      config.OrderTTL,
		queue:         queue,
		cleanupPeriod: config.CleanupPeriod,
	}
}

type CachedOrder struct {
	order      *models.Order
	lastAccess time.Time
}

func newCachedOrder(order *models.Order) *CachedOrder {
	return &CachedOrder{order: order, lastAccess: time.Now()}
}

type Queue struct {
	buffer []string
	head   int
	tail   int
	size   int
}

func newQueue(size int) *Queue {
	return &Queue{buffer: make([]string, size)}
}

func (q *Queue) enqueue(orderUID string) string {
	var tail string
	if q.size == len(q.buffer) {
		tail = q.buffer[q.tail]
		q.buffer[q.tail] = orderUID
		q.tail = q.moveIndex(q.tail)
		q.head = q.moveIndex(q.head)
		return tail
	}
	q.buffer[q.tail] = orderUID
	q.tail = q.moveIndex(q.tail)
	q.size++
	return orderUID
}

func (q *Queue) moveIndex(i int) int {
	i++
	if i == len(q.buffer) {
		i = 0
	}
	return i
}

func (c *Cache) GetCachedOrder(orderID string) (*models.Order, bool) {
	if c.queue == nil {
		return nil, false
	}
	c.mu.RLock()
	cachedOrder, found := c.cachedOrders[orderID]
	c.mu.RUnlock()
	if !found {
		return nil, false
	}
	c.mu.Lock()
	cachedOrder.lastAccess = time.Now()
	c.mu.Unlock()
	return cachedOrder.order, true
}

func (c *Cache) CacheOrder(order *models.Order, logger logger.Logger) {
	if c.queue == nil {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if cachedOrder, found := c.cachedOrders[order.OrderUID]; found {
		cachedOrder.order = order // to keep cache updated if the order changes in the database (it shouldn't but just in case)
		cachedOrder.lastAccess = time.Now()
	} else {
		rewriteId := c.queue.enqueue(order.OrderUID)
		if rewriteId != order.OrderUID {
			delete(c.cachedOrders, rewriteId)
		}
		c.cachedOrders[order.OrderUID] = newCachedOrder(order)
		logger.LogInfo("order saved", "orderUID", order.OrderUID, "layer", "cache.memory")
	}
}

func (c *Cache) CacheCleaner(ctx context.Context, logger logger.Logger) {
	if !c.bgCleanup {
		return
	}
	logger.LogInfo("cleaner started", "layer", "cache.memory")
	ticker := time.NewTicker(c.cleanupPeriod)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			logger.LogInfo("cleaner stopped", "layer", "cache.memory")
			return
		case <-ticker.C:
			logger.LogInfo("cleanup cycle started", "layer", "cache.memory")
			var expiredOrders []string
			c.mu.RLock()
			for orderUID, order := range c.cachedOrders {
				if time.Since(order.lastAccess) > c.orderTTL {
					expiredOrders = append(expiredOrders, orderUID)
				}
			}
			c.mu.RUnlock()
			if len(expiredOrders) > 0 {
				c.mu.Lock()
				for _, orderUID := range expiredOrders {
					delete(c.cachedOrders, orderUID)
					logger.LogInfo("order deleted", "orderUID", orderUID, "layer", "cache.memory")
				}
				c.mu.Unlock()
			}
			logger.LogInfo("cleanup cycle completed", "layer", "cache.memory")
		}
	}
}
