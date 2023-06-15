-- Create default database
-- CREATE DATABASE mydb;

-- Add some tables (for demo only)
-- Relationships: users <--1:N--> orders <--1:1--> shipments
CREATE TABLE IF NOT EXISTS users
(
    id   INTEGER       PRIMARY KEY,
    name VARCHAR(255)  NOT NULL
);

CREATE TABLE IF NOT EXISTS orders
(
    id      INTEGER       PRIMARY KEY,
    name    VARCHAR(255)  NOT NULL,
    user_id INTEGER       NOT NULL,

    CONSTRAINT fk_users_orders FOREIGN KEY (user_id) REFERENCES users (id)
);

CREATE TYPE shipment_status AS ENUM('NOT_SHIPPED', 'SHIPPED', 'DONE');

CREATE TABLE IF NOT EXISTS shipments
(
    id       INTEGER           PRIMARY KEY,
    status   shipment_status   NOT NULL,
    order_id INTEGER           NOT NULL,

    CONSTRAINT fk_shipments_orders FOREIGN KEY (order_id) REFERENCES orders (id)
);
