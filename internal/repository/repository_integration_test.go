//go:build integration
// +build integration

package repository_test

import (
	"testing"

	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/broker/kafka"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/repository"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/repository/postgres"
	"github.com/jmoiron/sqlx"
)

func TestPostgresStorer_SaveOrder_Integration(t *testing.T) {
	db, err := sqlx.Connect("postgres", "postgres://Neo:0451@localhost:5434/wb-service-db-test?sslmode=disable")
	if err != nil {
		t.Fatalf("failed to connect to db: %v", err)
	}
	defer db.Close()

	ps := repository.NewStorage(db).Storer.(*postgres.PostgresStorer)

	order := kafka.CreateOrder()
	order.OrderUID = "1"
	order.Payment.Transaction = order.OrderUID

	if err := ps.SaveOrder(&order); err != nil {
		t.Fatalf("SaveOrder failed: %v", err)
	}
}

func TestPostgresStorer_GetOrder_Integration(t *testing.T) {
	db, err := sqlx.Connect("postgres", "postgres://Neo:0451@localhost:5434/wb-service-db-test?sslmode=disable")
	if err != nil {
		t.Fatalf("failed to connect to db: %v", err)
	}
	defer db.Close()

	ps := repository.NewStorage(db).Storer.(*postgres.PostgresStorer)

	gotOrder, err := ps.GetOrder("1")
	if err != nil {
		t.Fatalf("GetOrder failed: %v", err)
	}

	if gotOrder.OrderUID != "1" {
		t.Fatalf("expected orderUID 1, got %s", gotOrder.OrderUID)
	}
}

func TestPostgresStorer_SaveAndGetOrder_Integration(t *testing.T) {
	db, err := sqlx.Connect("postgres", "postgres://Neo:0451@localhost:5434/wb-service-db-test?sslmode=disable")
	if err != nil {
		t.Fatalf("failed to connect to db: %v", err)
	}
	defer db.Close()

	ps := repository.NewStorage(db).Storer.(*postgres.PostgresStorer)

	order := kafka.CreateOrder()

	if err := ps.SaveOrder(&order); err != nil {
		t.Fatalf("SaveOrder failed: %v", err)
	}

	gotOrder, err := ps.GetOrder(order.OrderUID)
	if err != nil {
		t.Fatalf("GetOrder failed: %v", err)
	}

	if gotOrder.OrderUID != order.OrderUID {
		t.Fatalf("expected orderUID %s, got %s", order.OrderUID, gotOrder.OrderUID)
	}

	if len(gotOrder.Items) != len(order.Items) {
		t.Fatalf("expected %d items, got %d", len(order.Items), len(gotOrder.Items))
	}
}
