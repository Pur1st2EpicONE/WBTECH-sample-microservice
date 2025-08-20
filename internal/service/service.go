package service

import (
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/cache"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/models"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/repository"
)

type ServiceProvider interface {
	GetOrder(orderID string) (*models.Order, bool, error)
}

type Service struct {
	Storage repository.Storer
	Cache   cache.Cache
}

func NewService(storage repository.Storer, cache cache.Cache) *Service {
	return &Service{Storage: storage, Cache: cache}
}
