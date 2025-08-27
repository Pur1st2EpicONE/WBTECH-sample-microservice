package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/models"
)

func (s *Storage) GetOrder(orderUID string) (*models.Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var orderId int
	order := new(models.Order)
	if err := queryAllButItems(ctx, s, order, orderUID, &orderId); err != nil {
		return nil, err
	}
	if err := queryItems(ctx, s, &order.Items, orderId); err != nil {
		return nil, err
	}
	return order, nil
}

func queryAllButItems(ctx context.Context, s *Storage, order *models.Order, orderUID string, orderId *int) error {
	query := `SELECT 

        orders.id, 
        orders.order_uid, 
        orders.track_number, 
        orders.entry, 
        orders.locale, 
        orders.internal_signature,
        orders.customer_id,
        orders.delivery_service,
        orders.shardkey,
        orders.sm_id,
        orders.date_created,
        orders.oof_shard,

        deliveries.name, 
        deliveries.phone, 
        deliveries.zip, 
        deliveries.city, 
        deliveries.address,
        deliveries.region,
        deliveries.email,

        payments.transaction, 
        payments.request_id, 
        payments.currency, 
        payments.provider, 
        payments.amount,
        payments.payment_dt,
        payments.bank,
        payments.delivery_cost,
        payments.goods_total,
        payments.custom_fee

        FROM orders 
        JOIN deliveries ON orders.id = deliveries.order_id
        JOIN payments ON orders.id = payments.order_id
        WHERE orders.order_uid = $1`

	row := s.db.QueryRowContext(ctx, query, orderUID)
	var paymentTime time.Time
	if err := row.Scan(orderId,

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
		&order.OofShard,

		&order.Delivery.Name,
		&order.Delivery.Phone,
		&order.Delivery.Zip,
		&order.Delivery.City,
		&order.Delivery.Address,
		&order.Delivery.Region,
		&order.Delivery.Email,

		&order.Payment.Transaction,
		&order.Payment.RequestID,
		&order.Payment.Currency,
		&order.Payment.Provider,
		&order.Payment.Amount,
		&paymentTime,
		&order.Payment.Bank,
		&order.Payment.DeliveryCost,
		&order.Payment.GoodsTotal,
		&order.Payment.CustomFee,
	); err != nil {
		return err
	}
	order.Payment.PaymentDT = paymentTime.Unix()
	return nil
}

func queryItems(ctx context.Context, s *Storage, items *[]models.Item, orderId int) error {
	query := `SELECT 
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
        FROM items WHERE order_id = $1`

	rows, err := s.db.QueryContext(ctx, query, orderId)
	if err != nil {
		return err
	}

	defer rows.Close()

	for rows.Next() {
		var item models.Item
		err := rows.Scan(
			&item.ChrtID,
			&item.TrackNumber,
			&item.Price,
			&item.Rid,
			&item.Name,
			&item.Sale,
			&item.Size,
			&item.TotalPrice,
			&item.NmID,
			&item.Brand,
			&item.Status,
		)
		if err != nil {
			return err
		}
		*items = append(*items, item)
	}
	return rows.Err()
}

func (s *Storage) GetOrders(amount ...int) ([]*models.Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	query := `SELECT 

        orders.id,
        orders.order_uid, 
        orders.track_number, 
        orders.entry, 
        orders.locale, 
        orders.internal_signature,
        orders.customer_id,
        orders.delivery_service,
        orders.shardkey,
        orders.sm_id,
        orders.date_created,
        orders.oof_shard,

        deliveries.name, 
        deliveries.phone, 
        deliveries.zip, 
        deliveries.city, 
        deliveries.address,
        deliveries.region,
        deliveries.email,

        payments.transaction, 
        payments.request_id, 
        payments.currency, 
        payments.provider, 
        payments.amount,
        payments.payment_dt,
        payments.bank,
        payments.delivery_cost,
        payments.goods_total,
        payments.custom_fee

    FROM orders 
    JOIN deliveries ON orders.id = deliveries.order_id
    JOIN payments ON orders.id = payments.order_id`

	if len(amount) > 0 && amount[0] > 0 {
		query += fmt.Sprintf("\nLIMIT %d", amount[0])
	}

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*models.Order
	var orderId int

	for rows.Next() {
		order := new(models.Order)
		var paymentTime time.Time

		err := rows.Scan(
			&orderId,
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
			&order.OofShard,

			&order.Delivery.Name,
			&order.Delivery.Phone,
			&order.Delivery.Zip,
			&order.Delivery.City,
			&order.Delivery.Address,
			&order.Delivery.Region,
			&order.Delivery.Email,

			&order.Payment.Transaction,
			&order.Payment.RequestID,
			&order.Payment.Currency,
			&order.Payment.Provider,
			&order.Payment.Amount,
			&paymentTime,
			&order.Payment.Bank,
			&order.Payment.DeliveryCost,
			&order.Payment.GoodsTotal,
			&order.Payment.CustomFee,
		)
		if err != nil {
			return nil, err
		}
		order.Payment.PaymentDT = paymentTime.Unix()
		if err := queryItems(ctx, s, &order.Items, orderId); err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, rows.Err()
}
