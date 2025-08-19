package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/broker"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/broker/kafka"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/configs"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/logger"
)

func main() {

	config, err := configs.ProdConfig()
	if err != nil {
		logger.LogFatal("producer — failed to load config", err)
	}

	producer, err := broker.NewProducer(config)
	if err != nil {
		logger.LogFatal("producer — creation failed", err)
	}

	checkArgs(&config.TotalMessages)
	orders := kafka.GetOrders(config.TotalMessages)

	for i, order := range orders {
		orderJSON, err := json.MarshalIndent(order, "", "   ")
		if err != nil {
			logger.LogFatal("producer — failed to marshal order with indent", err)
		}
		logger.LogInfo(fmt.Sprintf("order-producer — sending order %d to Kafka", i+1))
		producer.Produce(orderJSON, config.Topic)
	}
}

func checkArgs(amount *int) {
	if len(os.Args) > 1 {
		newAmount, err := strconv.Atoi(os.Args[1])
		if err != nil {
			logger.LogError("producer — failed to convert argument to string", err)
			*amount = 10
		} else {
			*amount = newAmount
		}
	}
}
