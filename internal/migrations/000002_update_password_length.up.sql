-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
ALTER COLUMN password TYPE VARCHAR(200);
-- +goose StatementEnd
