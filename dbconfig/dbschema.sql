
-- salesanalytics.customers definition

CREATE TABLE `customers` (
  `customer_id` varchar(100)  NOT NULL,
  `name` varchar(100)  DEFAULT NULL,
  `email` varchar(100) DEFAULT NULL,
  `address` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`customer_id`)
) 

-- salesanalytics.products definition

CREATE TABLE `products` (
  `product_id` varchar(100)  NOT NULL,
  `product_name` varchar(255)  DEFAULT NULL,
  `category` varchar(100) DEFAULT NULL,
  PRIMARY KEY (`product_id`)
) 

-- salesanalytics.orders definition

CREATE TABLE `orders` (
  `order_id` bigint NOT NULL,
  `customer_id` varchar(100) DEFAULT NULL,
  `payment_method` varchar(100)  DEFAULT NULL,
  `date_of_sale` datetime DEFAULT NULL,
  `region` varchar(100) DEFAULT NULL,
  PRIMARY KEY (`order_id`)
) 

-- salesanalytics.order_items definition

CREATE TABLE `order_items` (
  `order_item_id` bigint NOT NULL AUTO_INCREMENT,
  `order_id` bigint NOT NULL,
  `product_id` varchar(100) DEFAULT NULL,
  `quantity_sold` int DEFAULT NULL,
  `unit_price` decimal(10,2) DEFAULT NULL,
  `shipping_cost` decimal(10,2) DEFAULT NULL,
  `discount` decimal(10,2) DEFAULT NULL,
  PRIMARY KEY (`order_item_id`)
)
-- salesanalytics.sales_refresh_logs definition

CREATE TABLE `sales_refresh_logs` (
  `id` int NOT NULL AUTO_INCREMENT,
  `status` varchar(20) DEFAULT NULL,
  `type` varchar(20) DEFAULT NULL,
  `created_date` datetime DEFAULT NULL,
  `updated_date` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
)