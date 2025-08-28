package memory

import (
	"context"
	"sync"
	"sync/atomic"
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
	pauseCleaner  bool
	pauseDuration time.Duration
}

func NewCache(storage repository.Storage, config configs.Cache, logger logger.Logger) *Cache {
	if !config.SaveInCache || config.CacheSize < 1 {
		return new(Cache)
	}

	var queue *Queue
	cachedOrders := make(map[string]*CachedOrder, config.CacheSize)
	queue = newQueue(config.CacheSize)

	allOrders, err := storage.GetOrders(config.CacheSize)
	if err != nil {
		logger.LogError("cache — failed to load orders from database: %v", err, "layer", "cache.memory")
	} else {
		for _, order := range allOrders {
			cachedOrders[order.OrderUID] = newCachedOrder(order)
			queue.enqueue(order.OrderUID)
		}
		logger.LogInfo("cache — load from database completed", "layer", "cache.memory")
	}

	return &Cache{
		bgCleanup:     config.BgCleanup,
		cachedOrders:  cachedOrders,
		orderTTL:      config.OrderTTL,
		queue:         queue,
		cleanupPeriod: config.CleanupPeriod,
		pauseDuration: config.PauseDuration,
	}
}

type CachedOrder struct {
	order      *models.Order
	lastAccess atomic.Int64
}

func newCachedOrder(order *models.Order) *CachedOrder {
	cachedOrder := &CachedOrder{order: order}
	cachedOrder.lastAccess.Store(time.Now().UnixNano())
	return cachedOrder
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
	cachedOrder.lastAccess.Store(time.Now().UnixNano())
	return cachedOrder.order, true
}

func (c *Cache) CacheOrder(order *models.Order, logger logger.Logger) {
	if c.queue == nil {
		return
	}
	c.mu.Lock()
	if cachedOrder, found := c.cachedOrders[order.OrderUID]; found {
		cachedOrder.lastAccess.Store(time.Now().UnixNano())
	} else {
		rewriteId := c.queue.enqueue(order.OrderUID)
		if rewriteId != order.OrderUID {
			delete(c.cachedOrders, rewriteId)
		}
		c.cachedOrders[order.OrderUID] = newCachedOrder(order)
		logger.LogInfo("cache — order saved", "orderUID", order.OrderUID, "layer", "cache.memory")
	}
	c.mu.Unlock()
}

func (c *Cache) CacheCleaner(ctx context.Context, logger logger.Logger, dbStatus chan bool) {
	if !c.bgCleanup {
		return
	}
	ticker := time.NewTicker(c.cleanupPeriod)
	defer ticker.Stop()
	logger.LogInfo("cache — cleaner started", "layer", "cache.memory")
	for {
		if c.pauseCleaner {
			time.Sleep(c.pauseDuration)
		}
		select {
		case <-ctx.Done():
			logger.LogInfo("cache — cleaner stopped", "layer", "cache.memory")
			return
		case connected := <-dbStatus:
			if connected {
				if c.pauseCleaner {
					logger.LogInfo("cache — connection to database restored, cleaner resumed", "layer", "cache.memory")
				}
				c.pauseCleaner = false
			} else {
				if !c.pauseCleaner {
					logger.LogInfo("cache — lost connection to database, cleaner paused", "layer", "cache.memory")
				}
				c.pauseCleaner = true
				continue
			}
		case <-ticker.C:
			if c.pauseCleaner {
				continue
			}
			logger.Debug("cache — cleanup cycle started", "layer", "cache.memory")
			var expiredOrders []string
			c.mu.RLock()
			for orderUID, order := range c.cachedOrders {
				lastAccess := time.Unix(0, order.lastAccess.Load())
				if time.Since(lastAccess) > c.orderTTL {
					expiredOrders = append(expiredOrders, orderUID)
				}
			}
			c.mu.RUnlock()
			if len(expiredOrders) > 0 {
				c.mu.Lock()
				for _, orderUID := range expiredOrders {
					delete(c.cachedOrders, orderUID)
					logger.Debug("cache — order deleted", "orderUID", orderUID, "layer", "cache.memory")
				}
				c.mu.Unlock()
			}
			logger.Debug("cache — cleanup cycle completed", "layer", "cache.memory")
		}
	}
}
