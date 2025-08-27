package postgres_test

import (
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/models"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/repository/postgres"
	mock_logger "github.com/Pur1st2EpicONE/WBTECH-sample-microservice/pkg/logger/mocks"
	"github.com/golang/mock/gomock"
	"github.com/jmoiron/sqlx"
)

func TestPostgresStorer_SaveOrder_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %v", err)
	}
	defer db.Close()

	logger := mock_logger.NewMockLogger(gomock.NewController(t))
	s := postgres.NewStorage(sqlx.NewDb(db, "postgres"), logger)

	order := &models.Order{
		Items: []models.Item{
			{
				ChrtID:      1,
				TrackNumber: "test123",
				Price:       100,
				Rid:         "test123",
				Name:        "Test Testov",
				Sale:        0,
				Size:        "BIG_TEST",
				TotalPrice:  100,
				NmID:        1,
				Brand:       "cool_test_brand",
				Status:      1,
			},
		},
	}

	mock.ExpectBegin()

	mock.ExpectQuery("INSERT INTO orders").
		WithArgs(
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
		).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	mock.ExpectExec("INSERT INTO deliveries").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO payments").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO items").WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	err = s.SaveOrder(order)
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unfulfilled expectations: %v", err)
	}
}

func TestPostgresStorer_SaveOrder_BeginTxError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %v", err)
	}
	defer db.Close()

	logger := mock_logger.NewMockLogger(gomock.NewController(t))
	s := postgres.NewStorage(sqlx.NewDb(db, "postgres"), logger)

	mock.ExpectBegin().WillReturnError(fmt.Errorf("begin failed"))

	err = s.SaveOrder(new(models.Order))
	if err == nil {
		t.Fatalf("expected begin error, got %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unfulfilled expectations: %v", err)
	}
}

func TestPostgresStorer_InsertOrder_Rollback(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %v", err)
	}
	defer db.Close()

	logger := mock_logger.NewMockLogger(gomock.NewController(t))
	s := postgres.NewStorage(sqlx.NewDb(db, "postgres"), logger)

	mock.ExpectBegin()

	mock.ExpectQuery("INSERT INTO orders").
		WithArgs(
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
		).
		WillReturnError(fmt.Errorf("failed to insert order"))

	mock.ExpectRollback()

	err = s.SaveOrder(new(models.Order))
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unfulfilled expectations: %v", err)
	}
}

func TestPostgresStorer_InsertDelivery_Rollback(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %v", err)
	}
	defer db.Close()

	logger := mock_logger.NewMockLogger(gomock.NewController(t))
	s := postgres.NewStorage(sqlx.NewDb(db, "postgres"), logger)

	mock.ExpectBegin()

	mock.ExpectQuery("INSERT INTO orders").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectExec("INSERT INTO deliveries").
		WithArgs(
			1,
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
		).
		WillReturnError(fmt.Errorf("failed to insert delivery"))

	mock.ExpectRollback()

	err = s.SaveOrder(new(models.Order))
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unfulfilled expectations: %v", err)
	}
}

func TestPostgresStorer_InsertPayment_Rollback(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %v", err)
	}
	defer db.Close()

	logger := mock_logger.NewMockLogger(gomock.NewController(t))
	s := postgres.NewStorage(sqlx.NewDb(db, "postgres"), logger)

	mock.ExpectBegin()

	mock.ExpectQuery("INSERT INTO orders").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectExec("INSERT INTO deliveries").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO payments").WithArgs(
		sqlmock.AnyArg(),
		sqlmock.AnyArg(),
		sqlmock.AnyArg(),
		sqlmock.AnyArg(),
		sqlmock.AnyArg(),
		sqlmock.AnyArg(),
		sqlmock.AnyArg(),
		sqlmock.AnyArg(),
		sqlmock.AnyArg(),
		sqlmock.AnyArg(),
		sqlmock.AnyArg(),
	).WillReturnError(fmt.Errorf("failed to insert payment"))

	mock.ExpectRollback()

	err = s.SaveOrder(new(models.Order))
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unfulfilled expectations: %v", err)
	}
}

func TestPostgresStorer_InsertItem_Rollback(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %v", err)
	}
	defer db.Close()

	logger := mock_logger.NewMockLogger(gomock.NewController(t))
	s := postgres.NewStorage(sqlx.NewDb(db, "postgres"), logger)

	order := &models.Order{
		Items: []models.Item{
			{
				ChrtID:      1,
				TrackNumber: "test123",
				Price:       100,
				Rid:         "test123",
				Name:        "Test Testov",
				Sale:        0,
				Size:        "BIG_TEST",
				TotalPrice:  100,
				NmID:        1,
				Brand:       "cool_test_brand",
				Status:      1,
			},
		},
	}

	mock.ExpectBegin()

	mock.ExpectQuery("INSERT INTO orders").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectExec("INSERT INTO deliveries").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO payments").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO items").
		WithArgs(
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
		).
		WillReturnError(fmt.Errorf("failed to insert item"))

	mock.ExpectRollback()

	err = s.SaveOrder(order)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unfulfilled expectations: %v", err)
	}
}

func TestPostgresStorer_SaveOrder_CommitError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %v", err)
	}
	defer db.Close()

	logger := mock_logger.NewMockLogger(gomock.NewController(t))
	s := postgres.NewStorage(sqlx.NewDb(db, "postgres"), logger)

	order := &models.Order{
		Items: []models.Item{
			{
				ChrtID:      1,
				TrackNumber: "test123",
				Price:       100,
				Rid:         "test123",
				Name:        "Test Testov",
				Sale:        0,
				Size:        "BIG_TEST",
				TotalPrice:  100,
				NmID:        1,
				Brand:       "cool_test_brand",
				Status:      1,
			},
		},
	}

	mock.ExpectBegin()

	mock.ExpectQuery("INSERT INTO orders").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectExec("INSERT INTO deliveries").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO payments").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO items").WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit().WillReturnError(fmt.Errorf("commit failed"))

	err = s.SaveOrder(order)
	if err == nil {
		t.Fatalf("expected commit error, got %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unfulfilled expectations: %v", err)
	}
}
