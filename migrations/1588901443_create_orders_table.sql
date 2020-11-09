use `market`
DROP TABLE IF EXISTS `orders`;
CREATE TABLE `orders` (
  `id` INT(11) NOT NULL AUTO_INCREMENT,
  `order_public_id` CHAR(100) NOT NULL,
  `active` BOOLEAN NOT NULL DEFAULT true,
  `status` char(15) NOT NULL DEFAULT 'open',
  `cents` INT NOT NULL DEFAULT 0,
  `side` char(5) NOT NULL, 
  `gold_tenth_milli_grams` INT NOT NULL DEFAULT 0,
  `created_at` INT(11) NOT NULL,
  `updated_at` INT(11) NOT NULL,
  `user_public_id` CHAR(100) NOT NULL,
PRIMARY KEY (`id`)
) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
