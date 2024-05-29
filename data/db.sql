CREATE DATABASE IF NOT EXISTS web_golang_api;
USE web_golang_api;

CREATE TABLE IF NOT EXISTS customers(
    customer_id INT NOT NULL PRIMARY KEY AUTO_INCREMENT,
    customer_number INT NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS users (
    user_id VARCHAR(100) NOT NULL PRIMARY KEY,
    email VARCHAR(100) NOT NULL,
    password VARCHAR(100) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);

INSERT INTO `web_golang_api`.`customers` (`customer_number`, `first_name`, `last_name`, `created_at`) VALUES (1, 'Danilo', 'Sano', '2024-05-29 00:00:00');
INSERT INTO `web_golang_api`.`customers` (`customer_number`, `first_name`, `last_name`, `created_at`) VALUES (2, 'Cliente', 'Teste', '2024-05-04 00:00:00');

INSERT INTO `web_golang_api`.`users` (`user_id`, `email`, `password`, `created_at`, `updated_at`, `deleted_at`) VALUES ('4d1315b1-62ad-4711-8082-bb07f3bbc35f', 'admin', '$2a$10$F.G8FJMaAlVBuTIve1B.M.nLAAwxFH2ftandpEM9ymg76dRJYe.pa', '2024-05-08 20:49:40', '0000-00-00 00:00:00', '0000-00-00 00:00:00');