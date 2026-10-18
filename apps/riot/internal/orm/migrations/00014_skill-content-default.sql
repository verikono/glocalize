-- +goose Up
PRAGMA foreign_keys = off;
CREATE TABLE `new_Skill` (
  `id` text NOT NULL,
  `name` text NOT NULL,
  `description` text NOT NULL,
  `url` text NOT NULL,
  `content` text NOT NULL DEFAULT '',
  `keywords` text NOT NULL DEFAULT '',
  `createdAt` timestamp NOT NULL DEFAULT (CURRENT_TIMESTAMP),
  `updatedAt` timestamp NOT NULL,
  PRIMARY KEY (`id`)
);
INSERT INTO `new_Skill` (`id`, `name`, `description`, `url`, `content`, `keywords`, `createdAt`, `updatedAt`) SELECT `id`, `name`, `description`, `url`, IFNULL(`content`, '') AS `content`, `keywords`, `createdAt`, `updatedAt` FROM `Skill`;
DROP TABLE `Skill`;
ALTER TABLE `new_Skill` RENAME TO `Skill`;
CREATE UNIQUE INDEX `Skill_url_key` ON `Skill` (`url`);
CREATE TABLE `new_StyleToken` (
  `style` text NOT NULL,
  `value` text NOT NULL,
  `mode` text NOT NULL,
  `type` text NOT NULL,
  `token` text NOT NULL,
  `app_id` text NOT NULL,
  `variable` text NULL,
  PRIMARY KEY (`app_id`, `mode`, `token`, `style`),
  CONSTRAINT `StyleToken_app_id_fkey` FOREIGN KEY (`app_id`) REFERENCES `App` (`name`) ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT `StyleToken_mode_check` CHECK ("mode" IN ('light', 'dark', 'primitive')),
  CONSTRAINT `StyleToken_type_check` CHECK ("type" IN ('color', 'spatial', 'size', 'dimension', 'font', 'shadow', 'zindex', 'gradient', 'weight', 'css')),
  CONSTRAINT `StyleToken_token_check` CHECK ("token" IN ('semantic', 'applied', 'primitive'))
);
INSERT INTO `new_StyleToken` (`style`, `value`, `mode`, `type`, `token`, `app_id`, `variable`) SELECT `style`, `value`, `mode`, `type`, `token`, `app_id`, `variable` FROM `StyleToken`;
DROP TABLE `StyleToken`;
ALTER TABLE `new_StyleToken` RENAME TO `StyleToken`;
PRAGMA foreign_keys = on;

-- +goose Down
PRAGMA foreign_keys = off;
CREATE TABLE `new_Skill` (
  `id` text NOT NULL,
  `name` text NOT NULL,
  `description` text NOT NULL,
  `url` text NOT NULL,
  `content` text NOT NULL,
  `createdAt` timestamp NOT NULL DEFAULT (CURRENT_TIMESTAMP),
  `updatedAt` timestamp NOT NULL,
  `keywords` text NOT NULL DEFAULT '',
  PRIMARY KEY (`id`)
);
INSERT INTO `new_Skill` (`id`, `name`, `description`, `url`, `content`, `createdAt`, `updatedAt`, `keywords`) SELECT `id`, `name`, `description`, `url`, `content`, `createdAt`, `updatedAt`, `keywords` FROM `Skill`;
DROP TABLE `Skill`;
ALTER TABLE `new_Skill` RENAME TO `Skill`;
CREATE UNIQUE INDEX `Skill_url_key` ON `Skill` (`url`);
CREATE TABLE `new_StyleToken` (
  `style` text NOT NULL,
  `value` text NOT NULL,
  `mode` text NOT NULL,
  `type` text NOT NULL,
  `token` text NOT NULL,
  `app_id` text NOT NULL,
  `variable` text NULL,
  PRIMARY KEY (`style`, `mode`, `token`, `app_id`),
  CONSTRAINT `StyleToken_app_id_fkey` FOREIGN KEY (`app_id`) REFERENCES `App` (`name`) ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT `StyleToken_mode_check` CHECK ("mode" IN ('light', 'dark', 'primitive')),
  CONSTRAINT `StyleToken_type_check` CHECK ("type" IN ('color', 'spatial', 'size', 'dimension', 'font', 'shadow', 'zindex', 'gradient', 'weight', 'css')),
  CONSTRAINT `StyleToken_token_check` CHECK ("token" IN ('semantic', 'applied', 'primitive'))
);
INSERT INTO `new_StyleToken` (`style`, `value`, `mode`, `type`, `token`, `app_id`, `variable`) SELECT `style`, `value`, `mode`, `type`, `token`, `app_id`, `variable` FROM `StyleToken`;
DROP TABLE `StyleToken`;
ALTER TABLE `new_StyleToken` RENAME TO `StyleToken`;
PRAGMA foreign_keys = on;
