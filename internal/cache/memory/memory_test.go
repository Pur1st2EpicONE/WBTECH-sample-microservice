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
	mockLogger.EXPECT().LogInfo("cache — order saved", gomock.Any())

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
	mockLogger.EXPECT().LogInfo("cache — load from database completed", "layer", "cache.memory")
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
	mockLogger.EXPECT().LogError("cache — failed to load orders from database: %v", gomock.Any(), "layer", "cache.memory")

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

	oldAccess := cached.lastAccess.Load()

	newOrder := &models.Order{OrderUID: "1"}
	cache.CacheOrder(newOrder, mockLogger)

	newAccess := cache.cachedOrders["1"].lastAccess.Load()
	if newAccess <= oldAccess {
		t.Errorf("expected lastAccess to be updated, got old=%d new=%d", oldAccess, newAccess)
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
	mockLogger.EXPECT().LogInfo("cache — load from database completed", "layer", "cache.memory")
	mockLogger.EXPECT().LogInfo("cache — order saved", "orderUID", "1", "layer", "cache.memory")
	mockLogger.EXPECT().LogInfo("cache — order saved", "orderUID", "2", "layer", "cache.memory")
	mockLogger.EXPECT().LogInfo("cache — order saved", "orderUID", "3", "layer", "cache.memory")

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

	mockLogger.EXPECT().LogInfo("cache — order saved", "orderUID", "1", "layer", "cache.memory")

	order := &models.Order{OrderUID: "1"}
	cache.CacheOrder(order, mockLogger)

	mockLogger.EXPECT().LogInfo("cache — cleaner started", "layer", "cache.memory")
	mockLogger.EXPECT().Debug("cache — cleanup cycle started", "layer", "cache.memory").AnyTimes()
	mockLogger.EXPECT().Debug("cache — order deleted", "orderUID", "1", "layer", "cache.memory")
	mockLogger.EXPECT().Debug("cache — cleanup cycle completed", "layer", "cache.memory").AnyTimes()
	mockLogger.EXPECT().LogInfo("cache — cleaner stopped", "layer", "cache.memory")

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	go cache.CacheCleaner(ctx, mockLogger, make(chan bool))

	time.Sleep(150 * time.Millisecond)

	if _, ok := cache.cachedOrders["1"]; ok {
		t.Error("order 1 should have been deleted by CacheCleaner")
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

	cache.CacheCleaner(context.Background(), mockLogger, make(chan bool))
}

func TestCacheCleaner_DBDown(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	logger := mock_logger.NewMockLogger(controller)
	logger.EXPECT().LogInfo(gomock.Any(), gomock.Any()).AnyTimes()
	logger.EXPECT().Debug(gomock.Any(), gomock.Any()).AnyTimes()

	cache := &Cache{
		bgCleanup:     true,
		cleanupPeriod: 50 * time.Millisecond,
		cachedOrders: map[string]*CachedOrder{
			"1": newCachedOrder(&models.Order{OrderUID: "1"}),
		},
	}

	ctx := t.Context()

	dbStatus := make(chan bool, 1)

	go cache.CacheCleaner(ctx, logger, dbStatus)

	dbStatus <- false
	time.Sleep(1 * time.Second)
	dbStatus <- true
	time.Sleep(1 * time.Second)

	if cache.pauseCleaner {
		t.Errorf("expected pauseCleaner to be false after DB is restored")
	}
}
