package kafka

import (
	"encoding/json"
	"fmt"

	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/logger"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/models"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/repository"
)

type PgHandler struct {
}

func NewOrderHandler() *PgHandler {
	return &PgHandler{}
}

func (h *PgHandler) SaveOrder(jsonMsg []byte, db repository.Storage) error {
	order := new(models.Order)
	if err := json.Unmarshal(jsonMsg, order); err != nil {
		return fmt.Errorf("pg-handler — failed to unmarshal the order: %v", err)
	}
	if err := db.Ping(); err != nil {
		return fmt.Errorf("pg-handler — lost connection to database: %v", err)
	}
	if err := db.SaveOrder(order); err != nil {
		return fmt.Errorf("pg-handler — failed to save order %s to database: %v", order.OrderUID, err)
	}
	logger.LogInfo("pg-handler — saved order", "orderUID", order.OrderUID)
	return nil
}
