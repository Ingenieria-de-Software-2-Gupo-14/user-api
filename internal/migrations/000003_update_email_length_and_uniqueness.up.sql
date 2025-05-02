-- +goose Up
-- +goose StatementBegin
ALTER TABLE Users
    ALTER COLUMN email TYPE VARCHAR(240);
ALTER TABLE Users
    ADD CONSTRAINT unique_email UNIQUE (email)
-- +goose StatementEnd
