package main

import (
	"fmt"

	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/kafka"
)

func main() {
	producer, err := kafka.NewProducer([]string{"localhost:9092"})
	if err != nil {
		fmt.Println(err)
	}
	orders := getOrders()
	for _, order := range orders {
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
			"order_uid": "b563feb7b2b84b7test",
			"track_number": "WBILMTESTTRACK2",
			"entry": "WBIL",
			"delivery": {
				"name": "Alice Smith",
				"phone": "+9721111111",
				"zip": "2639810",
				"city": "Kiryat Mozkin",
				"address": "Main St 12",
				"region": "Kraiot",
				"email": "alice@gmail.com"
			},
			"payment": {
				"transaction": "b563feb7b2b84b7test",
				"request_id": "",
				"currency": "USD",
				"provider": "wbpay",
				"amount": 1200,
				"payment_dt": 1637908727,
				"bank": "alpha",
				"delivery_cost": 500,
				"goods_total": 700,
				"custom_fee": 0
			},
			"items": [
				{
					"chrt_id": 9934931,
					"track_number": "WBILMTESTTRACK2",
					"price": 700,
					"rid": "ab4219087a764ae1btest",
					"name": "Lipstick",
					"sale": 0,
					"size": "0",
					"total_price": 700,
					"nm_id": 2389213,
					"brand": "L’Oreal",
					"status": 202
				}
			],
			"locale": "en",
			"internal_signature": "",
			"customer_id": "alice",
			"delivery_service": "dhl",
			"shardkey": "2",
			"sm_id": 50,
			"date_created": "2021-11-27T10:12:00Z",
			"oof_shard": "1"
		}`,
		`{
			"order_uid": "b563feb7b2b84b8test",
			"track_number": "WBILMTESTTRACK3",
			"entry": "WBIL",
			"delivery": {
				"name": "Bob Johnson",
				"phone": "+9722222222",
				"zip": "2639811",
				"city": "Kiryat Mozkin",
				"address": "Park Ave 5",
				"region": "Kraiot",
				"email": "bob@gmail.com"
			},
			"payment": {
				"transaction": "b563feb7b2b84b8test",
				"request_id": "",
				"currency": "USD",
				"provider": "wbpay",
				"amount": 2300,
				"payment_dt": 1637910000,
				"bank": "beta",
				"delivery_cost": 300,
				"goods_total": 2000,
				"custom_fee": 0
			},
			"items": [
				{
					"chrt_id": 9934932,
					"track_number": "WBILMTESTTRACK3",
					"price": 2000,
					"rid": "ab4219087a764ae2ctest",
					"name": "Foundation",
					"sale": 10,
					"size": "0",
					"total_price": 1800,
					"nm_id": 2389214,
					"brand": "Maybelline",
					"status": 202
				}
			],
			"locale": "en",
			"internal_signature": "",
			"customer_id": "bob",
			"delivery_service": "ups",
			"shardkey": "5",
			"sm_id": 60,
			"date_created": "2021-11-28T14:00:00Z",
			"oof_shard": "2"
		}`,
		`{
			"order_uid": "b563feb7b2b84b9test",
			"track_number": "WBILMTESTTRACK4",
			"entry": "WBIL",
			"delivery": {
				"name": "Charlie Brown",
				"phone": "+9723333333",
				"zip": "2639812",
				"city": "Kiryat Mozkin",
				"address": "Central St 8",
				"region": "Kraiot",
				"email": "charlie@gmail.com"
			},
			"payment": {
				"transaction": "b563feb7b2b84b9test",
				"request_id": "",
				"currency": "USD",
				"provider": "wbpay",
				"amount": 1500,
				"payment_dt": 1637911000,
				"bank": "gamma",
				"delivery_cost": 200,
				"goods_total": 1300,
				"custom_fee": 0
			},
			"items": [
				{
					"chrt_id": 9934933,
					"track_number": "WBILMTESTTRACK4",
					"price": 1300,
					"rid": "ab4219087a764ae3dtest",
					"name": "Eyeliner",
					"sale": 0,
					"size": "0",
					"total_price": 1300,
					"nm_id": 2389215,
					"brand": "Rimmel",
					"status": 202
				}
			],
			"locale": "en",
			"internal_signature": "",
			"customer_id": "charlie",
			"delivery_service": "fedex",
			"shardkey": "7",
			"sm_id": 70,
			"date_created": "2021-11-29T09:30:00Z",
			"oof_shard": "2"
		}`,
		`{
			"order_uid": "b563feb7b2b84b10test",
			"track_number": "WBILMTESTTRACK5",
			"entry": "WBIL",
			"delivery": {
				"name": "Diana Prince",
				"phone": "+9724444444",
				"zip": "2639813",
				"city": "Kiryat Mozkin",
				"address": "Liberty St 20",
				"region": "Kraiot",
				"email": "diana@gmail.com"
			},
			"payment": {
				"transaction": "b563feb7b2b84b10test",
				"request_id": "",
				"currency": "USD",
				"provider": "wbpay",
				"amount": 3200,
				"payment_dt": 1637912000,
				"bank": "delta",
				"delivery_cost": 400,
				"goods_total": 2800,
				"custom_fee": 0
			},
			"items": [
				{
					"chrt_id": 9934934,
					"track_number": "WBILMTESTTRACK5",
					"price": 2800,
					"rid": "ab4219087a764ae4etest",
					"name": "Blush",
					"sale": 10,
					"size": "0",
					"total_price": 2520,
					"nm_id": 2389216,
					"brand": "Sephora",
					"status": 202
				}
			],
			"locale": "en",
			"internal_signature": "",
			"customer_id": "diana",
			"delivery_service": "dhl",
			"shardkey": "3",
			"sm_id": 80,
			"date_created": "2021-11-30T12:00:00Z",
			"oof_shard": "3"
		}`,
		`{
			"order_uid": "b563feb7b2b84b11test",
			"track_number": "WBILMTESTTRACK6",
			"entry": "WBIL",
			"delivery": {
				"name": "Eve Adams",
				"phone": "+9725555555",
				"zip": "2639814",
				"city": "Kiryat Mozkin",
				"address": "Green St 4",
				"region": "Kraiot",
				"email": "eve@gmail.com"
			},
			"payment": {
				"transaction": "b563feb7b2b84b11test",
				"request_id": "",
				"currency": "USD",
				"provider": "wbpay",
				"amount": 1450,
				"payment_dt": 1637913000,
				"bank": "omega",
				"delivery_cost": 200,
				"goods_total": 1250,
				"custom_fee": 0
			},
			"items": [
				{"chrt_id": 9934935,"track_number":"WBILMTESTTRACK6","price":500,"rid":"rid001","name":"Mascara","sale":10,"size":"0","total_price":450,"nm_id":2389217,"brand":"L’Oreal","status":202},
				{"chrt_id": 9934936,"track_number":"WBILMTESTTRACK6","price":800,"rid":"rid002","name":"Lipstick","sale":20,"size":"0","total_price":640,"nm_id":2389218,"brand":"Maybelline","status":202},
				{"chrt_id": 9934937,"track_number":"WBILMTESTTRACK6","price":100,"rid":"rid003","name":"Brush","sale":0,"size":"0","total_price":100,"nm_id":2389219,"brand":"Sephora","status":202}
			],
			"locale":"en","internal_signature":"","customer_id":"eve","delivery_service":"ups","shardkey":"1","sm_id":101,"date_created":"2021-12-01T08:00:00Z","oof_shard":"1"
		}`,
		`{
			"order_uid": "b563feb7b2b84b12test",
			"track_number": "WBILMTESTTRACK7",
			"entry": "WBIL",
			"delivery": {
				"name": "Frank Miller",
				"phone": "+9726666666",
				"zip": "2639815",
				"city": "Kiryat Mozkin",
				"address": "Blue St 7",
				"region": "Kraiot",
				"email": "frank@gmail.com"
			},
			"payment": {
				"transaction": "b563feb7b2b84b12test",
				"request_id": "",
				"currency": "USD",
				"provider": "wbpay",
				"amount": 2100,
				"payment_dt": 1637914000,
				"bank": "alpha",
				"delivery_cost": 150,
				"goods_total": 1950,
				"custom_fee": 0
			},
			"items": [
				{"chrt_id": 9934938,"track_number":"WBILMTESTTRACK7","price":1200,"rid":"rid004","name":"Perfume","sale":0,"size":"0","total_price":1200,"nm_id":2389220,"brand":"Chanel","status":202},
				{"chrt_id": 9934939,"track_number":"WBILMTESTTRACK7","price":800,"rid":"rid005","name":"Eyeshadow","sale":10,"size":"0","total_price":720,"nm_id":2389221,"brand":"MAC","status":202},
				{"chrt_id": 9934940,"track_number":"WBILMTESTTRACK7","price":50,"rid":"rid006","name":"Eyeliner","sale":0,"size":"0","total_price":50,"nm_id":2389222,"brand":"Sephora","status":202}
			],
			"locale":"en","internal_signature":"","customer_id":"frank","delivery_service":"dhl","shardkey":"3","sm_id":102,"date_created":"2021-12-02T09:30:00Z","oof_shard":"1"
		}`,
		`{
			"order_uid": "b563feb7b2b84b13test",
			"track_number": "WBILMTESTTRACK8",
			"entry": "WBIL",
			"delivery": {
				"name": "Grace Lee",
				"phone": "+9727777777",
				"zip": "2639816",
				"city": "Kiryat Mozkin",
				"address": "Red St 10",
				"region": "Kraiot",
				"email": "grace@gmail.com"
			},
			"payment": {
				"transaction": "b563feb7b2b84b13test",
				"request_id": "",
				"currency": "USD",
				"provider": "wbpay",
				"amount": 3050,
				"payment_dt": 1637915000,
				"bank": "beta",
				"delivery_cost": 300,
				"goods_total": 2750,
				"custom_fee": 0
			},
			"items": [
				{"chrt_id": 9934941,"track_number":"WBILMTESTTRACK8","price":1500,"rid":"rid007","name":"Foundation","sale":0,"size":"0","total_price":1500,"nm_id":2389223,"brand":"Maybelline","status":202},
				{"chrt_id": 9934942,"track_number":"WBILMTESTTRACK8","price":1000,"rid":"rid008","name":"Blush","sale":10,"size":"0","total_price":900,"nm_id":2389224,"brand":"Sephora","status":202},
				{"chrt_id": 9934943,"track_number":"WBILMTESTTRACK8","price":250,"rid":"rid009","name":"Mascara","sale":0,"size":"0","total_price":250,"nm_id":2389225,"brand":"L’Oreal","status":202},
				{"chrt_id": 9934944,"track_number":"WBILMTESTTRACK8","price":200,"rid":"rid010","name":"Eyeliner","sale":0,"size":"0","total_price":200,"nm_id":2389226,"brand":"MAC","status":202}
			],
			"locale":"en","internal_signature":"","customer_id":"grace","delivery_service":"ups","shardkey":"2","sm_id":103,"date_created":"2021-12-03T11:00:00Z","oof_shard":"2"
		}`,
		`{
			"order_uid": "b563feb7b2b84b14test",
			"track_number": "WBILMTESTTRACK9",
			"entry": "WBIL",
			"delivery": {
				"name": "Henry King",
				"phone": "+9728888888",
				"zip": "2639817",
				"city": "Kiryat Mozkin",
				"address": "Yellow St 3",
				"region": "Kraiot",
				"email": "henry@gmail.com"
			},
			"payment": {
				"transaction": "b563feb7b2b84b14test",
				"request_id": "",
				"currency": "USD",
				"provider": "wbpay",
				"amount": 1700,
				"payment_dt": 1637916000,
				"bank": "gamma",
				"delivery_cost": 150,
				"goods_total": 1550,
				"custom_fee": 0
			},
			"items": [
				{"chrt_id": 9934945,"track_number":"WBILMTESTTRACK9","price":800,"rid":"rid011","name":"Lipstick","sale":0,"size":"0","total_price":800,"nm_id":2389227,"brand":"MAC","status":202},
				{"chrt_id": 9934946,"track_number":"WBILMTESTTRACK9","price":750,"rid":"rid012","name":"Blush","sale":0,"size":"0","total_price":750,"nm_id":2389228,"brand":"Sephora","status":202}
			],
			"locale":"en","internal_signature":"","customer_id":"henry","delivery_service":"dhl","shardkey":"4","sm_id":104,"date_created":"2021-12-04T13:15:00Z","oof_shard":"2"
		}`,
		`{
			"order_uid": "b563feb7b2b84b15test",
			"track_number": "WBILMTESTTRACK10",
			"entry": "WBIL",
			"delivery": {
				"name": "Ivy Watson",
				"phone": "+9729999999",
				"zip": "2639818",
				"city": "Kiryat Mozkin",
				"address": "White St 12",
				"region": "Kraiot",
				"email": "ivy@gmail.com"
			},
			"payment": {
				"transaction": "b563feb7b2b84b15test",
				"request_id": "",
				"currency": "USD",
				"provider": "wbpay",
				"amount": 4000,
				"payment_dt": 1637917000,
				"bank": "delta",
				"delivery_cost": 300,
				"goods_total": 3700,
				"custom_fee": 0
			},
			"items": [
				{"chrt_id": 9934947,"track_number":"WBILMTESTTRACK10","price":1000,"rid":"rid013","name":"Foundation","sale":0,"size":"0","total_price":1000,"nm_id":2389229,"brand":"Maybelline","status":202},
				{"chrt_id": 9934948,"track_number":"WBILMTESTTRACK10","price":1200,"rid":"rid014","name":"Mascara","sale":10,"size":"0","total_price":1080,"nm_id":2389230,"brand":"L’Oreal","status":202},
				{"chrt_id": 9934949,"track_number":"WBILMTESTTRACK10","price":1500,"rid":"rid015","name":"Lipstick","sale":5,"size":"0","total_price":1425,"nm_id":2389231,"brand":"Sephora","status":202},
				{"chrt_id": 9934950,"track_number":"WBILMTESTTRACK10","price":200,"rid":"rid016","name":"Brush","sale":0,"size":"0","total_price":200,"nm_id":2389232,"brand":"MAC","status":202}
			],
			"locale":"en","internal_signature":"","customer_id":"ivy","delivery_service":"ups","shardkey":"5","sm_id":105,"date_created":"2021-12-05T15:45:00Z","oof_shard":"3"
		}`,
	}
	return orders
}
