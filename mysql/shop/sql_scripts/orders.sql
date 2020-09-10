drop table if exists shop.`orders`;

CREATE TABLE shop.`orders` (
                              id INT UNSIGNED auto_increment NOT NULL,
                              order_number varchar(100) NOT NULL,
                              user_id INT UNSIGNED NOT NULL,
                              product_id INT UNSIGNED NOT NULL,
                              purchase_number INT UNSIGNED NOT NULL,
                              order_datetime DATETIME NOT NULL,
                              status varchar(100) NOT NULL,
                              `version` varchar(32) default 'alpha 0.1' NOT NULL,
                              is_delete BOOL default false NOT NULL,
                              gmt_create DATETIME default current_timestamp NOT NULL,
                              gmt_update DATETIME on update current_timestamp,
                              primary key(id)
)
    ENGINE=InnoDB
    DEFAULT CHARSET=utf8mb4
    COLLATE=utf8mb4_0900_ai_ci
    COMMENT='订单表';