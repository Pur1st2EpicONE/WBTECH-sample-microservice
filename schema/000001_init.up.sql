CREATE TABLE IF NOT EXISTS orders (
    id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    order_uid VARCHAR(255) UNIQUE NOT NULL,
    track_number VARCHAR(255) UNIQUE NOT NULL,
    entry VARCHAR(255) NOT NULL,
    locale VARCHAR(10) NOT NULL,
    internal_signature VARCHAR(255) NULL,
    customer_id VARCHAR(255) NOT NULL,
    delivery_service VARCHAR(255) NOT NULL,
    shardkey VARCHAR(10) NOT NULL,
    sm_id INTEGER NOT NULL,
    date_created TIMESTAMP NOT NULL,
    oof_shard VARCHAR(10) NOT NULL
);

CREATE TABLE IF NOT EXISTS deliveries (
    id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    phone VARCHAR(50) NOT NULL,
    zip VARCHAR(50) NOT NULL,
    city VARCHAR(100) NOT NULL,
    address VARCHAR(255) NOT NULL,
    region VARCHAR(255) NOT NULL,
    email VARCHAR(100) NOT NULL,
    order_id INTEGER UNIQUE NOT NULL,
    CONSTRAINT fk_delivery_order_id FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS payments (
    id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    transaction VARCHAR(255) NOT NULL,
    request_id VARCHAR(255) DEFAULT NULL,
    currency VARCHAR(10) NOT NULL,
    provider VARCHAR(50) NOT NULL,
    amount NUMERIC(10,2) NOT NULL,
    payment_dt TIMESTAMP NOT NULL,
    bank VARCHAR(50) NULL,
    delivery_cost NUMERIC(10,2) NOT NULL,
    goods_total NUMERIC(10,2) NOT NULL,
    custom_fee NUMERIC(10,2) NOT NULL DEFAULT 0,
    order_id INTEGER UNIQUE NOT NULL,
    CONSTRAINT fk_payment_order_id FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE,
    CONSTRAINT fk_payment_transaction FOREIGN KEY (transaction) REFERENCES orders(order_uid) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS items (
    id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    chrt_id INTEGER NOT NULL,
    track_number VARCHAR(255) NOT NULL,
    price NUMERIC(10,2) NOT NULL,
    rid VARCHAR(255) NOT NULL,
    name VARCHAR(100) NOT NULL,
    sale INTEGER NOT NULL DEFAULT 0,
    size VARCHAR(10) DEFAULT NULL,
    total_price NUMERIC(10,2) NOT NULL,
    nm_id INTEGER NOT NULL,
    brand VARCHAR(100) NOT NULL,
    status INTEGER NOT NULL,
    order_id INTEGER NOT NULL,
    CONSTRAINT fk_item_order_id FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE,
    CONSTRAINT fk_item_track_number FOREIGN KEY (track_number) REFERENCES orders(track_number) ON DELETE CASCADE
);

-- Queries will primarily search orders, deliveries, and payments by their primary keys (id columns),
-- while lookups for items within an order will use the idx_items index on order_id for faster access.
CREATE INDEX idx_items_order_id ON items(order_id);
