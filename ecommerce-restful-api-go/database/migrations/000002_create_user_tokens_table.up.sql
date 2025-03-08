create table `user_tokens` (
    `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    `user_id` bigint UNSIGNED NOT NULL,
    `token` varchar(255) NOT NULL,
    `expires` timestamp NULL,
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE INDEX `user_tokens_token_unique`(`token`),
    CONSTRAINT `user_tokens_user_id_foreign` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
);