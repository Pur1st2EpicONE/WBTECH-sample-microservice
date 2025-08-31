// Package cache provides an abstraction for caching orders and managing cache cleanup.
//
// It defines a Cache interface for storing and retrieving orders,
// as well as running a background cache cleaner that can adapt
// to database connectivity status.
//
// The cache implementation ensures fast access to frequently used orders
// and works together with the storage layer to maintain consistency.
//
// NewCache returns a concrete in-memory cache implementation.
package cache

import (
	"context"

	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/cache/memory"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/configs"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/models"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/repository"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/pkg/logger"
)

// Cache defines the behavior for caching orders and managing background cleanup.
type Cache interface {
	GetCachedOrder(orderID string) (*models.Order, bool)
	CacheOrder(order *models.Order, logger logger.Logger)
	CacheCleaner(ctx context.Context, logger logger.Logger, dbStatus chan bool)
}

// NewCache creates a new in-memory cache instance, wired to the storage and logger.
func NewCache(storage repository.Storage, config configs.Cache, logger logger.Logger) Cache {
	return memory.NewCache(storage, config, logger)
}
