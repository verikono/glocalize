-- +goose Up
PRAGMA foreign_keys = off;
ALTER TABLE `Skill` ADD COLUMN `keywords` text NOT NULL DEFAULT '';
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
  PRIMARY KEY (`id`)
);
INSERT INTO `new_Skill` (`id`, `name`, `description`, `url`, `content`, `createdAt`, `updatedAt`) SELECT `id`, `name`, `description`, `url`, `content`, `createdAt`, `updatedAt` FROM `Skill`;
DROP TABLE `Skill`;
ALTER TABLE `new_Skill` RENAME TO `Skill`;
CREATE UNIQUE INDEX `Skill_url_key` ON `Skill` (`url`);
PRAGMA foreign_keys = on;
