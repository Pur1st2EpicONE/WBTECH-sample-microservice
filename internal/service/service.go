package service

import (
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/cache"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/models"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/repository"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type ServiceProvider interface {
	GetOrder(orderID string) (*models.Order, bool, error)
}

type Service struct {
	repository.Storage
	cache.Cache
}

func NewService(storage *repository.Storage, cache *cache.Cache) *Service {
	return &Service{Storage: *storage, Cache: *cache}
}
