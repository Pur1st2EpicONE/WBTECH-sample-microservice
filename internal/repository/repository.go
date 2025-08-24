package repository

import (
	"fmt"

	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/configs"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/logger"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/models"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/repository/postgres"
	"github.com/jmoiron/sqlx"
)

//go:generate mockgen -source=repository.go -destination=mocks/mock.go

type Storage interface {
	SaveOrder(order *models.Order) error
	GetOrder(id string) (*models.Order, error)
	GetOrders(amount ...int) ([]*models.Order, error)
	Ping() error
	Close()
}

func NewStorage(db *sqlx.DB, logger logger.Logger) Storage {
	return postgres.NewStorage(db, logger)
}

func ConnectDB(config configs.Database) (*sqlx.DB, error) {
	db, err := sqlx.Open(config.Driver, fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.Username, config.Password, config.DBName, config.SSLMode))
	if err != nil {
		return nil, fmt.Errorf("database driver not found or DSN invalid: %v", err)
	}
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("database ping failed: %v", err)
	}
	return db, nil
}
