CREATE DATABASE jobhunter;
USE jobhunter;
CREATE TABLE IF NOT EXISTS `role` (
  `role_id` INT NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(15) CHARACTER SET 'utf8mb4' NOT NULL,
  PRIMARY KEY (`role_id`),
  UNIQUE INDEX `Rolecol_UNIQUE` (`name` ASC) );


CREATE TABLE IF NOT EXISTS `user` (
  `user_id` INT NOT NULL AUTO_INCREMENT,
  `login` VARCHAR(40) CHARACTER SET 'utf8mb4' NOT NULL,
  `email` VARCHAR(255) CHARACTER SET 'utf8mb4' NOT NULL,
  `date_created` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `password` BINARY(60) NOT NULL,
  `last_check` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `full_name` VARCHAR(150) CHARACTER SET 'utf8mb4' NOT NULL,
  PRIMARY KEY (`user_id`),
  UNIQUE INDEX `login_UNIQUE` (`login` ASC) ,
  UNIQUE INDEX `email_UNIQUE` (`email` ASC) );


CREATE TABLE IF NOT EXISTS `user_has_role` (
  `fk_user_id` INT NOT NULL,
  `fk_role_id` INT NOT NULL,
  PRIMARY KEY (`fk_user_id`, `fk_role_id`),
  INDEX `fk_User_has_Role_Role1_idx` (`fk_role_id` ASC) ,
  INDEX `fk_User_has_Role_User_idx` (`fk_user_id` ASC) ,
  CONSTRAINT `fk_User_has_Role_User`
    FOREIGN KEY (`fk_user_id`)
    REFERENCES `user` (`user_id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_User_has_Role_Role1`
    FOREIGN KEY (`fk_role_id`)
    REFERENCES `role` (`role_id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION);