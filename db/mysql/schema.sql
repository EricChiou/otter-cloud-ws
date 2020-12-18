CREATE TABLE IF NOT EXISTS `role` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `code` VARCHAR(32) NOT NULL,
  `name` VARCHAR(32) NOT NULL,
  `lv` INT NOT NULL DEFAULT 1,
  `sort_no` INT NOT NULL DEFAULT 1,
  `enable` TINYINT NOT NULL DEFAULT 1,
  `created_date` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_date` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `code_UNIQUE` (`code` ASC) VISIBLE
);

CREATE TABLE IF NOT EXISTS `acl` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `code` VARCHAR(32) NOT NULL,
  `name` VARCHAR(32) NOT NULL,
  `type` VARCHAR(32) NOT NULL,
  `lv` INT NOT NULL DEFAULT 1,
  `sort_no` INT NOT NULL DEFAULT 1,
  `enable` TINYINT NOT NULL DEFAULT 1,
  `created_date` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_date` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `code_UNIQUE` (`code` ASC) VISIBLE
);

CREATE TABLE IF NOT EXISTS `user` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `acc` VARCHAR(256) NOT NULL,
  `pwd` VARCHAR(64) NOT NULL,
  `name` VARCHAR(64) NOT NULL,
  `role_code` VARCHAR(32) NOT NULL DEFAULT 'normal',
  `status` VARCHAR(32) NOT NULL DEFAULT 'inactive',
  `created_date` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_date` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  INDEX `user_role_code_FK_idx` (`role_code` ASC) INVISIBLE,
  UNIQUE INDEX `acc_UNIQUE` (`acc` ASC) VISIBLE,
  CONSTRAINT `user_role_code_FK`
    FOREIGN KEY (`role_code`)
    REFERENCES `role` (`code`)
    ON DELETE NO ACTION
    ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS `role_acl` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `role_code` VARCHAR(32) NOT NULL,
  `acl_code` VARCHAR(32) NOT NULL,
  `created_date` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_date` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  INDEX `role_acl_role_code_FK_idx` (`role_code` ASC) INVISIBLE,
  INDEX `role_acl_acl_code_FK_idx` (`acl_code` ASC) INVISIBLE,
  UNIQUE INDEX `role_acl_UNIQUE` (`role_code` ASC, `acl_code` ASC) VISIBLE,
  CONSTRAINT `role_acl_role_code_FK`
    FOREIGN KEY (`role_code`)
    REFERENCES `role` (`code`)
    ON DELETE NO ACTION
    ON UPDATE CASCADE,
  CONSTRAINT `role_acl_acl_code_FK`
    FOREIGN KEY (`acl_code`)
    REFERENCES `acl` (`code`)
    ON DELETE NO ACTION
    ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS `codemap` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `type` VARCHAR(32) NOT NULL,
  `code` VARCHAR(32) NOT NULL,
  `name` VARCHAR(32) NOT NULL,
  `sort_no` INT NOT NULL DEFAULT 1,
  `enable` TINYINT NOT NULL DEFAULT 1,
  `created_date` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_date` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `code_UNIQUE` (`code` ASC) VISIBLE
);

INSERT INTO `role` ( `code`, `name` ) VALUES ( 'manager', '管理員' );
INSERT INTO `role` ( `code`, `name` ) VALUES ( 'normal', '一般' );

INSERT INTO `acl` ( `code`, `name`, `type` ) VALUES ( 'addCodemap', '新增代碼', 'codemap' );
INSERT INTO `acl` ( `code`, `name`, `type` ) VALUES ( 'updateCodemap', '修改代碼', 'codemap' );
INSERT INTO `acl` ( `code`, `name`, `type` ) VALUES ( 'deleteCodemap', '刪除代碼', 'codemap' );
INSERT INTO `acl` ( `code`, `name`, `type` ) VALUES ( 'updateUser', '修改使用者', 'user' );
INSERT INTO `acl` ( `code`, `name`, `type` ) VALUES ( 'deleteUser', '刪除使用者', 'user' );

INSERT INTO `role_acl` ( `role_code`, `acl_code` ) VALUES ( 'manager', 'addCodemap' );
INSERT INTO `role_acl` ( `role_code`, `acl_code` ) VALUES ( 'manager', 'updateCodemap' );
INSERT INTO `role_acl` ( `role_code`, `acl_code` ) VALUES ( 'manager', 'deleteCodemap' );
INSERT INTO `role_acl` ( `role_code`, `acl_code` ) VALUES ( 'manager', 'updateUser' );
INSERT INTO `role_acl` ( `role_code`, `acl_code` ) VALUES ( 'manager', 'deleteUser' );