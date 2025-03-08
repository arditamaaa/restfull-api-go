CREATE TABLE `users` (
	`id` bigint UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
	`name` varchar(200) NOT NULL,
	`email` varchar(100) NOT NULL,
	`password` varchar(200) NOT NULL,
	`role` enum('user', 'admin') NOT NULL DEFAULT 'user',
	`created_at` timestamp NULL,
	`updated_at` timestamp NULL,
	`deleted_at` timestamp NULL,
	UNIQUE INDEX `users_email_unique`(`email`)
);