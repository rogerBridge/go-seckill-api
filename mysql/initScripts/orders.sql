CREATE TABLE `orders` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `order_number` varchar(100) NOT NULL,
  `user_id` varchar(100) NOT NULL,
  `product_id` int unsigned NOT NULL,
  `purchase_number` int unsigned NOT NULL,
  `order_datetime` datetime NOT NULL,
  `status` varchar(100) NOT NULL,
  `version` varchar(32) NOT NULL DEFAULT 'alpha 0.1',
  `is_delete` tinyint(1) NOT NULL DEFAULT '0',
  `gmt_create` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `gmt_update` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=34928 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='订单表'