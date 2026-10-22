-- +goose Up
PRAGMA foreign_keys = off;
CREATE TABLE `new_SkillDependency` (
  `dependencyId` text NOT NULL,
  `dependentId` text NOT NULL,
  `agent` text NOT NULL DEFAULT '',
  `recursive` boolean NOT NULL DEFAULT true,
  PRIMARY KEY (`dependencyId`, `dependentId`, `agent`),
  CONSTRAINT `SkillDependency_dependencyId_fkey` FOREIGN KEY (`dependencyId`) REFERENCES `Skill` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT `SkillDependency_dependentId_fkey` FOREIGN KEY (`dependentId`) REFERENCES `Skill` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION
);
INSERT INTO `new_SkillDependency` (`dependencyId`, `dependentId`, `agent`, `recursive`) SELECT `dependencyId`, `dependentId`, '', `recursive` FROM `SkillDependency`;
DROP TABLE `SkillDependency`;
ALTER TABLE `new_SkillDependency` RENAME TO `SkillDependency`;
PRAGMA foreign_keys = on;

-- +goose Down
PRAGMA foreign_keys = off;
CREATE TABLE `new_SkillDependency` (
  `dependencyId` text NOT NULL,
  `dependentId` text NOT NULL,
  `recursive` boolean NOT NULL DEFAULT true,
  PRIMARY KEY (`dependencyId`, `dependentId`),
  CONSTRAINT `SkillDependency_dependencyId_fkey` FOREIGN KEY (`dependencyId`) REFERENCES `Skill` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT `SkillDependency_dependentId_fkey` FOREIGN KEY (`dependentId`) REFERENCES `Skill` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION
);
INSERT INTO `new_SkillDependency` (`dependencyId`, `dependentId`, `recursive`) SELECT `dependencyId`, `dependentId`, `recursive` FROM `SkillDependency` WHERE `agent` = '';
DROP TABLE `SkillDependency`;
ALTER TABLE `new_SkillDependency` RENAME TO `SkillDependency`;
PRAGMA foreign_keys = on;
