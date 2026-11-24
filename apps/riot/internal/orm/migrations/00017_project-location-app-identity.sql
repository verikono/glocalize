-- +goose Up
PRAGMA foreign_keys = off;

-- Project gains an absolute repo location. "N/A" marks a project with no
-- checkout of its own (the shared "all" project).
CREATE TABLE `new_Project` (
  `name` text NOT NULL,
  `location` text NOT NULL DEFAULT 'N/A',
  `skillsPath` text NOT NULL,
  `agentsPath` text NOT NULL,
  PRIMARY KEY (`name`)
);
INSERT INTO `new_Project` (`name`, `location`, `skillsPath`, `agentsPath`)
SELECT `name`,
  CASE `name`
    WHEN 'all' THEN 'N/A'
    WHEN 'launchboard' THEN '/opt/code/launchboard/'
    WHEN 'kaltra' THEN '/opt/code/kaltra/'
    WHEN 'glocalize' THEN '~/glocalize/'
    ELSE 'N/A'
  END,
  `skillsPath`, `agentsPath`
FROM `Project`;

-- the "verikon" project is retired; its join rows go with it.
DELETE FROM `ProjectSkill` WHERE `projectName` = 'verikon';
DELETE FROM `_ProjectAgents` WHERE `projectName` = 'verikon';
DELETE FROM `ProjectSpecialization` WHERE `generalizesName` = 'verikon' OR `specializesName` = 'verikon';

DROP TABLE `Project`;
ALTER TABLE `new_Project` RENAME TO `Project`;
CREATE UNIQUE INDEX `Project_name_key` ON `Project` (`name`);

-- App is re-keyed from (name) to (language, name). The two existing rows are
-- remapped onto the name their own manifest declares.
CREATE TABLE `new_App` (
  `type` text NOT NULL DEFAULT 'app',
  `language` text NOT NULL DEFAULT 'typescript',
  `name` text NOT NULL,
  `location` text NOT NULL DEFAULT '',
  `project_id` text NOT NULL,
  PRIMARY KEY (`language`, `name`),
  CONSTRAINT `App_project_id_fkey` FOREIGN KEY (`project_id`) REFERENCES `Project` (`name`) ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT `App_type_check` CHECK ("type" IN ('app', 'package')),
  CONSTRAINT `App_language_check` CHECK ("language" IN ('typescript', 'go', 'python'))
);
INSERT INTO `new_App` (`type`, `language`, `name`, `location`, `project_id`)
SELECT 'app',
  CASE `name` WHEN 'mcpikon' THEN 'go' ELSE 'typescript' END,
  CASE `name` WHEN 'app-kaltra-fi' THEN '@kaltra/app-kaltra-fi'
              WHEN 'mcpikon' THEN 'github.com/verikono/mcpikon'
              ELSE `name` END,
  CASE `name` WHEN 'app-kaltra-fi' THEN 'apps/app-kaltra-fi'
              WHEN 'mcpikon' THEN '.'
              ELSE '' END,
  `project_id`
FROM `App`;

-- the child tables must carry the parent's new key.
CREATE TABLE `new_AppSkill` (
  `app_language` text NOT NULL,
  `app_name` text NOT NULL,
  `skillId` text NOT NULL,
  `priority` text NOT NULL DEFAULT 'NORMAL',
  PRIMARY KEY (`app_language`, `app_name`, `skillId`),
  CONSTRAINT `AppSkill_app_fkey` FOREIGN KEY (`app_language`, `app_name`) REFERENCES `App` (`language`, `name`) ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT `AppSkill_skillId_fkey` FOREIGN KEY (`skillId`) REFERENCES `Skill` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT `AppSkill_priority_check` CHECK ("priority" IN ('LOW', 'NORMAL', 'HIGH', 'ESSENTIAL'))
);
INSERT INTO `new_AppSkill` (`app_language`, `app_name`, `skillId`, `priority`)
SELECT
  CASE `app_id` WHEN 'mcpikon' THEN 'go' ELSE 'typescript' END,
  CASE `app_id` WHEN 'app-kaltra-fi' THEN '@kaltra/app-kaltra-fi'
                WHEN 'mcpikon' THEN 'github.com/verikono/mcpikon'
                ELSE `app_id` END,
  `skillId`, `priority`
FROM `AppSkill`;

