CREATE DATABASE IF NOT EXISTS `market` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
USE `market`;

DROP TABLE IF EXISTS `vaults`;
CREATE TABLE `vaults` (
  `id` INT(11) NOT NULL AUTO_INCREMENT,
  `gold_tenth_milli_grams` INT NOT NULL DEFAULT 0,
  `gold_tenth_milli_grams_available` INT NOT NULL DEFAULT 0,
  `cents` INT NOT NULL DEFAULT 0,
  `cents_available` INT NOT NULL DEFAULT 0,
  `user_id` BIGINT NOT NULL,
  `updated_at` INT(11) NOT NULL,
  `vault_public_id` CHAR(100) NOT NULL,
  `user_public_id` CHAR(100) NOT NULL,
PRIMARY KEY (`id`)
) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci;

DROP TABLE IF EXISTS `deposits`;
CREATE TABLE `deposits` (
  `id` INT(11) NOT NULL AUTO_INCREMENT,
  `gold_tenth_milli_grams` INT NOT NULL DEFAULT 0,
  `cents` INT NOT NULL DEFAULT 0,
  `user_id` BIGINT NOT NULL,
  `created_at` INT(11) NOT NULL,
  `deposit_public_id` CHAR(100) NOT NULL,
  `user_public_id` CHAR(100) NOT NULL,
PRIMARY KEY (`id`)
) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
