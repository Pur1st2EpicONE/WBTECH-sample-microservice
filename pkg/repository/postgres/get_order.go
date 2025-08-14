package postgres

import (
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/models"
)

func (ps *PostgresStorer) GetOrder(orderUID string) (*models.Order, error) {
	order := new(models.Order)
	if err := queryOrder(ps, order, orderUID); err != nil {
		return nil, err
	}

	return order, nil
}

func queryOrder(ps *PostgresStorer, order *models.Order, orderUID string) error {
	query := `SELECT 
		order_uid, 
		track_number, 
		entry, 
		locale, 
		internal_signature,
		customer_id,
		delivery_service,
		shardkey,
		sm_id,
		date_created,
		oof_shard
		FROM orders WHERE order_uid=$1`
	row := ps.db.QueryRow(query, orderUID)
	if err := row.Scan(
		&order.OrderUID,
		&order.TrackNumber,
		&order.Entry,
		&order.Locale,
		&order.InternalSignature,
		&order.CustomerID,
		&order.DeliveryService,
		&order.ShardKey,
		&order.SmID,
		&order.DateCreated,
		&order.OofShard); err != nil {
		return err
	}
	return nil
}
