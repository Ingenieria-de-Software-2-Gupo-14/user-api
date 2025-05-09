-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS verification (
    id SERIAL PRIMARY KEY,
    email TYPE VARCHAR(255) UNIQUE,
    password TYPE VARCHAR(60),
    name VARCHAR(50) NOT NULL,
    surname VARCHAR(50) NOT NULL,
    verification_pin VARCHAR(6) NOT NULL,
    pin_expiration TYPE TIMESTAMP WITH TIME ZONE,
-- +goose StatementEnd