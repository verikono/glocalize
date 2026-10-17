-- +goose Up
PRAGMA foreign_keys = off;
CREATE TABLE `_ProjectAgents` (
  `agentId` text NOT NULL,
  `projectName` text NOT NULL,
  PRIMARY KEY (`agentId`, `projectName`),
  CONSTRAINT `_ProjectAgents_agentId_fkey` FOREIGN KEY (`agentId`) REFERENCES `Agent` (`id`) ON UPDATE CASCADE ON DELETE CASCADE,
  CONSTRAINT `_ProjectAgents_projectName_fkey` FOREIGN KEY (`projectName`) REFERENCES `Project` (`name`) ON UPDATE CASCADE ON DELETE CASCADE
);
CREATE INDEX `_ProjectAgents_projectName_idx` ON `_ProjectAgents` (`projectName`);
PRAGMA foreign_keys = on;

-- +goose Down
PRAGMA foreign_keys = off;
DROP TABLE `_ProjectAgents`;
PRAGMA foreign_keys = on;
