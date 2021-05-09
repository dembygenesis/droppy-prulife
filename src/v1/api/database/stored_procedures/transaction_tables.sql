/**
  Transaction
 */

CREATE TABLE `transaction` (
                               `id` int NOT NULL AUTO_INCREMENT,
                               `created_by` int DEFAULT NULL,
                               `updated_by` int DEFAULT NULL,
                               `is_active` tinyint DEFAULT NULL,
                               `amount` decimal(15,2) DEFAULT NULL,
                               `coin_amount` decimal(15,2) DEFAULT NULL,
                               `admin_allotted_id` int DEFAULT NULL,
                               `user_allotted_id` int DEFAULT NULL,
                               `money_in` tinyint(1) DEFAULT NULL,
                               PRIMARY KEY (`id`),
                               KEY `created_by` (`created_by`),
                               KEY `updated_by` (`updated_by`),
                               KEY `admin_alotted_id` (`admin_allotted_id`),
                               KEY `user_allotted_id` (`user_allotted_id`),
                               CONSTRAINT `transaction_ibfk_1` FOREIGN KEY (`created_by`) REFERENCES `user` (`id`),
                               CONSTRAINT `transaction_ibfk_2` FOREIGN KEY (`updated_by`) REFERENCES `user` (`id`),
                               CONSTRAINT `transaction_ibfk_3` FOREIGN KEY (`admin_allotted_id`) REFERENCES `user` (`id`),
                               CONSTRAINT `transaction_ibfk_4` FOREIGN KEY (`user_allotted_id`) REFERENCES `user` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

/**
  Coin Transaction
 */
CREATE TABLE `coin_transaction` (
                                    `id` int DEFAULT NULL,
                                    `created_by` int DEFAULT NULL,
                                    `updated_by` int DEFAULT NULL,
                                    `is_active` tinyint DEFAULT NULL,
                                    `transaction_id` int DEFAULT NULL COMMENT 'This is nullable, this transaction may come from a debit transaction, hence the need for this possible relationship. Or not, hence it can be nullable, and that "not" scenario is when we give coins to other droppshippers or sellers.',
                                    `user_id` int DEFAULT NULL,
                                    `type` enum('D','C') DEFAULT NULL,
                                    KEY `created_by` (`created_by`),
                                    KEY `updated_by` (`updated_by`),
                                    KEY `transaction_id` (`transaction_id`),
                                    KEY `user_id` (`user_id`),
                                    CONSTRAINT `coin_transaction_ibfk_1` FOREIGN KEY (`created_by`) REFERENCES `user` (`id`),
                                    CONSTRAINT `coin_transaction_ibfk_2` FOREIGN KEY (`updated_by`) REFERENCES `user` (`id`),
                                    CONSTRAINT `coin_transaction_ibfk_3` FOREIGN KEY (`transaction_id`) REFERENCES `transaction` (`id`),
                                    CONSTRAINT `coin_transaction_ibfk_4` FOREIGN KEY (`user_id`) REFERENCES `user` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci