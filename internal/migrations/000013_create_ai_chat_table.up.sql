-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS ai_chat (
    id SERIAL PRIMARY KEY,
    user_id int NOT NULL,
    sender TEXT CHECK (sender IN ('user', 'assistant', 'system')) NOT NULL,
    message TEXT NOT NULL,
    time_sent TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    rating int NOT NULL,
    feedback TEXT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
    );
-- +goose StatementEnd
-- +goose Down