package repository_test

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/configs"
	mock_logger "github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/logger/mocks"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/models"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/repository"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/repository/postgres"
	"github.com/golang/mock/gomock"
	"github.com/jmoiron/sqlx"
)

func TestPostgresStorer_SaveOrder_Success(t *testing.T) { // integration?
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %v", err)
	}
	defer db.Close()
	logger := mock_logger.NewMockLogger(gomock.NewController(t))
	ps := postgres.NewPostgresStorage(sqlx.NewDb(db, "postgres"), logger)

	order := &models.Order{
		OrderUID: "123",
		Items:    []models.Item{{ChrtID: 1, Name: "Cool hat", Price: 100}},
		Delivery: models.Delivery{Name: "Aboba"},
		Payment:  models.Payment{Transaction: "1", Amount: 100},
	}

	mock.ExpectBegin()
	mock.ExpectQuery("INSERT INTO orders").
		WithArgs(
			order.OrderUID,
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	mock.ExpectExec("INSERT INTO deliveries").
		WithArgs(
			1,
			order.Delivery.Name,
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectExec("INSERT INTO payments").
		WithArgs(
			1,
			order.Payment.Transaction,
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			order.Payment.Amount,
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectExec("INSERT INTO items").
		WithArgs(
			1,
			order.Items[0].ChrtID,
			sqlmock.AnyArg(),
			order.Items[0].Price,
			sqlmock.AnyArg(),
			order.Items[0].Name,
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	if err := ps.SaveOrder(order); err != nil {
		t.Fatalf("SaveOrder failed: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %v", err)
	}
}

func TestConnectDB_SQLXOpenFail(t *testing.T) {
	config := configs.Database{
		Driver:   "bad driver name",
		Host:     "localhost",
		Port:     "5434",
		Username: "Neo",
		Password: "0452",
		DBName:   "wb-service-db-test",
		SSLMode:  "disable",
	}

	_, err := repository.ConnectDB(config)
	if err == nil {
		t.Fatalf("ConnectPostgres fail to fail: %v", err)
	}
}

func TestConnectDB_PingFail(t *testing.T) {
	cfg := configs.Database{
		Driver:   "postgres",
		Host:     "localhost",
		Port:     "bad port",
		Username: "user",
		Password: "pass",
		DBName:   "test",
		SSLMode:  "disable",
	}

	_, err := repository.ConnectDB(cfg)
	if err == nil {
		t.Fatal("expected ping to fail, got nil")
	}
}
