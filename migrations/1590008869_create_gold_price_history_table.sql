USE `market`;

DROP TABLE IF EXISTS `gold_price_history`;
CREATE TABLE `gold_price_history` (
  `id` INT(11) NOT NULL AUTO_INCREMENT,
  `cents` INT NOT NULL DEFAULT 0,
  `price_date` INT(11) NOT NULL,
  `date_string` char(10) NOT NULL,
  `fast_search` boolean DEFAULT 0, 
   PRIMARY KEY (`id`)
) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci;

