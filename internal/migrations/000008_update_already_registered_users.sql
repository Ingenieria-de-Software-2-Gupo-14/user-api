-- +goose Up
-- +goose StatementBegin
UPDATE users SET verfied = true;
-- +goose StatementEnd
