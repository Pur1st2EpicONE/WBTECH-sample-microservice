package postgres

import (
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/logger"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type PostgresStorage struct {
	db     *sqlx.DB
	logger logger.Logger
}

func NewPostgresStorage(db *sqlx.DB, logger logger.Logger) *PostgresStorage {
	return &PostgresStorage{db: db, logger: logger}
}

func (p *PostgresStorage) Ping() error {
	return p.db.Ping()
}

func (p *PostgresStorage) Close() {
	if err := p.db.Close(); err != nil {
		p.logger.LogError("postgres — failed to close properly", err)
	} else {
		p.logger.LogInfo("postgres — stopped")
	}
}
