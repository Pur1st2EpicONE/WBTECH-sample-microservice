package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	model "github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/models"
)

func (ps *PostgresStorer) SaveOrder(order *model.Order) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	tx, err := ps.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return err
	}
	orderID, err := insertOrder(ctx, tx, order)
	if err != nil {
		tx.Rollback()
		return err
	}
	if err := insertDelivery(ctx, tx, order.Delivery, orderID); err != nil {
		tx.Rollback()
		return err
	}
	if err := insertPayment(ctx, tx, order.Payment, orderID); err != nil {
		tx.Rollback()
		return err
	}
	for i := range order.Items {
		if err := insertItem(ctx, tx, &order.Items[i], orderID); err != nil {
			tx.Rollback()
			return err
		}
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit failed: %w", err)
	}
	return nil
}

func insertOrder(ctx context.Context, tx *sql.Tx, order *model.Order) (int, error) {
	var id int
	query := `
	INSERT INTO order_table (
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
	) 
	VALUES (
		$1, 
		$2, 
		$3, 
		$4, 
		$5, 
		$6, 
		$7, 
		$8, 
		$9, 
		$10, 
		$11
	) 
	RETURNING id`

	row := tx.QueryRowContext(
		ctx,
		query,
		order.OrderUID,
		order.TrackNumber,
		order.Entry,
		order.Locale,
		order.InternalSignature,
		order.CustomerID,
		order.DeliveryService,
		order.ShardKey,
		order.SmID,
		order.DateCreated,
		order.OofShard)
	if err := row.Scan(&id); err != nil {
		return id, err
	}
	return id, nil
}

func insertDelivery(ctx context.Context, tx *sql.Tx, delivery model.Delivery, orderID int) error {
	query := `
	INSERT INTO delivery (
		order_id,
		name,
		phone, 
		zip,
		city,
		address,
		region, 
		email
	) 
	VALUES (
		$1, 
		$2, 
		$3, 
		$4, 
		$5, 
		$6, 
		$7, 
		$8
	)`

	_, err := tx.ExecContext(
		ctx,
		query,
		orderID,
		delivery.Name,
		delivery.Phone,
		delivery.Zip,
		delivery.City,
		delivery.Address,
		delivery.Region,
		delivery.Email)

	return err
}

func insertPayment(ctx context.Context, tx *sql.Tx, payment model.Payment, orderID int) error {
	paymentTime := time.Unix(payment.PaymentDT, 0)
	query := `
	INSERT INTO payment (
		order_id,
		transaction, 
		request_id, 
		currency, 
		provider, 
		amount, 
		payment_dt, 
		bank,
		delivery_cost,
		goods_total,
		custom_fee
	) 
	VALUES (
		$1, 
		$2, 
		$3, 
		$4, 
		$5, 
		$6, 
		$7, 
		$8,
		$9,
		$10,
		$11 
	)`

	_, err := tx.ExecContext(
		ctx,
		query,
		orderID,
		payment.Transaction,
		payment.RequestID,
		payment.Currency,
		payment.Provider,
		payment.Amount,
		paymentTime,
		payment.Bank,
		payment.DeliveryCost,
		payment.GoodsTotal,
		payment.CustomFee)

	return err
}

func insertItem(ctx context.Context, tx *sql.Tx, item *model.Item, orderID int) error {
	query := `
	INSERT INTO item (
		order_id,
		chrt_id,
    	track_number,
    	price,
    	rid,
    	name,
    	sale,
    	size,
    	total_price,
    	nm_id,
    	brand,
    	status
	) 
	VALUES (
		$1, 
		$2, 
		$3, 
		$4, 
		$5, 
		$6, 
		$7, 
		$8,
		$9,
		$10,
		$11,
		$12
	)`

	_, err := tx.ExecContext(
		ctx,
		query,
		orderID,
		item.ChrtID,
		item.TrackNumber,
		item.Price,
		item.Rid,
		item.Name,
		item.Sale,
		item.Size,
		item.TotalPrice,
		item.NmID,
		item.Brand,
		item.Status,
	)
	return err
}
