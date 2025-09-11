package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/models"
)

// SaveOrder inserts a complete order with delivery, payment, and items into the database as a single transaction
func (s *Storage) SaveOrder(order *models.Order) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %v", err)
	}
	defer func() { _ = tx.Rollback() }()

	orderId, err := insertOrder(ctx, tx, order)
	if err != nil {
		return fmt.Errorf("failed to insert order: %v", err)
	}
	if err := insertDelivery(ctx, tx, &order.Delivery, orderId); err != nil {
		return fmt.Errorf("failed to insert delivery: %v", err)
	}
	if err := insertPayment(ctx, tx, &order.Payment, orderId); err != nil {
		return fmt.Errorf("failed to insert payment: %v", err)
	}
	for i := range order.Items {
		if err := insertItem(ctx, tx, &order.Items[i], orderId); err != nil {
			return fmt.Errorf("failed to insert item: %v", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("transaction commit failed: %v", err)
	}
	return nil
}

// insertOrder inserts the main order record and returns the generated order ID
func insertOrder(ctx context.Context, tx *sql.Tx, order *models.Order) (int, error) {
	var id int
	query := `
	INSERT INTO orders (
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
		return id, fmt.Errorf("row.Scan failed to get order id: %v", err)
	}
	return id, nil
}

// insertDelivery inserts delivery details associated with the given order ID
func insertDelivery(ctx context.Context, tx *sql.Tx, delivery *models.Delivery, orderID int) error {
	query := `
	INSERT INTO deliveries (
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

// insertPayment inserts payment details associated with the given order ID
func insertPayment(ctx context.Context, tx *sql.Tx, payment *models.Payment, orderID int) error {
	paymentTime := time.Unix(payment.PaymentDT, 0)
	query := `
	INSERT INTO payments (
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

// insertItem inserts a single item associated with the given order ID
func insertItem(ctx context.Context, tx *sql.Tx, item *models.Item, orderID int) error {
	query := `
	INSERT INTO items (
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
