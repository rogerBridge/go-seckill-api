create database if not exists shop;
use shop;

drop table if exists shop.`goods`;

CREATE TABLE shop.goods (
    `id` INT UNSIGNED auto_increment NOT NULL,
    `product_id` INT UNSIGNED NOT NULL,
    `product_name` varchar(100) NOT NULL,
    `inventory` INT UNSIGNED NOT NULL,
    `version` varchar(32) DEFAULT 'alpha 0.1' NOT NULL COMMENT '版本号',
    `is_delete` BOOL DEFAULT FALSE COMMENT 'is delete or not',
    `gmt_create` DATETIME DEFAULT current_timestamp NOT NULL,
    `gmt_update` DATETIME on update current_timestamp,
    primary key(`id`),
    unique key (`product_id`)
)
    ENGINE=InnoDB
    DEFAULT CHARSET=utf8mb4
    COLLATE=utf8mb4_0900_ai_ci
    COMMENT='存放商品信息的table';
