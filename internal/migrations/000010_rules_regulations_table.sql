-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS rules (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    effective_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    application_condition TEXT NOT NULL
    );
CREATE TABLE IF NOT EXISTS rules_audit (
    id SERIAL PRIMARY KEY,
    rule_id INTEGER NULL,
    user_id INTEGER NULL,
    modification_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    nature_of_modification TEXT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL,
    FOREIGN KEY (rule_id) REFERENCES rules(id) ON DELETE SET NULL
    );
-- +goose StatementEnd
