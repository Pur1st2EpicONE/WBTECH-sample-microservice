package postgres_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/repository/postgres"
	mock_logger "github.com/Pur1st2EpicONE/WBTECH-sample-microservice/pkg/logger/mocks"
	"github.com/golang/mock/gomock"
	"github.com/jmoiron/sqlx"
)

func TestPostgresStorer_GetOrder_Success(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("failed to open mock db: %v", err)
	}
	defer func() { _ = db.Close() }()

	logger := mock_logger.NewMockLogger(gomock.NewController(t))
	ps := postgres.NewStorage(sqlx.NewDb(db, "postgres"), logger)

	orderUID := "test-uid"
	orderID := 1
	paymentTime := time.Now()

	orderQuery := `SELECT 

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

	mock.ExpectQuery(orderQuery).WithArgs(orderUID).WillReturnRows(sqlmock.NewRows([]string{

		"id",
		"order_uid",
		"track_number",
		"entry",
		"locale",
		"internal_signature",
		"customer_id",
		"delivery_service",
		"shardkey", "sm_id",
		"date_created",
		"oof_shard",
		"name",
		"phone",
		"zip",
		"city",
		"address",
		"region",
		"email",
		"transaction",
		"request_id",
		"currency",
		"provider",
		"amount",
		"payment_dt",
		"bank",
		"delivery_cost",
		"goods_total",
		"custom_fee",
	}).AddRow(

		orderID,
		orderUID,
		"track1",
		"entry",
		"en",
		"sig",
		123,
		"d_service",
		"shard",
		1,
		time.Now(),
		"oof",
		"name",
		"phone",
		"zip",
		"city",
		"address",
		"region",
		"email",
		"tx",
		"req",
		"USD",
		"prov",
		100,
		paymentTime,
		"bank",
		10,
		90,
		0,
	))

	itemsQuery := `SELECT 

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

	mock.ExpectQuery(itemsQuery).WithArgs(orderID).WillReturnRows(sqlmock.NewRows([]string{

		"chrt_id",
		"track_number",
		"price",
		"rid",
		"name",
		"sale",
		"size",
		"total_price",
		"nm_id",
		"brand",
		"status",
	}).AddRow(
		1,
		"track1",
		100,
		"rid1",
		"item1",
		0,
		"M",
		100,
		1,
		"brand1",
		1,
	))

	order, err := ps.GetOrder(orderUID)
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if len(order.Items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(order.Items))
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unfulfilled expectations: %v", err)
	}
}

