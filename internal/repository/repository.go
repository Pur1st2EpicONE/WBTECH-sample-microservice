package repository

import (
	"fmt"

	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/configs"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/models"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/repository/postgres"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/pkg/logger"
	"github.com/jmoiron/sqlx"
)

// Storage defines methods for interacting with order storage (DB).
type Storage interface {
	SaveOrder(order *models.Order) error
	GetOrder(id string) (*models.Order, error)
	GetOrders(amount ...int) ([]*models.Order, error)
	Ping() error
	Close()
}

// NewStorage wraps a specific storage implementation (Postgres) into the Storage interface.
func NewStorage(db *sqlx.DB, logger logger.Logger) Storage {
	return postgres.NewStorage(db, logger)
}

// ConnectDB establishes a connection to the database using the given configuration.
// Configures connection pool parameters and verifies connectivity with Ping.
func ConnectDB(config configs.Database) (*sqlx.DB, error) {
	db, err := sqlx.Open(config.Driver, fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.Username, config.Password, config.DBName, config.SSLMode))
	if err != nil {
		return nil, fmt.Errorf("database driver not found or DSN invalid: %v", err)
	}
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("database ping failed: %v", err)
	}
	db.SetMaxOpenConns(config.MaxOpenConns)
	db.SetMaxIdleConns(config.MaxIdleConns)
	db.SetConnMaxLifetime(config.ConnMaxLifetime)
	db.SetConnMaxIdleTime(config.ConnMaxIdleTime)
	return db, nil
}
