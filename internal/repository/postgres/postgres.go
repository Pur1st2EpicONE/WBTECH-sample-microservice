// Package postgres provides a PostgreSQL storage implementation for the Storage interface.
// It wraps the database connection and logger, and offers methods for managing orders, deliveries,
// payments, and items in a transactional and safe way. It also includes connection management utilities.
package postgres

import (
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/pkg/logger"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// Storage wraps the database connection and logger for interacting with PostgreSQL
type Storage struct {
	db     *sqlx.DB
	logger logger.Logger
}

// NewStorage creates a new Storage instance with the provided database connection and logger
func NewStorage(db *sqlx.DB, logger logger.Logger) *Storage {
	return &Storage{db: db, logger: logger}
}

// Ping checks the database connection to ensure it is alive
func (s *Storage) Ping() error {
	return s.db.Ping()
}

// Close safely closes the database connection and logs the result
func (s *Storage) Close() {
	if err := s.db.Close(); err != nil {
		s.logger.LogError("postgres — failed to close properly", err, "layer", "repository.postgres")
	} else {
		s.logger.LogInfo("postgres — stopped", "layer", "repository.postgres")
	}
}
