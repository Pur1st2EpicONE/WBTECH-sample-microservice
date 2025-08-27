package memory

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/configs"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/models"
	mock_repository "github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/repository/mocks"
	mock_logger "github.com/Pur1st2EpicONE/WBTECH-sample-microservice/pkg/logger/mocks"
	"github.com/golang/mock/gomock"
)

func TestCacheOrder_WithGoMock(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	mockLogger := mock_logger.NewMockLogger(controller)
	mockLogger.EXPECT().LogInfo("order saved", gomock.Any())

	cache := &Cache{
		cachedOrders: make(map[string]*CachedOrder),
		queue:        newQueue(2),
	}

	order := &models.Order{OrderUID: "1"}
	cache.CacheOrder(order, mockLogger)

	if _, ok := cache.cachedOrders["1"]; !ok {
		t.Error("order 1 not cached")
	}
}

func TestNewCache_WithStorageMock(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	storageMock := mock_repository.NewMockStorage(controller)
	mockLogger := mock_logger.NewMockLogger(controller)
	mockLogger.EXPECT().LogInfo("load from database completed", "layer", "cache.memory")
	storageMock.EXPECT().GetOrders(5).Return([]*models.Order{{OrderUID: "1"}, {OrderUID: "2"}}, nil)

	cache := NewCache(storageMock, configs.Cache{
		SaveInCache:   true,
		CacheSize:     5,
		BgCleanup:     false,
		OrderTTL:      time.Minute,
		CleanupPeriod: time.Minute,
	}, mockLogger)

	if len(cache.cachedOrders) != 2 {
		t.Errorf("expected 2 cached orders, got %d", len(cache.cachedOrders))
	}
}

func TestNewCache_DisabledCache(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	storageMock := mock_repository.NewMockStorage(controller)
	mockLogger := mock_logger.NewMockLogger(controller)

	config := configs.Cache{
		SaveInCache: false,
		CacheSize:   0,
	}

	cache := NewCache(storageMock, config, mockLogger)

	if cache == nil || cache.cachedOrders != nil {
		t.Errorf("expected empty cache, got %+v", cache)
	}
}

func TestNewCache_StorageError(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	storageMock := mock_repository.NewMockStorage(controller)
	mockLogger := mock_logger.NewMockLogger(controller)

	storageMock.EXPECT().GetOrders(5).Return(nil, fmt.Errorf("db error"))
	mockLogger.EXPECT().LogError("failed to load orders from database -> %v", gomock.Any(), "layer", "cache.memory")

	config := configs.Cache{
		SaveInCache:   true,
		CacheSize:     5,
		BgCleanup:     false,
		OrderTTL:      time.Minute,
		CleanupPeriod: time.Minute,
	}

	cache := NewCache(storageMock, config, mockLogger)

	if len(cache.cachedOrders) != 0 {
		t.Errorf("expected 0 cached orders on error, got %d", len(cache.cachedOrders))
	}
}

func TestQueue_Enqueue_Overflow(t *testing.T) {
	q := newQueue(2)
	ret1 := q.enqueue("1")
	if ret1 != "1" {
		t.Errorf("expected 1, got %s", ret1)
	}
	ret2 := q.enqueue("2")
	if ret2 != "2" {
		t.Errorf("expected 2, got %s", ret2)
	}
	ret3 := q.enqueue("3")
	if ret3 != "1" {
		t.Errorf("expected tail 1, got %s", ret3)
	}
	if q.buffer[0] != "3" && q.buffer[1] != "3" {
		t.Errorf("expected 3 to be in buffer, got %+v", q.buffer)
	}
}

func TestGetCachedOrder(t *testing.T) {
	cache := &Cache{
		cachedOrders: make(map[string]*CachedOrder),
		queue:        newQueue(10),
	}
	order := &models.Order{OrderUID: "1"}
	cache.cachedOrders["1"] = newCachedOrder(order)
	gotOrder, ok := cache.GetCachedOrder("1")
	if !ok || gotOrder != order {
		t.Errorf("expected to find order, got %+v, %v", gotOrder, ok)
	}
	gotOrder, ok = cache.GetCachedOrder("2")
	if ok || gotOrder != nil {
		t.Errorf("expected not found, got %+v, %v", gotOrder, ok)
	}
	cache.queue = nil
	gotOrder, ok = cache.GetCachedOrder("1")
	if ok || gotOrder != nil {
		t.Errorf("expected not found with nil queue, got %+v, %v", gotOrder, ok)
	}
}

