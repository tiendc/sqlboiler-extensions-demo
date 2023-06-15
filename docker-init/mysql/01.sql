-- Create default database
CREATE DATABASE IF NOT EXISTS `mydb` CHARACTER SET utf8mb4 COLLATE utf8mb4_bin;

-- Add some tables (for demo only)
-- Relationships: users <--1:N--> orders <--1:1--> shipments
CREATE TABLE IF NOT EXISTS `users`
(
    `id`   INTEGER       PRIMARY KEY,
    `name` VARCHAR(255)  NOT NULL
);

CREATE TABLE IF NOT EXISTS `orders`
(
    `id`      INTEGER       PRIMARY KEY,
    `name`    VARCHAR(255)  NOT NULL,
    `user_id` INTEGER       NOT NULL,

    CONSTRAINT `fk_users_orders` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
);

CREATE TABLE IF NOT EXISTS `shipments`
(
    `id`       INTEGER                                 PRIMARY KEY,
    `status`   ENUM('NOT_SHIPPED', 'SHIPPED', 'DONE')  NOT NULL,
    `order_id` INTEGER                                 NOT NULL,

    CONSTRAINT `fk_shipments_orders` FOREIGN KEY (`order_id`) REFERENCES `orders` (`id`)
);
