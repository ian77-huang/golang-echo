CREATE TABLE IF NOT EXISTS `user` (
    `id` INTEGER PRIMARY KEY AUTOINCREMENT,
    `account` VARCHAR(100) NOT NULL,
    `password` VARCHAR(255) NOT NULL,
    `created_at` DATETIME NULL,
    `updated_at` DATETIME NULL,
    `is_active` TINYINT NOT NULL DEFAULT 1,
    `is_admin` TINYINT NOT NULL DEFAULT 0,
    UNIQUE (`account`)
);