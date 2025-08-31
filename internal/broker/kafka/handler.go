package kafka

import (
	"encoding/json"
	"fmt"

	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/models"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/repository"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/pkg/logger"
	"github.com/go-playground/validator/v10"
)

// MessageHandler defines the contract for processing Kafka messages.
// Each message is expected to represent an order in JSON format.
type MessageHandler interface {
	SaveOrder(jsonMsg []byte, storage repository.Storage, logger logger.Logger, workerID int) error
}

// Handler is a concrete implementation of MessageHandler.
// It provides logic for parsing, validating, and storing incoming Kafka messages.
type Handler struct{}

func newHandler() *Handler {
	return new(Handler)
}

// SaveOrder parses a JSON message into an Order, validates it,
// and persists it into the provided storage.
//
// Steps:
//  1. Unmarshal JSON into a models.Order struct.
//  2. Validate the struct fields using go-playground/validator.
//  3. Save the validated order to the storage.
//  4. Log a debug message on success.
//
// If unmarshaling, validation, or saving fails, an error is returned.
// The workerID is included in logs for easier debugging in multi-worker setups.
func (h *Handler) SaveOrder(jsonMsg []byte, storage repository.Storage, logger logger.Logger, workerID int) error {
	validate := validator.New()
	order := new(models.Order)
	if err := json.Unmarshal(jsonMsg, order); err != nil {
		return fmt.Errorf("failed to unmarshal the order: %w", err)
	}
	if err := validate.Struct(order); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}
	if err := storage.SaveOrder(order); err != nil {
		return fmt.Errorf("failed to save order %s to database: %w", order.OrderUID, err)
	}
	logger.Debug(fmt.Sprintf("worker %d — saved order to DB", workerID), "orderUID", order.OrderUID, "workerID", fmt.Sprintf("%d", workerID), "layer", "broker.kafka")
	return nil
}
