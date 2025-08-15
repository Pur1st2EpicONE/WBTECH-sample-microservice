package main

import (
	"fmt"

	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/kafka"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/logger"
)

func main() {
	producer, err := kafka.NewProducer([]string{"localhost:9092"})
	if err != nil {
		logger.LogError("producer creation failed", err)
	}
	orders := getOrders()
	for i, order := range orders {
		logger.LogInfo(fmt.Sprintf("order-producer — sending order %d to Kafka", i))
		producer.Produce(order, "orders")
	}
}

func getOrders() []string {
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
		`{
			"order_uid": "b77HelloThere77test",
			"track_number": "WBILMGESTTRACK",
			"entry": "WBIL",
			"delivery": {
				"name": "Max Payne",
				"phone": "88005553535",
				"zip": "2639398",
				"city": "Kaliningrad",
				"address": "Pushkina",
				"region": "Kaliningrad",
				"email": "test2@mail.com"
			},
			"payment": {
				"transaction": "b77HelloThere77test",
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
					"rid": "ab4211327a764ae0btest",
					"name": "John",
					"sale": 20,
					"size": "0",
					"total_price": 214,
					"nm_id": 2389421,
					"brand": "Logitech",
					"status": 202
				},
				{
					"chrt_id": 9934931,
					"track_number": "WBILMTESTTRACK",
					"price": 100,
					"rid": "ab4219087a764ae0btest",
					"name": "Aboba",
					"sale": 50,
					"size": "0",
					"total_price": 254,
					"nm_id": 2389212,
					"brand": "SONY",
					"status": 202
				}
			],
			"locale": "ru",
			"internal_signature": "",
			"customer_id": "test",
			"delivery_service": "meest",
			"shardkey": "9",
			"sm_id": 99,
			"date_created": "2021-12-26T06:22:19Z",
			"oof_shard": "1"
		}`,
	}
	return orders
}
