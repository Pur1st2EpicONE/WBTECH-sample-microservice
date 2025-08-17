package postgres

import (
	"fmt"

	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/configs"
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

func ConnectPostgres(config configs.Database) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.Username, config.Password, config.DBName, config.SSLMode))
	if err != nil {
		return nil, fmt.Errorf("sqlx.Open failed to open database: %v", err)
	}
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("lost connection to database: %v", err)
	}
	logger.LogInfo("postgres — connected to database")
	return db, nil
}

func (p *PostgresStorer) Ping() error {
	return p.db.Ping()
}

func (p *PostgresStorer) Close() {
	logger.LogInfo("postgres — stopped")
	if err := p.db.Close(); err != nil {
		logger.LogError("postgres — failed to close properly", err)
	}
}
