-- +goose Up
-- +goose StatementBegin

ALTER TABLE users
    DROP COLUMN username,
    ALTER COLUMN email TYPE VARCHAR(255),
    ADD CONSTRAINT users_email_unique UNIQUE (email),
    ALTER COLUMN password TYPE VARCHAR(60),
    ALTER COLUMN admin SET DEFAULT FALSE,
    DROP COLUMN blocked_user,
    ALTER COLUMN profile_photo TYPE VARCHAR(255),
    ALTER COLUMN profile_photo DROP NOT NULL, -- Assuming profile_photo can be null based on original CREATE statement
    ADD COLUMN phone VARCHAR(50),
    ADD COLUMN created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    ADD COLUMN updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP;



-- Function to update updated_at column
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
   NEW.updated_at = NOW();
   RETURN NEW;
END;
$$ language 'plpgsql';

-- Trigger to execute the function before update on users table
CREATE TRIGGER trigger_users_updated_at
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

CREATE TABLE IF NOT EXISTS login_attempts (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    ip_address VARCHAR(50) NOT NULL,
    user_agent VARCHAR(255) NOT NULL,
    user_id INTEGER NOT NULL, -- Keep as INTEGER to match users.id (SERIAL)
    successful BOOLEAN NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_login_attempts_user_id ON login_attempts(user_id);

CREATE TABLE IF NOT EXISTS blocked_users (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    blocked_until TIMESTAMP, -- if null, it means the user is blocked forever
    reason VARCHAR(255) NOT NULL,
    blocker_id INTEGER, -- This can be null if the user is blocked by the system
    blocked_user_id INTEGER NOT NULL,
    FOREIGN KEY (blocker_id) REFERENCES users(id) ON DELETE SET NULL,
    FOREIGN KEY (blocked_user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_blocked_users_blocked_user_id ON blocked_users(blocked_user_id);

-- +goose StatementEnd
