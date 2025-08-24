package postgres

import (
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/logger"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Storage struct {
	db     *sqlx.DB
	logger logger.Logger
}

func NewStorage(db *sqlx.DB, logger logger.Logger) *Storage {
	return &Storage{db: db, logger: logger}
}

func (s *Storage) Ping() error {
	return s.db.Ping()
}

func (s *Storage) Close() {
	if err := s.db.Close(); err != nil {
		s.logger.LogError("postgres — failed to close properly", err)
	} else {
		s.logger.LogInfo("postgres — stopped")
	}
}
