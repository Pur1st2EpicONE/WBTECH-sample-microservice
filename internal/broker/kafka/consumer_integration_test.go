package kafka_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/cmd/producer/order"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/broker"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/configs"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/models"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/repository"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/repository/postgres"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/pkg/logger"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

func TestKafkaConsumer_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	cfg := configs.App{
		Database: configs.Database{
			Driver:   "postgres",
			Host:     "localhost",
			Port:     "5434",
			Username: "Neo",
			Password: "0451",
			DBName:   "wb-service-db-test",
			SSLMode:  "disable",
		},
		Consumer: configs.Consumer{
			Brokers:             []string{"localhost:9092"},
			Topic:               "test-orders",
			ClientID:            "test-client",
			GroupID:             "test-group",
			SaveOrderRetryDelay: 100 * time.Millisecond,
			SaveOrderRetryMax:   2,
			CommitRetryDelay:    100 * time.Millisecond,
			CommitRetryMax:      2,
			DLQ: configs.Producer{
				Brokers:  []string{"localhost:9092"},
				Topic:    "test-dlq",
				ClientID: "dlq-client",
				Kafka: &configs.KafkaProducer{
					Retries:           1,
					CompressionType:   "none",
					EnableIdempotence: false,
					BatchSize:         1, Acks: "all"},
			},
			Kafka: &configs.Kafka{
				EnableAutoCommit: false,
				AutoOffsetReset:  "earliest",
			},
		},
	}

	log, _ := logger.NewLogger(cfg.Logger)
	db, err := repository.ConnectDB(cfg.Database)
	if err != nil {
		t.Fatalf("failed to connect to test DB: %v", err)
	}
	storage := postgres.NewStorage(db, log)
	time.Sleep(5 * time.Second)

	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
	})
	if err != nil {
		t.Fatalf("failed to create producer: %v", err)
	}
	defer producer.Close()

	topic := cfg.Consumer.Topic
	order := order.CreateOrder(log)
	orderJSON, _ := json.Marshal(order)
	msg := &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Key:            []byte("1"),
		Value:          orderJSON,
	}

	if err := producer.Produce(msg, nil); err != nil {
		t.Fatalf("failed to produce message: %v", err)
	}

	if err := producer.Produce(msg, nil); err != nil {
		t.Fatalf("failed to produce message for the second time: %v", err)
	}

	badMsg := bad()
	badkey, _ := json.Marshal("b563feb7b2b84b6test")
	badmsg, _ := json.Marshal(badMsg)
	badmessage := &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Key:            badkey,
		Value:          badmsg,
	}
	if err := producer.Produce(badmessage, nil); err != nil {
		t.Fatalf("failed to produce message: %v", err)
	}

	evenWorse := worse(log)
	worseKey, _ := json.Marshal("b563feb7b2b84b6test")
	worseMsg, _ := json.Marshal(evenWorse)
	badermessage := &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Key:            worseKey,
		Value:          worseMsg,
	}
	if err := producer.Produce(badermessage, nil); err != nil {
		t.Fatalf("failed to produce message: %v", err)
	}

	kc, err := broker.NewConsumer(cfg.Consumer, log)
	if err != nil {
		t.Fatalf("failed to create consumer: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	go kc.Run(ctx, storage, log, 1)
	producer.Flush(7000)

	time.Sleep(10 * time.Second)

	_, err = storage.GetOrder(order.OrderUID)
	if err != nil {
		t.Fatalf("failed to fetch order: %v", err)
	}

	kc.Close(log)

}

func bad() []string {
	orders := []string{"aboba"}
	return orders
}

func worse(log logger.Logger) models.Order {
	bo := order.CreateOrder(log)
	bo.Delivery.Email = "E-male"
	return bo
}
