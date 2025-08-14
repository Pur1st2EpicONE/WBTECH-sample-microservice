package repository

import (
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/models"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/repository/postgres"
	"github.com/jmoiron/sqlx"
)

type Storer interface {
	SaveOrder(order *models.Order) error
	GetOrder(id string) (*models.Order, error)
	Ping() error
	Close() error
}

type Storage struct {
	Storer
}

func NewStorage(db *sqlx.DB) *Storage {
	return &Storage{Storer: postgres.NewPostgresStorer(db)}
}
