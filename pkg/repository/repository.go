package repository

import (
	model "github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/models"
	"github.com/jmoiron/sqlx"
)

type Storer interface {
	SaveOrder(order *model.Order) error
}

type Storage struct {
	Storer
}

func NewStorage(db *sqlx.DB) *Storage {
	return &Storage{Storer: NewPostgresStorer(db)}
}
