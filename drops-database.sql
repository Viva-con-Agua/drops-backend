-- MySQL Workbench Forward Engineering

SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0;
SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0;
SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='TRADITIONAL,ALLOW_INVALID_DATES';

-- -----------------------------------------------------
-- Schema drops
-- -----------------------------------------------------

-- -----------------------------------------------------
-- Schema drops
-- -----------------------------------------------------
CREATE SCHEMA IF NOT EXISTS `drops` DEFAULT CHARACTER SET utf8mb4 ;
USE `drops` ;

-- -----------------------------------------------------
-- Table `drops`.`city`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `drops`.`city` (
  `id` BIGINT(20) NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(64) NOT NULL,
  `country` VARCHAR(64) NOT NULL,
  `state` VARCHAR(64) NOT NULL,
  `google_id` VARCHAR(255) NOT NULL,
  `updated` BIGINT(20) NOT NULL,
  `created` BIGINT(20) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `google_id_UNIQUE` (`google_id` ASC))
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `drops`.`entity`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `drops`.`entity` (
  `id` BIGINT(20) NOT NULL AUTO_INCREMENT,
  `public_id` BINARY(16) NOT NULL,
  `name` VARCHAR(64) NOT NULL,
  `email` VARCHAR(64) NULL,
  `phone` VARCHAR(64) NULL,
  `fax` VARCHAR(64) NULL,
  `type` ENUM('organisation', 'crew', 'workgroup') NOT NULL,
  `abbreviation` VARCHAR(8) NOT NULL,
  `updated` BIGINT(20) NOT NULL,
  `created` BIGINT(20) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `public_id_UNIQUE` (`public_id` ASC))
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `drops`.`profile`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `drops`.`profile` (
  `id` BIGINT(20) NOT NULL AUTO_INCREMENT,
  `public_id` BINARY(16) NOT NULL,
  `firstname` VARCHAR(64) NOT NULL,
  `lastname` VARCHAR(64) NOT NULL,
  `email` VARCHAR(64) NOT NULL,
  `mobile` VARCHAR(64) NOT NULL,
  `birthdate` BIGINT(20) NOT NULL,
  `sex` ENUM('divers', 'male', 'female', 'none') NOT NULL DEFAULT 'none',
  `updated` BIGINT(20) NOT NULL,
  `created` BIGINT(20) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `email_UNIQUE` (`email` ASC),
  UNIQUE INDEX `public_id_UNIQUE` (`public_id` ASC))
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `drops`.`address`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `drops`.`address` (
  `id` BIGINT(20) NOT NULL AUTO_INCREMENT,
  `public_id` BINARY(16) NOT NULL,
  `street` VARCHAR(64) NOT NULL,
  `additional` VARCHAR(64) NOT NULL,
  `zip` VARCHAR(16) NOT NULL,
  `city` VARCHAR(64) NOT NULL,
  `country` VARCHAR(64) NOT NULL,
  `google_id` VARCHAR(255) NOT NULL,
  `updated` BIGINT(20) NOT NULL,
  `created` BIGINT(20) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `public_id_UNIQUE` (`public_id` ASC))
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `drops`.`avatar`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `drops`.`avatar` (
  `id` BIGINT(20) NOT NULL AUTO_INCREMENT,
  `type` VARCHAR(64) NOT NULL,
  `url` VARCHAR(64) NOT NULL,
  `updated` BIGINT(20) NOT NULL,
  `created` BIGINT(20) NOT NULL,
  `profile_id` BIGINT(20) NOT NULL,
  PRIMARY KEY (`id`, `profile_id`),
  INDEX `fk_avatar_profile_idx` (`profile_id` ASC),
  CONSTRAINT `fk_avatar_profile`
    FOREIGN KEY (`profile_id`)
    REFERENCES `drops`.`profile` (`id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION)
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `drops`.`entity_has_city`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `drops`.`entity_has_city` (
  `entity_id` BIGINT(20) NOT NULL,
  `city_id` BIGINT(20) NOT NULL,
  PRIMARY KEY (`entity_id`, `city_id`),
  INDEX `fk_entity_has_city_city1_idx` (`city_id` ASC),
  INDEX `fk_entity_has_city_entity1_idx` (`entity_id` ASC),
  CONSTRAINT `fk_entity_has_city_entity1`
    FOREIGN KEY (`entity_id`)
    REFERENCES `drops`.`entity` (`id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_entity_has_city_city1`
    FOREIGN KEY (`city_id`)
    REFERENCES `drops`.`city` (`id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION)
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `drops`.`profile_has_address`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `drops`.`profile_has_address` (
  `profile_id` BIGINT(20) NOT NULL,
  `address_id` BIGINT(20) NOT NULL,
  PRIMARY KEY (`profile_id`, `address_id`),
  INDEX `fk_profile_has_address_address1_idx` (`address_id` ASC),
  CONSTRAINT `fk_profile_has_address_profile1`
    FOREIGN KEY (`profile_id`)
    REFERENCES `drops`.`profile` (`id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_profile_has_address_address1`
    FOREIGN KEY (`address_id`)
    REFERENCES `drops`.`address` (`id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION)
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `drops`.`profile_has_entity`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `drops`.`profile_has_entity` (
  `profile_id` BIGINT(20) NOT NULL,
  `entity_id` BIGINT(20) NOT NULL,
  PRIMARY KEY (`profile_id`, `entity_id`),
  INDEX `fk_profile_has_entity_entity1_idx` (`entity_id` ASC),
  CONSTRAINT `fk_profile_has_entity_profile1`
    FOREIGN KEY (`profile_id`)
    REFERENCES `drops`.`profile` (`id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_profile_has_entity_entity1`
    FOREIGN KEY (`entity_id`)
    REFERENCES `drops`.`entity` (`id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION)
ENGINE = InnoDB;


SET SQL_MODE=@OLD_SQL_MODE;
SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS;
SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS;
