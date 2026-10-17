-- +goose Up
PRAGMA foreign_keys = off;
CREATE TABLE `AgentSkill` (
  `agentId` text NOT NULL,
  `skillId` text NOT NULL,
  `priority` text NOT NULL DEFAULT 'NORMAL',
  PRIMARY KEY (`agentId`, `skillId`),
  CONSTRAINT `AgentSkill_agentId_fkey` FOREIGN KEY (`agentId`) REFERENCES `Agent` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT `AgentSkill_skillId_fkey` FOREIGN KEY (`skillId`) REFERENCES `Skill` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT `AgentSkill_priority_check` CHECK ("priority" IN ('LOW', 'NORMAL', 'HIGH', 'ESSENTIAL'))
);
CREATE TABLE `AppSkill` (
  `app_id` text NOT NULL,
  `skillId` text NOT NULL,
  `priority` text NOT NULL DEFAULT 'NORMAL',
  PRIMARY KEY (`app_id`, `skillId`),
  CONSTRAINT `AppSkill_app_id_fkey` FOREIGN KEY (`app_id`) REFERENCES `App` (`name`) ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT `AppSkill_skillId_fkey` FOREIGN KEY (`skillId`) REFERENCES `Skill` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT `AppSkill_priority_check` CHECK ("priority" IN ('LOW', 'NORMAL', 'HIGH', 'ESSENTIAL'))
);
PRAGMA foreign_keys = on;

-- +goose Down
PRAGMA foreign_keys = off;
DROP TABLE `AppSkill`;
DROP TABLE `AgentSkill`;
PRAGMA foreign_keys = on;
