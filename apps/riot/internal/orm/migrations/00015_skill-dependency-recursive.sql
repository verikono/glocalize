-- +goose Up
PRAGMA foreign_keys = off;
ALTER TABLE `SkillDependency` ADD COLUMN `recursive` boolean NOT NULL DEFAULT true;
PRAGMA foreign_keys = on;

-- +goose Down
PRAGMA foreign_keys = off;
CREATE TABLE `new_SkillDependency` (
  `dependencyId` text NOT NULL,
  `dependentId` text NOT NULL,
  PRIMARY KEY (`dependencyId`, `dependentId`),
  CONSTRAINT `SkillDependency_dependencyId_fkey` FOREIGN KEY (`dependencyId`) REFERENCES `Skill` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT `SkillDependency_dependentId_fkey` FOREIGN KEY (`dependentId`) REFERENCES `Skill` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION
);
INSERT INTO `new_SkillDependency` (`dependencyId`, `dependentId`) SELECT `dependencyId`, `dependentId` FROM `SkillDependency`;
DROP TABLE `SkillDependency`;
ALTER TABLE `new_SkillDependency` RENAME TO `SkillDependency`;
PRAGMA foreign_keys = on;
