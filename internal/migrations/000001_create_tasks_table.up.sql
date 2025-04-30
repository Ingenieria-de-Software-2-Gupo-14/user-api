-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL, -- TODO: check if we need username and if it should be unique
    name VARCHAR(50) NOT NULL,
    surname VARCHAR(50) NOT NULL,
    password VARCHAR(50) NOT NULL,
    UNIQUE email VARCHAR(50) NOT NULL, -- TODO: make it larger and unique
    location VARCHAR(50) NOT NULL,
    admin BOOLEAN NOT NULL,
    blocked_user BOOLEAN NOT NULL,
    profile_photo INTEGER, -- TODO: ADD  AT THE END AFTER MONGODB IS DONE FOREIGN KEY (profile_photo) REFERENCES photos(id)
    description VARCHAR(240) NOT NULL
    );
CREATE TABLE IF NOT EXISTS users_privacy{
    id SERIAL PRIMARY KEY,
    userid INT,
    account BOOLEAN NOT NULL,
    name BOOLEAN NOT NULL,
    surname BOOLEAN NOT NULL,
    email BOOLEAN NOT NULL,
    location BOOLEAN NOT NULL,
    description BOOLEAN NOT NULL,
    FOREIGN KEY (userid) REFERENCES users(id)
    };
-- +goose StatementEnd
