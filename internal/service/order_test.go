package service

import (
	"fmt"
	"testing"

	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/cache"
	mock_cache "github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/cache/mocks"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/models"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/repository"
	mock_repo "github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/repository/mocks"
	"github.com/golang/mock/gomock"
)

func TestNewService(t *testing.T) {
	mockStorer := new(repository.Storage)
	mockCacher := new(cache.Cache)

	service := NewService(mockStorer, mockCacher)
	if service.Storage != *mockStorer {
		t.Error("storage assigning error")
	}
	if service.Cache != *mockCacher {
		t.Error("cache assigning error")
	}
}

func TestService_GetOrder_CacheHit(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	mockStorer := mock_repo.NewMockStorer(controller)
	mockCacher := mock_cache.NewMockCacher(controller)

	mockStorage := repository.Storage{Storer: mockStorer}
	mockCache := cache.Cache{Cacher: mockCacher}

	service := &Service{Storage: mockStorage, Cache: mockCache}

	orderID := "1703"
	cachedOrder := &models.Order{OrderUID: orderID}

	mockCacher.EXPECT().GetCachedOrder(orderID).Return(cachedOrder, true)

	order, found, err := service.GetOrder(orderID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !found {
		t.Fatalf("order not found in cache")
	}
	if order.OrderUID != orderID {
		t.Fatalf("expected order with orderID %s, got order with orderID %s", orderID, order.OrderUID)
	}
}

func TestService_GetOrder_CacheMiss(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	mockStorer := mock_repo.NewMockStorer(controller)
	mockCacher := mock_cache.NewMockCacher(controller)

	mockStorage := repository.Storage{Storer: mockStorer}
	mockCache := cache.Cache{Cacher: mockCacher}

	service := &Service{Storage: mockStorage, Cache: mockCache}

	orderID := "1"
	expectedOrder := &models.Order{OrderUID: orderID}

	mockCacher.EXPECT().GetCachedOrder(orderID).Return(nil, false)
	mockStorer.EXPECT().GetOrder(orderID).Return(expectedOrder, nil)
	mockCacher.EXPECT().CacheOrder(expectedOrder)

	order, found, err := service.GetOrder(orderID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if found {
		t.Fatalf("cache returned an order when it should be empty")
	}
	if order.OrderUID != orderID {
		t.Fatalf("expected order with orderID %s, got order with orderID %s", orderID, order.OrderUID)
	}
}

func TestService_GetOrder_FromDB_NotFound(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	mockStorer := mock_repo.NewMockStorer(controller)
	mockCacher := mock_cache.NewMockCacher(controller)

	mockStorage := repository.Storage{Storer: mockStorer}
	mockCache := cache.Cache{Cacher: mockCacher}

	service := &Service{Storage: mockStorage, Cache: mockCache}

	orderID := "0"

	mockCacher.EXPECT().GetCachedOrder(orderID).Return(nil, false)
	mockStorer.EXPECT().GetOrder(orderID).Return(nil, fmt.Errorf("order not found in storage"))

	order, fromCache, err := service.GetOrder(orderID)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if order != nil {
		t.Fatalf("expected nil order, got %v", order)
	}
	if fromCache {
		t.Fatalf("expected fromCache value to be false, got true")
	}
}
