CREATE TABLE `user` (
  `id` bigint(20) PRIMARY KEY AUTO_INCREMENT,
  `uuid` varchar(50) NOT NULL,
  `first_name` varchar(50),
  `last_name` varchar(50) DEFAULT null,
  `email` varchar(80) NOT NULL,
  `phone` varchar(15) DEFAULT null,
  `password` text NOT NULL,
  `intro` tinytext DEFAULT null,
  `status` varchar(8) NOT NULL DEFAULT "inactive",
  `profile` text DEFAULT null COMMENT 'The user details.',
  `user_groups_id` bigint(20) NOT NULL,
  `last_login` timestamp(6) DEFAULT "now()",
  `created_at` timestamp(6) DEFAULT "now()",
  `updated_at` timestamp(6),
  `deleted_at` timestamp(6)
);

CREATE TABLE `roles` (
  `id` bigint(20) PRIMARY KEY AUTO_INCREMENT,
  `title` varchar(75) NOT NULL,
  `action` varchar(100) NOT NULL,
  `description` tinytext,
  `active` tinyint(1) NOT NULL DEFAULT 0,
  `content` text DEFAULT null,
  `created_at` timestamp(6) DEFAULT "now()",
  `updated_at` timestamp(6),
  `deleted_at` timestamp(6)
);

CREATE TABLE `mapping_roles` (
  `id` bigint(20) PRIMARY KEY AUTO_INCREMENT,
  `role_id` bigint(20) NOT NULL,
  `user_groups_id` bigint(20) NOT NULL,
  `created_at` timestamp(6) DEFAULT "now()",
  `updated_at` timestamp(6),
  `deleted_at` timestamp(6)
);

CREATE TABLE `user_groups` (
  `id` bigint(20) PRIMARY KEY AUTO_INCREMENT,
  `name` varchar(45) NOT NULL,
  `active` tinyint(1) NOT NULL DEFAULT 0,
  `created_at` timestamp(6) DEFAULT "now()",
  `updated_at` timestamp(6),
  `deleted_at` timestamp(6)
);

CREATE UNIQUE INDEX `user_index_0` ON `user` (`email`);

CREATE UNIQUE INDEX `user_index_1` ON `user` (`phone`);

CREATE UNIQUE INDEX `user_groups_index_2` ON `user_groups` (`name`);

ALTER TABLE `user` ADD FOREIGN KEY (`user_groups_id`) REFERENCES `user_groups` (`id`);

ALTER TABLE `mapping_roles` ADD FOREIGN KEY (`role_id`) REFERENCES `roles` (`id`);

ALTER TABLE `mapping_roles` ADD FOREIGN KEY (`user_groups_id`) REFERENCES `user_groups` (`id`);