CREATE TABLE `new_StyleToken` (
  `style` text NOT NULL,
  `value` text NOT NULL,
  `mode` text NOT NULL,
  `type` text NOT NULL,
  `token` text NOT NULL,
  `app_language` text NOT NULL,
  `app_name` text NOT NULL,
  `variable` text NULL,
  PRIMARY KEY (`app_language`, `app_name`, `mode`, `token`, `style`),
  CONSTRAINT `StyleToken_app_fkey` FOREIGN KEY (`app_language`, `app_name`) REFERENCES `App` (`language`, `name`) ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT `StyleToken_mode_check` CHECK ("mode" IN ('light', 'dark', 'primitive')),
  CONSTRAINT `StyleToken_type_check` CHECK ("type" IN ('color', 'spatial', 'size', 'dimension', 'font', 'shadow', 'zindex', 'gradient', 'weight', 'css')),
  CONSTRAINT `StyleToken_token_check` CHECK ("token" IN ('semantic', 'applied', 'primitive'))
);
INSERT INTO `new_StyleToken` (`style`, `value`, `mode`, `type`, `token`, `app_language`, `app_name`, `variable`)
SELECT `style`, `value`, `mode`, `type`, `token`,
  CASE `app_id` WHEN 'mcpikon' THEN 'go' ELSE 'typescript' END,
  CASE `app_id` WHEN 'app-kaltra-fi' THEN '@kaltra/app-kaltra-fi'
                WHEN 'mcpikon' THEN 'github.com/verikono/mcpikon'
                ELSE `app_id` END,
  `variable`
FROM `StyleToken`;

DROP TABLE `AppSkill`;
DROP TABLE `StyleToken`;
ALTER TABLE `new_AppSkill` RENAME TO `AppSkill`;
ALTER TABLE `new_StyleToken` RENAME TO `StyleToken`;

CREATE TABLE `AppAgent` (
  `app_language` text NOT NULL,
  `app_name` text NOT NULL,
  `agentId` text NOT NULL,
  PRIMARY KEY (`app_language`, `app_name`, `agentId`),
  CONSTRAINT `AppAgent_app_fkey` FOREIGN KEY (`app_language`, `app_name`) REFERENCES `App` (`language`, `name`) ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT `AppAgent_agentId_fkey` FOREIGN KEY (`agentId`) REFERENCES `Agent` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION
);

PRAGMA foreign_keys = on;

-- +goose Down
PRAGMA foreign_keys = off;
DROP TABLE `AppAgent`;

CREATE TABLE `new_App` (
  `name` text NOT NULL,
  `project_id` text NOT NULL,
  PRIMARY KEY (`name`),
  CONSTRAINT `App_project_id_fkey` FOREIGN KEY (`project_id`) REFERENCES `Project` (`name`) ON UPDATE NO ACTION ON DELETE NO ACTION
);
INSERT INTO `new_App` (`name`, `project_id`)
SELECT CONCAT(`language`, ':', `name`), `project_id` FROM `App`;

CREATE TABLE `new_AppSkill` (
  `app_id` text NOT NULL,
  `skillId` text NOT NULL,
  `priority` text NOT NULL DEFAULT 'NORMAL',
  PRIMARY KEY (`app_id`, `skillId`),
  CONSTRAINT `AppSkill_app_id_fkey` FOREIGN KEY (`app_id`) REFERENCES `App` (`name`) ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT `AppSkill_skillId_fkey` FOREIGN KEY (`skillId`) REFERENCES `Skill` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT `AppSkill_priority_check` CHECK ("priority" IN ('LOW', 'NORMAL', 'HIGH', 'ESSENTIAL'))
);
INSERT INTO `new_AppSkill` (`app_id`, `skillId`, `priority`)
SELECT CONCAT(`app_language`, ':', `app_name`), `skillId`, `priority` FROM `AppSkill`;

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
INSERT INTO `new_StyleToken` (`style`, `value`, `mode`, `type`, `token`, `app_id`, `variable`)
SELECT `style`, `value`, `mode`, `type`, `token`, CONCAT(`app_language`, ':', `app_name`), `variable` FROM `StyleToken`;

DROP TABLE `AppSkill`;
DROP TABLE `StyleToken`;
DROP TABLE `App`;
ALTER TABLE `new_App` RENAME TO `App`;
ALTER TABLE `new_AppSkill` RENAME TO `AppSkill`;
ALTER TABLE `new_StyleToken` RENAME TO `StyleToken`;
CREATE UNIQUE INDEX `App_name_key` ON `App` (`name`);

CREATE TABLE `new_Project` (
  `name` text NOT NULL,
  `skillsPath` text NOT NULL,
  `agentsPath` text NOT NULL,
  PRIMARY KEY (`name`)
);
INSERT INTO `new_Project` (`name`, `skillsPath`, `agentsPath`) SELECT `name`, `skillsPath`, `agentsPath` FROM `Project`;
DROP TABLE `Project`;
ALTER TABLE `new_Project` RENAME TO `Project`;
CREATE UNIQUE INDEX `Project_name_key` ON `Project` (`name`);
PRAGMA foreign_keys = on;
