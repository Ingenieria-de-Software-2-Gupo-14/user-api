-- +goose Up
-- +goose StatementBegin
UPDATE users SET verified = true;
-- +goose StatementEnd
