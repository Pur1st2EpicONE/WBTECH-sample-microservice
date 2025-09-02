// This is a standalone Kafka producer for testing purposes.
// Everything inside this package is not part of the main service and
// is just quick-and-dirty code to generate and send test orders.
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
		logger.LogInfo("order-producer — sending bad order to Kafka — b563feb7b2b84b6test")
		producer.Produce(sendBad())
		producer.Close()
		return
	}

	orders := order.GetOrders(config.MsgsToSend, logger)
	for i, order := range orders {
		orderJSON, err := json.MarshalIndent(order, "", "   ")
		if err != nil {
			logger.LogFatal("producer — failed to marshal order with indent", err)
		}
		keyJSON, err := json.Marshal(order.OrderUID)
		if err != nil {
			logger.LogFatal("producer — failed to marshal key", err)
		}
		logger.LogInfo(fmt.Sprintf("order-producer — sending order %s to Kafka", orders[i].OrderUID))
		msg := configs.Message{Topic: config.Topic, Key: keyJSON, Value: orderJSON}
		producer.Produce(msg)
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

func sendBad() configs.Message {
	val, _ := json.Marshal(bad())
	key, _ := json.Marshal("b563feb7b2b84b6test")
	return configs.Message{Topic: "orders", Key: key, Value: val}
}

func bad() []string {
	orders := []string{
		`{
			"order_uid": "b563feb7b2b84b6test",
			"track_number": "WBILMTESTTRACK",
			"entry": "WBIL",
			"delivery": {
				"name": "Test Testov",
				"phone": "+9720000000",
				"zip": "2639809",
				"city": "Kiryat Mozkin",
				"address": "Ploshad Mira 15",
				"region": "Kraiot",
				"email": "test@gmail.com"
			},
			"payment": {
				"transaction": "b563feb7b2b84b6test",
				"request_id": "",
				"currency": "USD",
				"provider": "wbpay",
				"amount": 1817,
				"payment_dt": 1637907727,
				"bank": "alpha",
				"delivery_cost": 1500,
				"goods_total": 317,
				"custom_fee": 0
			},
			"items": [
				{
					"chrt_id": 9934930,
					"track_number": "WBILMTESTTRACK",
					"price": 453,
					"rid": "ab4219087a764ae0btest",
					"name": "Mascaras",
					"sale": 30,
					"size": "0",
					"total_price": 317,
					"nm_id": 2389212,
					"brand": "Vivienne Sabo",
					"status": 202
				},
				{
					"chrt_id": 9934931,
					"track_number": "WBILMTESTTRACK",
					"price": 100,
					"rid": "ab4219087a764ae0btest",
					"name": "Arnold",
					"sale": 50,
					"size": "0",
					"total_price": 254,
					"nm_id": 2389212,
					"brand": "Amogus",
					"status": 202
				}
			],
			"locale": "en",
			"internal_signature": "",
			"customer_id": "test",
			"delivery_service": "meest",
			"shardkey": "9",
			"sm_id": 99,
			"date_created": "2021-11-26T06:22:19Z",
			"oof_shard": "1"
		}`,
	}
	return orders
}
