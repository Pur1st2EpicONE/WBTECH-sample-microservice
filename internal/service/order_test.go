package service

import (
	"fmt"
	"testing"

	mock_cache "github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/cache/mocks"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/models"
	mock_repo "github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/repository/mocks"
	"github.com/golang/mock/gomock"
)

func TestNewService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorer := mock_repo.NewMockStorer(ctrl)
	mockCacher := mock_cache.NewMockCacher(ctrl)

	service := NewService(mockStorer, mockCacher)
	if service == nil {
		t.Fatal("expected service, got nil")
	}
}

func TestService_GetOrder_CacheHit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorer := mock_repo.NewMockStorer(ctrl)
	mockCacher := mock_cache.NewMockCacher(ctrl)

	service := &Service{
		Storage: mockStorer,
		Cache:   mockCacher,
	}

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
		t.Fatalf("expected order with orderID %s, got %s", orderID, order.OrderUID)
	}
}

func TestService_GetOrder_CacheMiss(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorer := mock_repo.NewMockStorer(ctrl)
	mockCacher := mock_cache.NewMockCacher(ctrl)

	service := &Service{
		Storage: mockStorer,
		Cache:   mockCacher,
	}

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
		t.Fatalf("expected order with orderID %s, got %s", orderID, order.OrderUID)
	}
}

func TestService_GetOrder_FromDB_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorer := mock_repo.NewMockStorer(ctrl)
	mockCacher := mock_cache.NewMockCacher(ctrl)

	service := &Service{
		Storage: mockStorer,
		Cache:   mockCacher,
	}

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
