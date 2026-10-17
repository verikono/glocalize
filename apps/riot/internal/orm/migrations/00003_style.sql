-- +goose Up
PRAGMA foreign_keys = off;
CREATE TABLE `SkillDependency` (
  `dependencyId` text NOT NULL,
  `dependentId` text NOT NULL,
  PRIMARY KEY (`dependencyId`, `dependentId`),
  CONSTRAINT `SkillDependency_dependencyId_fkey` FOREIGN KEY (`dependencyId`) REFERENCES `Skill` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT `SkillDependency_dependentId_fkey` FOREIGN KEY (`dependentId`) REFERENCES `Skill` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION
);
CREATE TABLE `ProjectSkill` (
  `projectName` text NOT NULL,
  `skillId` text NOT NULL,
  PRIMARY KEY (`projectName`, `skillId`),
  CONSTRAINT `ProjectSkill_projectName_fkey` FOREIGN KEY (`projectName`) REFERENCES `Project` (`name`) ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT `ProjectSkill_skillId_fkey` FOREIGN KEY (`skillId`) REFERENCES `Skill` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION
);
CREATE TABLE `ProjectSpecialization` (
  `generalizesName` text NOT NULL,
  `specializesName` text NOT NULL,
  PRIMARY KEY (`generalizesName`, `specializesName`),
  CONSTRAINT `ProjectSpecialization_generalizesName_fkey` FOREIGN KEY (`generalizesName`) REFERENCES `Project` (`name`) ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT `ProjectSpecialization_specializesName_fkey` FOREIGN KEY (`specializesName`) REFERENCES `Project` (`name`) ON UPDATE NO ACTION ON DELETE NO ACTION
);
CREATE TABLE `StyleToken` (
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
CREATE TABLE `App` (
  `name` text NOT NULL,
  `project_id` text NOT NULL,
  PRIMARY KEY (`name`),
  CONSTRAINT `App_project_id_fkey` FOREIGN KEY (`project_id`) REFERENCES `Project` (`name`) ON UPDATE NO ACTION ON DELETE NO ACTION
);
CREATE UNIQUE INDEX `App_name_key` ON `App` (`name`);
PRAGMA foreign_keys = on;

-- +goose Down
PRAGMA foreign_keys = off;
DROP TABLE `SkillDependency`;
DROP TABLE `ProjectSkill`;
DROP TABLE `ProjectSpecialization`;
DROP TABLE `StyleToken`;
DROP TABLE `App`;
PRAGMA foreign_keys = on;
