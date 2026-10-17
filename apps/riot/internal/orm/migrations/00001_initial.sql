-- +goose Up
CREATE TABLE `Skill` (
  `id` text NOT NULL,
  `name` text NOT NULL,
  `description` text NOT NULL,
  `url` text NOT NULL,
  `content` text NOT NULL,
  `createdAt` timestamp NOT NULL DEFAULT (CURRENT_TIMESTAMP),
  `updatedAt` timestamp NOT NULL,
  PRIMARY KEY (`id`)
);
CREATE UNIQUE INDEX `Skill_url_key` ON `Skill` (`url`);

-- +goose Down
PRAGMA foreign_keys = off;
DROP TABLE `Skill`;
PRAGMA foreign_keys = on;
