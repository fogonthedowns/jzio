use `market`
DROP TABLE IF EXISTS `users`;
CREATE TABLE `users` (
  `id` INT(11) NOT NULL AUTO_INCREMENT,
  `active` BOOLEAN NOT NULL DEFAULT true,
  `created_at` INT(11) NOT NULL,
  `user_public_id` CHAR(100) NOT NULL,
PRIMARY KEY (`id`)
) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci;

