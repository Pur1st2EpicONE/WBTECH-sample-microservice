// Package service provides business logic for managing orders.
// It interacts with the storage layer and cache to retrieve and store orders efficiently.
package service

import (
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/cache"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/models"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/repository"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/pkg/logger"
)

// ServiceProvider defines the interface for service operations on orders.
type ServiceProvider interface {
	// GetOrder retrieves an order by its ID.
	// Returns the order, a boolean indicating if it was retrieved from cache, and an error if any.
	GetOrder(orderID string, logger logger.Logger) (*models.Order, bool, error)
}

// Service implements ServiceProvider using a storage backend and cache.
type Service struct {
	Storage repository.Storage
	Cache   cache.Cache
}

// NewService creates a new Service instance with the provided storage and cache.
func NewService(storage repository.Storage, cache cache.Cache) Service {
	return Service{Storage: storage, Cache: cache}
}
