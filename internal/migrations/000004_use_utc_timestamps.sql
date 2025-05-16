-- +goose Up
-- +goose StatementBegin

-- Set timezone to UTC for database session
SET timezone = 'UTC';

-- Modify existing timestamp columns to use UTC timezone
ALTER TABLE users
    ALTER COLUMN created_at TYPE TIMESTAMP WITH TIME ZONE,
    ALTER COLUMN updated_at TYPE TIMESTAMP WITH TIME ZONE;

ALTER TABLE login_attempts
    ALTER COLUMN created_at TYPE TIMESTAMP WITH TIME ZONE;

ALTER TABLE blocked_users
    ALTER COLUMN created_at TYPE TIMESTAMP WITH TIME ZONE,
    ALTER COLUMN blocked_until TYPE TIMESTAMP WITH TIME ZONE;

-- Update the function to use UTC timezone
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
   NEW.updated_at = NOW() AT TIME ZONE 'UTC';
   RETURN NEW;
END;
$$ language 'plpgsql';

-- +goose StatementEnd
