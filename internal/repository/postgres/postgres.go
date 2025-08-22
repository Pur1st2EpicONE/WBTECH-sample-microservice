package postgres

import (
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/logger"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type PostgresStorer struct {
	db     *sqlx.DB
	logger logger.Logger
}

func NewPostgresStorer(db *sqlx.DB, logger logger.Logger) *PostgresStorer {
	return &PostgresStorer{db: db, logger: logger}
}

func (p *PostgresStorer) Ping() error {
	return p.db.Ping()
}

func (p *PostgresStorer) Close() {
	if err := p.db.Close(); err != nil {
		p.logger.LogError("postgres — failed to close properly", err)
	} else {
		p.logger.LogInfo("postgres — stopped")
	}
}
