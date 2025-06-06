-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
    ADD COLUMN exam_notification BOOLEAN DEFAULT true,
    ADD COLUMN homework_notification BOOLEAN DEFAULT true,
    ADD COLUMN social_notification BOOLEAN DEFAULT true;

UPDATE users
SET
    exam_notification = TRUE,
    homework_notification = TRUE,
    social_notification = TRUE;
-- +goose StatementEnd