SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0;
SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0;
SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='TRADITIONAL,ALLOW_INVALID_DATES';


-- -----------------------------------------------------
-- Table `users`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `users` ;

CREATE TABLE IF NOT EXISTS `users` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `username` VARCHAR(75) NULL,
  `created` DATETIME NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `username_UNIQUE` (`username` ASC))
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `user_passwords`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `user_passwords` ;

CREATE TABLE IF NOT EXISTS `user_passwords` (
  `user_id` INT UNSIGNED NOT NULL,
  `timestamp` BIGINT UNSIGNED NOT NULL,
  `password` VARCHAR(255) NULL,
  PRIMARY KEY (`user_id`, `timestamp`),
  CONSTRAINT `fk_user_passwords_user_id`
    FOREIGN KEY (`user_id`)
    REFERENCES `users` (`id`)
    ON DELETE RESTRICT
    ON UPDATE CASCADE)
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `mods`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `mods` ;

CREATE TABLE IF NOT EXISTS `mods` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `title` VARCHAR(75) NULL,
  `created` DATETIME NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `title_UNIQUE` (`title` ASC))
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `mod_users`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `mod_users` ;

CREATE TABLE IF NOT EXISTS `mod_users` (
  `mod_id` INT UNSIGNED NOT NULL,
  `user_id` INT UNSIGNED NOT NULL,
  `timestamp` BIGINT UNSIGNED NOT NULL,
  `role` VARCHAR(45) NULL,
  PRIMARY KEY (`mod_id`, `user_id`, `timestamp`),
  INDEX `fk_mod_users_user_id_idx` (`user_id` ASC),
  INDEX `roles` (`role` ASC),
  CONSTRAINT `fk_mod_users_mod_id`
    FOREIGN KEY (`mod_id`)
    REFERENCES `mods` (`id`)
    ON DELETE RESTRICT
    ON UPDATE CASCADE,
  CONSTRAINT `fk_mod_users_user_id`
    FOREIGN KEY (`user_id`)
    REFERENCES `users` (`id`)
    ON DELETE RESTRICT
    ON UPDATE CASCADE)
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `mod_metadata`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `mod_metadata` ;

CREATE TABLE IF NOT EXISTS `mod_metadata` (
  `mod_id` INT UNSIGNED NOT NULL,
  `timestamp` BIGINT UNSIGNED NOT NULL,
  `url` VARCHAR(255) NULL,
  `modid` VARCHAR(125) NULL,
  `version` INT UNSIGNED NULL,
  `date` VARCHAR(75) NULL,
  `save` TINYINT(1) NULL,
  `enabled` TINYINT(1) NULL,
  PRIMARY KEY (`mod_id`, `timestamp`),
  CONSTRAINT `fk_mod_metadata_mod_id`
    FOREIGN KEY (`mod_id`)
    REFERENCES `mods` (`id`)
    ON DELETE RESTRICT
    ON UPDATE CASCADE)
ENGINE = InnoDB;


SET SQL_MODE=@OLD_SQL_MODE;
SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS;
SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS;
