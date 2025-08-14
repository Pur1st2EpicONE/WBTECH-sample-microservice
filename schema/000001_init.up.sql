CREATE TABLE IF NOT EXISTS orders (
    id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    order_uid VARCHAR(255) UNIQUE NOT NULL,
    track_number VARCHAR(255) UNIQUE NOT NULL,
    entry VARCHAR(255) NOT NULL,
    locale VARCHAR(10) NOT NULL,
    internal_signature VARCHAR(255),
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
    order_id INTEGER NOT NULL,
    CONSTRAINT fk_delivery_order_id FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS payments (
    id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    transaction VARCHAR(255) NOT NULL,
    request_id VARCHAR(255),
    currency VARCHAR(10) NOT NULL,
    provider VARCHAR(50) NOT NULL,
    amount NUMERIC(10,2) NOT NULL,
    payment_dt TIMESTAMP NOT NULL,
    bank VARCHAR(50),
    delivery_cost NUMERIC(10,2) NOT NULL,
    goods_total NUMERIC(10,2) NOT NULL,
    custom_fee NUMERIC(10,2) NOT NULL,
    order_id INTEGER NOT NULL,
    CONSTRAINT fk_payment_order_id FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS items (
    id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    chrt_id INTEGER NOT NULL,
    track_number VARCHAR(255) NOT NULL,
    price NUMERIC(10,2) NOT NULL,
    rid VARCHAR(255) NOT NULL,
    name VARCHAR(100) NOT NULL,
    sale NUMERIC(10,2) NOT NULL,
    size VARCHAR(10),
    total_price NUMERIC(10,2) NOT NULL,
    nm_id INTEGER NOT NULL,
    brand VARCHAR(100) NOT NULL,
    status INTEGER NOT NULL,
    order_id INTEGER NOT NULL,
    CONSTRAINT fk_item_order_id FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE
);
