-- +goose Up
PRAGMA foreign_keys = off;
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
CREATE TABLE `Reference` (
  `name` text NOT NULL,
  `description` text NOT NULL,
  `accessType` text NOT NULL,
  `resourceType` text NOT NULL,
  `content` text NULL,
  CONSTRAINT `Reference_accessType_check` CHECK ("accessType" IN ('static', 'url')),
  CONSTRAINT `Reference_resourceType_check` CHECK ("resourceType" IN ('literature', 'image', 'video', 'commands', 'script', 'configuration'))
);
CREATE UNIQUE INDEX `Reference_name_key` ON `Reference` (`name`);
PRAGMA foreign_keys = on;

-- +goose Down
PRAGMA foreign_keys = off;
DROP TABLE `Reference`;
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
