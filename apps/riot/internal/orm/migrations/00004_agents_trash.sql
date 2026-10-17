-- +goose Up
PRAGMA foreign_keys = off;
CREATE TABLE `Agent` (
  `id` text NOT NULL,
  `name` text NOT NULL,
  `description` text NOT NULL,
  `allowed-tools` text NOT NULL,
  `metadata` text NOT NULL,
  `createdAt` timestamp NOT NULL DEFAULT (CURRENT_TIMESTAMP),
  `updatedAt` timestamp NOT NULL,
  PRIMARY KEY (`id`)
);
CREATE UNIQUE INDEX `Agent_name_key` ON `Agent` (`name`);
CREATE TABLE `Trash` (
  `id` integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  `createdAt` timestamp NOT NULL DEFAULT (CURRENT_TIMESTAMP),
  `content` text NOT NULL
);
CREATE TABLE `TrashAgent` (
  `trashId` integer NOT NULL,
  `agentId` text NOT NULL,
  PRIMARY KEY (`trashId`, `agentId`),
  CONSTRAINT `TrashAgent_trashId_fkey` FOREIGN KEY (`trashId`) REFERENCES `Trash` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT `TrashAgent_agentId_fkey` FOREIGN KEY (`agentId`) REFERENCES `Agent` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION
);
CREATE TABLE `TrashSkill` (
  `trashId` integer NOT NULL,
  `skillId` text NOT NULL,
  PRIMARY KEY (`trashId`, `skillId`),
  CONSTRAINT `TrashSkill_trashId_fkey` FOREIGN KEY (`trashId`) REFERENCES `Trash` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT `TrashSkill_skillId_fkey` FOREIGN KEY (`skillId`) REFERENCES `Skill` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION
);
PRAGMA foreign_keys = on;

-- +goose Down
PRAGMA foreign_keys = off;
DROP TABLE `TrashAgent`;
DROP TABLE `TrashSkill`;
DROP TABLE `Trash`;
DROP TABLE `Agent`;
PRAGMA foreign_keys = on;
