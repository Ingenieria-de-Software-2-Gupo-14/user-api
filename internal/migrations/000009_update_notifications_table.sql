-- +goose Up
-- +goose StatementBegin
ALTER TABLE notifications RENAME COLUMN notification_text TO token;
ALTER TABLE notifications ADD CONSTRAINT unique_token UNIQUE (token);
-- +goose StatementEnd
