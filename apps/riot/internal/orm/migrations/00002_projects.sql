-- +goose Up
CREATE TABLE `Project` (
  `name` text NOT NULL,
  `skillsPath` text NOT NULL,
  `agentsPath` text NOT NULL,
  PRIMARY KEY (`name`)
);
CREATE UNIQUE INDEX `Project_name_key` ON `Project` (`name`);

-- +goose Down
PRAGMA foreign_keys = off;
DROP TABLE `Project`;
PRAGMA foreign_keys = on;
