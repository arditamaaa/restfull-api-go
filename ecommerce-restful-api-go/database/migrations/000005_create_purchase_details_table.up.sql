DROP TABLE IF EXISTS `purchase_details`;
CREATE TABLE `purchase_details` (
    `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    `purchase_id` bigint UNSIGNED NOT NULL,
    `product_id` bigint UNSIGNED NOT NULL,
    `qty` int NOT NULL,
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `deleted_at` timestamp NULL,
    INDEX `purchase_details_purchase_id_foreign`(`purchase_id` ASC) USING BTREE,
    INDEX `purchase_details_product_id_foreign`(`product_id` ASC) USING BTREE,
    CONSTRAINT `purchase_details_purchase_id_foreign` FOREIGN KEY (`purchase_id`) REFERENCES `purchases` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT,
    CONSTRAINT `purchase_details_product_id_foreign` FOREIGN KEY (`product_id`) REFERENCES `products` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT
);