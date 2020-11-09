USE `market`;

DROP TABLE IF EXISTS `daily_balance`;
CREATE TABLE `daily_balance` (
  `id` INT(11) NOT NULL AUTO_INCREMENT,
  `gold_tenth_milli_grams` INT NOT NULL DEFAULT 0,
  `cents` INT NOT NULL DEFAULT 0,
  `balance` INT NOT NULL DEFAULT 0,
  `user_id` BIGINT NOT NULL,
  `user_public_id` CHAR(100) NOT NULL,
`date_string` char(10) NOT NULL,
`unix_time` INT(11) NOT NULL,
PRIMARY KEY (`id`)
) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
