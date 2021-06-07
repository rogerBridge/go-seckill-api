CREATE TABLE shop.purchase_limits (
      id INT UNSIGNED auto_increment NOT NULL,
      product_id INT UNSIGNED NOT NULL,
      limit_num int unsigned default 1,
      start_purchase_time DATETIME NOT NULL,
      end_purchase_time DATETIME NOT NULL,
      version varchar(100) DEFAULT 'alpha 0.1' NOT NULL,
      is_delete BOOL DEFAULT true NOT NULL,
      gmt_create DATETIME DEFAULT current_timestamp NOT NULL,
      gmt_update DATETIME on update current_timestamp NULL,
      unique key(product_id),
      primary key(id)
)
    ENGINE=InnoDB
    DEFAULT CHARSET=utf8mb4
    COLLATE=utf8mb4_0900_ai_ci
    COMMENT='购买商品的限制, 每个人总共可以购买的数量 + 可购买时间段 + ...';