func TestCacheOrder_QueueNil(t *testing.T) {
	cache := &Cache{
		queue:        nil,
		cachedOrders: make(map[string]*CachedOrder),
	}
	mockLogger := mock_logger.NewMockLogger(nil)
	order := &models.Order{OrderUID: "1"}
	cache.CacheOrder(order, mockLogger)
}

func TestCacheOrder_UpdateExisting(t *testing.T) {
	cache := &Cache{
		queue:        newQueue(10),
		cachedOrders: make(map[string]*CachedOrder),
	}
	mockLogger := mock_logger.NewMockLogger(nil)

	order := &models.Order{OrderUID: "1"}
	cached := newCachedOrder(order)
	cache.cachedOrders["1"] = cached

	newOrder := &models.Order{OrderUID: "1"}
	cache.CacheOrder(newOrder, mockLogger)

	if cache.cachedOrders["1"].order != newOrder {
		t.Errorf("expected order to be updated")
	}
	if cache.cachedOrders["1"].lastAccess.Before(cached.lastAccess) {
		t.Errorf("expected lastAccess to be updated")
	}
}

func TestCacheOrder_NewOrder_Overflow(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	mockLogger := mock_logger.NewMockLogger(controller)
	mockStorage := mock_repository.NewMockStorage(controller)

	config := configs.Cache{
		SaveInCache:   true,
		CacheSize:     2,
		BgCleanup:     false,
		CleanupPeriod: time.Second * 1,
		OrderTTL:      time.Second * 5,
	}

	mockStorage.EXPECT().GetOrders(gomock.Any()).Return([]*models.Order{}, nil).AnyTimes()
	mockLogger.EXPECT().LogInfo("load from database completed", "layer", "cache.memory")
	mockLogger.EXPECT().LogInfo("order saved", "orderUID", "1", "layer", "cache.memory")
	mockLogger.EXPECT().LogInfo("order saved", "orderUID", "2", "layer", "cache.memory")
	mockLogger.EXPECT().LogInfo("order saved", "orderUID", "3", "layer", "cache.memory")

	cache := NewCache(mockStorage, config, mockLogger)

	order1 := &models.Order{OrderUID: "1"}
	order2 := &models.Order{OrderUID: "2"}
	order3 := &models.Order{OrderUID: "3"}

	cache.CacheOrder(order1, mockLogger)
	cache.CacheOrder(order2, mockLogger)
	cache.CacheOrder(order3, mockLogger)

	if _, ok := cache.GetCachedOrder("1"); ok {
		t.Errorf("expected order1 to be evicted, but it is still cached")
	}
	if _, ok := cache.GetCachedOrder("2"); !ok {
		t.Errorf("expected order2 to be cached")
	}
	if _, ok := cache.GetCachedOrder("3"); !ok {
		t.Errorf("expected order3 to be cached")
	}
}

func TestCacheCleaner(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	mockLogger := mock_logger.NewMockLogger(controller)

	cache := &Cache{
		bgCleanup:     true,
		cachedOrders:  make(map[string]*CachedOrder),
		queue:         newQueue(10),
		orderTTL:      50 * time.Millisecond,
		cleanupPeriod: 20 * time.Millisecond,
	}

	mockLogger.EXPECT().LogInfo("order saved", "orderUID", "1", "layer", "cache.memory")

	order := &models.Order{OrderUID: "1"}
	cache.CacheOrder(order, mockLogger)

	mockLogger.EXPECT().LogInfo("cleaner started", "layer", "cache.memory")
	mockLogger.EXPECT().LogInfo("cleanup cycle started", "layer", "cache.memory").AnyTimes()
	mockLogger.EXPECT().LogInfo("order deleted", "orderUID", "1", "layer", "cache.memory")
	mockLogger.EXPECT().LogInfo("cleanup cycle completed", "layer", "cache.memory").AnyTimes()
	mockLogger.EXPECT().LogInfo("cleaner stopped", "layer", "cache.memory")

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	go cache.CacheCleaner(ctx, mockLogger)

	time.Sleep(150 * time.Millisecond)

	if _, ok := cache.cachedOrders["1"]; ok {
		t.Error("order1 should have been deleted by CacheCleaner")
	}
}

func TestCacheCleaner_Disabled(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	mockLogger := mock_logger.NewMockLogger(controller)

	cache := &Cache{
		bgCleanup:    false,
		cachedOrders: make(map[string]*CachedOrder),
	}

	cache.CacheCleaner(context.Background(), mockLogger)
}
