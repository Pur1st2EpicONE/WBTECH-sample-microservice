package postgres_test

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/repository/postgres"
	mock_logger "github.com/Pur1st2EpicONE/WBTECH-sample-microservice/pkg/logger/mocks"
	"github.com/golang/mock/gomock"
	"github.com/jmoiron/sqlx"
)

func TestPostgresStorer_Ping(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	xdb := sqlx.NewDb(db, "sqlmock")
	logger := mock_logger.NewMockLogger(gomock.NewController(t))
	storer := postgres.NewStorage(xdb, logger)

	mock.ExpectPing()
	if err := storer.Ping(); err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	pingErr := errors.New("ping failed")
	mock.ExpectPing().WillReturnError(pingErr)
	if err := storer.Ping(); !errors.Is(err, pingErr) {
		t.Errorf("expected pingErr, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %v", err)
	}
}

func TestPostgresStorer_Close_Success(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	mockLogger := mock_logger.NewMockLogger(controller)

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	mock.ExpectClose()

	xdb := sqlx.NewDb(db, "sqlmock")
	storer := postgres.NewStorage(xdb, mockLogger)

	mockLogger.EXPECT().LogInfo("postgres — stopped", "layer", "repository.postgres").Times(1)

	storer.Close()

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled db expectations: %v", err)
	}
}
