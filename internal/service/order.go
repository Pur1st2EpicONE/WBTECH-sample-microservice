package service

import (
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/models"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/pkg/logger"
)

func (s Service) GetOrder(orderID string, logger logger.Logger) (*models.Order, bool, error) {
	if order, found := s.Cache.GetCachedOrder(orderID); found {
		return order, true, nil
	}
	order, err := s.Storage.GetOrder(orderID)
	if err != nil {
		return nil, false, err
	}
	s.Cache.CacheOrder(order, logger)
	return order, false, nil
}
