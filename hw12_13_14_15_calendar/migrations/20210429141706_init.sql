-- +goose Up
-- +goose StatementBegin
CREATE TABLE event (
	id INTEGER AUTO_INCREMENT,
	title TEXT NOT NULL,
	start TIMESTAMP NOT NULL,
	stop TIMESTAMP NOT NULL,
	description TEXT NULL,
	user_id INTEGER NOT NULL,
	notification INTEGER NULL DEFAULT NULL,
	PRIMARY KEY(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE event;
-- +goose StatementEnd
