use `market`
DROP TABLE IF EXISTS `transactions`;
CREATE TABLE `transactions` (
  `id` INT(11) NOT NULL AUTO_INCREMENT,
  `user_id` BIGINT NOT NULL,
  `user_public_id` CHAR(100) NOT NULL,
  `deposit_public_id` CHAR(100),
  `order_public_id` CHAR(100),
  `transaction_type` CHAR(12) NOT NULL,
  `unix_time` INT(11) NOT NULL,
  `date_string` char(10) NOT NULL,
  `cents` INT NOT NULL DEFAULT 0,
  `side` char(5) NOT NULL,
  `gold_tenth_milli_grams` INT NOT NULL DEFAULT 0,
PRIMARY KEY (`id`)
) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci;

