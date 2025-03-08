DROP TABLE IF EXISTS `purchases`;
create table `purchases` (
    `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    `user_id` bigint UNSIGNED NOT NULL,
    `code` varchar(20) NULL DEFAULT NULL,
    `payment_method` enum('Tunai', 'Transfer', 'QRIS') NOT NULL,
    `total` decimal(20,4) NOT NULL DEFAULT 0,
    `status` enum('DONE', 'CANCEL') NOT NULL DEFAULT 'DONE',
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `deleted_at` timestamp NULL,
    CONSTRAINT `purchases_user_id_foreign` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
);