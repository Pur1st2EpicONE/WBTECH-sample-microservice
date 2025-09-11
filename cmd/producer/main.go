// Package main implements a standalone Kafka producer for testing.
//
// This code is separate from the main service and is used only to generate
// and send test orders.
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/cmd/producer/order"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/broker"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/configs"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/pkg/logger"
)

// Entry point for the Kafka order producer application
// Loads configuration, initializes logger and Kafka producer,
// generates orders (including a bad order for testing), and sends them to Kafka
func main() {
	loggerConfig := configs.Logger{LogDir: "", Debug: false}
	logger, _ := logger.NewLogger(loggerConfig)

	config, err := configs.ProdConfig()
	if err != nil {
		logger.LogFatal("producer — failed to load config", err)
	}

	producer, err := broker.NewProducer(config, logger)
	if err != nil {
		logger.LogFatal("producer — creation failed", err)
	}

	sendBadOrder := checkArgs(&config.MsgsToSend)
	if sendBadOrder {
		badOrder := order.CreateBadOrder(logger)
		logger.LogInfo(fmt.Sprintf("order-producer — sending bad order %s to Kafka", badOrder.OrderUID))
		badOrderJSON, err := json.Marshal(badOrder)
		if err != nil {
			logger.LogFatal("producer — failed to marshal bad order", err)
		}
		badOrderKeyJSON, err := json.Marshal(badOrder.OrderUID)
		if err != nil {
			logger.LogFatal("producer — failed to marshal bad order's key", err)
		}
		badMsg := configs.Message{Topic: config.Topic, Key: badOrderKeyJSON, Value: badOrderJSON}
		_ = producer.Produce(badMsg)
		producer.Close()
		return
	}

	orders := order.GetOrders(config.MsgsToSend, logger)
	for i, order := range orders {
		orderJSON, err := json.Marshal(order)
		if err != nil {
			logger.LogFatal("producer — failed to marshal order", err)
		}
		keyJSON, err := json.Marshal(order.OrderUID)
		if err != nil {
			logger.LogFatal("producer — failed to marshal key", err)
		}
		logger.LogInfo(fmt.Sprintf("order-producer — sending order %s to Kafka", orders[i].OrderUID))
		msg := configs.Message{Topic: config.Topic, Key: keyJSON, Value: orderJSON}
		_ = producer.Produce(msg)
	}
	producer.Close()
}

func checkArgs(amount *int) bool {
	if len(os.Args) > 1 {
		if os.Args[1] == "bad" {
			return true
		}
		newAmount, err := strconv.Atoi(os.Args[1])
		if err != nil {
			*amount = 10
		} else {
			*amount = newAmount
		}
	}
	return false
}
