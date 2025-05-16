-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ADD COLUMN role VARCHAR(10) NOT NULL DEFAULT 'student';
ALTER TABLE users ADD COLUMN verified BOOLEAN NOT NULL DEFAULT false;


UPDATE users SET role = 'admin' WHERE admin = true;


ALTER TABLE users DROP COLUMN admin;
ALTER TABLE users DROP COLUMN IF EXISTS phone;

DROP TABLE IF EXISTS verification;

CREATE TABLE IF NOT EXISTS verification (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL UNIQUE,
    user_email VARCHAR(255) NOT NULL UNIQUE,
    verification_pin VARCHAR(6) NOT NULL,
    pin_expiration TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_verification_user_id ON verification(user_id);
CREATE INDEX IF NOT EXISTS idx_verification_user_email ON verification(user_email);

-- +goose StatementEnd
