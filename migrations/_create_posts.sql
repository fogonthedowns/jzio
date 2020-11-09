CREATE DATABASE IF NOT EXISTS `market` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
USE `market`;

DROP TABLE IF EXISTS `posts`;
CREATE TABLE `posts` (
  `id` INT(11) NOT NULL AUTO_INCREMENT,
  `url` VARCHAR(255) NOT NULL,
  `title` VARCHAR(255) NOT NULL,
  `thumbnail` VARCHAR(255) DEFAULT "",
  `likes` INT(11) DEFAULT 0,
  `views` INT(11) DEFAULT 0,
  `created_at` INT(11) NOT NULL,
  `duration_seconds` int(11) DEFAULT 0,
PRIMARY KEY (`id`)
)
