package cache

import (
	"context"

	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/cache/memory"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/configs"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/logger"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/models"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/repository"
)

//go:generate mockgen -source=cache.go -destination=mocks/mock.go

type Cache interface {
	GetCachedOrder(orderID string) (*models.Order, bool)
	CacheOrder(order *models.Order, logger logger.Logger)
	CacheCleaner(ctx context.Context, logger logger.Logger)
}

func NewCache(storage repository.Storage, config configs.Cache, logger logger.Logger) Cache {
	return memory.NewCache(storage, config, logger)
}
