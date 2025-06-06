-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS password_reset (
    id SERIAL PRIMARY KEY,
    user_id int,
    email VARCHAR(255),
    token VARCHAR(6) NOT NULL UNIQUE,
    token_expiration TIMESTAMP,
    used BOOLEAN DEFAULT false,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
    );
-- +goose StatementEnd