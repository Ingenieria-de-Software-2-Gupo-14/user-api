-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL,
    name VARCHAR(50) NOT NULL,
    surname VARCHAR(50) NOT NULL,
    password VARCHAR(50) NOT NULL,
    email VARCHAR(50) NOT NULL,
    location VARCHAR(50) NOT NULL,
    admin BOOLEAN NOT NULL,
    blocked_user BOOLEAN NOT NULL,
    profile_photo INTEGER, -- TODO: ADD  AT THE END AFTER MONGODB IS DONE FOREIGN KEY (profile_photo) REFERENCES photos(id)
    description VARCHAR(240) NOT NULL
    );
-- +goose StatementEnd
