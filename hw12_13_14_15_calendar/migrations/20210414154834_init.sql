-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE `event` (
	`id` INT(11) NOT NULL AUTO_INCREMENT,
	`title` TEXT NOT NULL,
	`start` TIMESTAMP NOT NULL,
	`stop` TIMESTAMP NOT NULL,
	`description` TEXT NULL,
	`user_id` INT(11) NOT NULL,
	`notification` BIGINT(20) NULL DEFAULT NULL,
	PRIMARY KEY (`id`)
)
COLLATE='utf8_general_ci'
ENGINE=MyISAM
;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE `event`;
