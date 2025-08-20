package postgres

import (
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/logger"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type PostgresStorer struct {
	db *sqlx.DB
}

func NewPostgresStorer(db *sqlx.DB) *PostgresStorer {
	return &PostgresStorer{db: db}
}

func (p *PostgresStorer) Ping() error {
	return p.db.Ping()
}

func (p *PostgresStorer) Close() {
	if err := p.db.Close(); err != nil {
		logger.LogError("postgres — failed to close properly", err)
	} else {
		logger.LogInfo("postgres — stopped")
	}
}
