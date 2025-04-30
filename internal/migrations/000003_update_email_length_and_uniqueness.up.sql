-- +goose Up
-- +goose StatementBegin
ALTER TABLE Users
    MODIFY COLUMN email VARCHAR(50) NOT NULL;
    ADD CONSTRAINT unique_email UNIQUE (email)
-- +goose StatementEnd
