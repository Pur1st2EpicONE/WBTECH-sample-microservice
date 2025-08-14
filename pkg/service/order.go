package service

import (
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/models"
)

func (s *Service) GetOrder(orderID string) (*models.Order, error) {
	if order, found := s.cache.GetCachedOrder(orderID); found {
		return order, nil
	}
	order, err := s.Storage.GetOrder(orderID)
	if err != nil {
		return nil, err
	}
	s.cache.CacheOrder(order)
	return order, nil
}
