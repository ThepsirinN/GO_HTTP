CREATE TABLE `barko`.`user_info`
  (
     `id`               VARCHAR(80) NOT NULL,
     `first_name`       VARCHAR(80) NOT NULL,
     `last_name`        VARCHAR(80) NOT NULL,
     `created_datetime` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
     `updated_datetime` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
     PRIMARY KEY (`id`)
  )
engine = innodb
charset=utf8mb4
COLLATE utf8mb4_general_ci; 

INSERT INTO `user_info` (`id`, `first_name`, `last_name`) VALUES
('123', 'Barko', 'Barko');s