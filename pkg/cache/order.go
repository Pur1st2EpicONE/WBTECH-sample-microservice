package cache

import (
	"time"

	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/models"
)

type CachedOrder struct {
	order      *models.Order
	lastAccess time.Time
}