func TestPostgresStorer_QueryAllButItems_RowsScanError(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("failed to open sqlmock: %v", err)
	}
	defer func() { _ = db.Close() }()

	logger := mock_logger.NewMockLogger(gomock.NewController(t))
	ps := postgres.NewStorage(sqlx.NewDb(db, "postgres"), logger)

	orderQuery := `SELECT

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

	mock.ExpectQuery(orderQuery).WithArgs(sqlmock.AnyArg()).WillReturnError(fmt.Errorf("query failed"))

	_, err = ps.GetOrder("some-uid")
	if err == nil {
		t.Fatalf("expected query error, got: %v", err)
	}
}

func TestPostgresStorer_QueryItems_RowsScanError(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("failed to open sqlmock: %v", err)
	}
	defer func() { _ = db.Close() }()

	logger := mock_logger.NewMockLogger(gomock.NewController(t))
	ps := postgres.NewStorage(sqlx.NewDb(db, "postgres"), logger)

	paymentTime := time.Now()

	orderQuery := `SELECT

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

	mock.ExpectQuery(orderQuery).
		WithArgs("uid1").
		WillReturnRows(sqlmock.NewRows([]string{
			"id",
			"order_uid",
			"track_number",
			"entry",
			"locale",
			"internal_signature",
			"customer_id",
			"delivery_service",
			"shardkey",
			"sm_id",
			"date_created",
			"oof_shard",
			"name",
			"phone",
			"zip",
			"city",
			"address",
			"region",
			"email",
			"transaction",
			"request_id",
			"currency",
			"provider",
			"amount",
			"payment_dt",
			"bank",
			"delivery_cost",
			"goods_total",
			"custom_fee",
		}).AddRow(
			1,
			"uid1",
			"track1",
			"entry",
			"en",
			"sig",
			123,
			"d_service",
			"shard",
			1,
			time.Now(),
			"oof",
			"name",
			"phone",
			"zip",
			"city",
			"address",
			"region",
			"email",
			"tx",
			"req",
			"USD",
			"prov",
			100,
			paymentTime,
			"bank",
			10,
			90,
			0,
		))

	itemsQuery := `SELECT

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

	mock.ExpectQuery(itemsQuery).WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{

		"chrt_id",
		"track_number",
		"price",
		"rid",
		"name",
		"sale",
		"size",
		"total_price",
		"nm_id",
		"brand",
		"status",
	}).AddRow(
		1,
		"track1",
		"not-an-int",
		"rid1",
		"item1",
		0,
		"M",
		100,
		1,
		"brand1",
		1,
	))

	_, err = ps.GetOrder("uid1")
	if err == nil {
		t.Fatalf("expected scan error, got nil")
	}
}

func TestPostgresStorer_QueryItems_QueryContextError(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("failed to open mock db: %v", err)
	}
	defer func() { _ = db.Close() }()

	logger := mock_logger.NewMockLogger(gomock.NewController(t))
	ps := postgres.NewStorage(sqlx.NewDb(db, "postgres"), logger)

	paymentTime := time.Now()

	allQuery := `SELECT

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

	mock.ExpectQuery(allQuery).WillReturnRows(sqlmock.NewRows([]string{

		"id",
		"order_uid",
		"track_number",
		"entry",
		"locale",
		"internal_signature",
		"customer_id",
		"delivery_service",
		"shardkey",
		"sm_id",
		"date_created",
		"oof_shard",
		"name",
		"phone",
		"zip",
		"city",
		"address",
		"region",
		"email",
		"transaction",
		"request_id",
		"currency",
		"provider",
		"amount",
		"payment_dt",
		"bank",
		"delivery_cost",
		"goods_total",
		"custom_fee",
	}).AddRow(
		1,
		"uid1",
		"track1",
		"entry",
		"en",
		"sig",
		123,
		"d_service",
		"shard",
		1,
		time.Now(),
		"oof",
		"name",
		"phone",
		"zip",
		"city",
		"address",
		"region",
		"email",
		"tx",
		"req",
		"USD",
		"prov",
		100,
		paymentTime,
		"bank",
		10,
		90,
		0,
	))

	itemsQuery := `SELECT

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

	mock.ExpectQuery(itemsQuery).WithArgs(sqlmock.AnyArg()).WillReturnError(fmt.Errorf("items query failed"))

	_, err = ps.GetOrders()
	if err == nil || err.Error() != "items query failed" {
		t.Fatalf("expected items query error, got: %v", err)
	}
}

func TestPostgresStorer_GetOrders_Success(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("failed to open mock db: %v", err)
	}
	defer func() { _ = db.Close() }()

	logger := mock_logger.NewMockLogger(gomock.NewController(t))
	ps := postgres.NewStorage(sqlx.NewDb(db, "postgres"), logger)

	orderID := 1
	paymentTime := time.Now()

	orderQuery := `SELECT 

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
		LIMIT 1`

	mock.ExpectQuery(orderQuery).WillReturnRows(sqlmock.NewRows([]string{

		"id",
		"order_uid",
		"track_number",
		"entry",
		"locale",
		"internal_signature",
		"customer_id",
		"delivery_service",
		"shardkey", "sm_id",
		"date_created",
		"oof_shard",
		"name",
		"phone",
		"zip",
		"city",
		"address",
		"region",
		"email",
		"transaction",
		"request_id",
		"currency",
		"provider",
		"amount",
		"payment_dt",
		"bank",
		"delivery_cost",
		"goods_total",
		"custom_fee",
	}).AddRow(
		orderID,
		"uid1",
		"track1",
		"entry",
		"en",
		"sig",
		123,
		"d_service",
		"shard",
		1,
		time.Now(),
		"oof",
		"name",
		"phone",
		"zip",
		"city",
		"address",
		"region",
		"email",
		"tx",
		"req",
		"USD",
		"prov",
		100,
		paymentTime,
		"bank",
		10,
		90,
		0,
	))

	itemsQuery := `SELECT 

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

	mock.ExpectQuery(itemsQuery).WithArgs(orderID).WillReturnRows(sqlmock.NewRows([]string{

		"chrt_id",
		"track_number",
		"price",
		"rid",
		"name",
		"sale",
		"size",
		"total_price",
		"nm_id",
		"brand",
		"status",
	}).AddRow(
		1,
		"track1",
		100,
		"rid1",
		"item1",
		0,
		"M",
		100,
		1,
		"brand1",
		1,
	))

	orders, err := ps.GetOrders(1)
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if len(orders) != 1 || len(orders[0].Items) != 1 {
		t.Fatalf("expected 1 order with 1 item, got %v", orders)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unfulfilled expectations: %v", err)
	}
}

func TestPostgresStorer_GetAllOrders_QueryContextError(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("failed to open mock db: %v", err)
	}
	defer func() { _ = db.Close() }()

	logger := mock_logger.NewMockLogger(gomock.NewController(t))
	ps := postgres.NewStorage(sqlx.NewDb(db, "postgres"), logger)

	allQuery := `SELECT

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

	mock.ExpectQuery(allQuery).WillReturnError(fmt.Errorf("all orders query failed"))

	_, err = ps.GetOrders()
	if err == nil || err.Error() != "all orders query failed" {
		t.Fatalf("expected all orders query error, got: %v", err)
	}
}

func TestPostgresStorer_GetAllOrders_RowsScanError(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("failed to open mock db: %v", err)
	}
	defer func() { _ = db.Close() }()

	logger := mock_logger.NewMockLogger(gomock.NewController(t))
	ps := postgres.NewStorage(sqlx.NewDb(db, "postgres"), logger)

	allQuery := `SELECT

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

	mock.ExpectQuery(allQuery).WillReturnRows(sqlmock.NewRows([]string{

		"id",
		"order_uid",
		"track_number",
		"entry",
		"locale",
		"internal_signature",
		"customer_id",
		"delivery_service",
		"shardkey",
		"sm_id",
		"date_created",
		"oof_shard",
		"name",
		"phone",
		"zip",
		"city",
		"address",
		"region",
		"email",
		"transaction",
		"request_id",
		"currency",
		"provider",
		"amount",
		"payment_dt",
		"bank",
		"delivery_cost",
		"goods_total",
		"custom_fee",
	}).AddRow(
		"string instead of an int",
		"uid1",
		"track1",
		"entry",
		"en",
		"sig",
		123,
		"d_service",
		"shard",
		1,
		time.Now(),
		"oof",
		"name",
		"phone",
		"zip",
		"city",
		"address",
		"region",
		"email",
		"tx",
		"req",
		"USD",
		"prov",
		100,
		time.Now(),
		"bank",
		10,
		90,
		0,
	))

	_, err = ps.GetOrders()
	if err == nil {
		t.Fatalf("expected scan error, got nil")
	}
}
