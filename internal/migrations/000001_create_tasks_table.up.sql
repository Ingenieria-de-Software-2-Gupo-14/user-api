-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL, -- TODO: check if we need username and if it should be unique
    name VARCHAR(50) NOT NULL,
    surname VARCHAR(50) NOT NULL,
    password VARCHAR(50) NOT NULL,
    UNIQUE email VARCHAR(240) NOT NULL, -- TODO: make it larger and unique (DONE)
    location VARCHAR(50) NOT NULL,
    admin BOOLEAN NOT NULL,
    blocked_user BOOLEAN NOT NULL,
    profile_photo INTEGER, -- TODO: Modify to accept URLS (VARCHARS) instead of integers
    description VARCHAR(240) NOT NULL,
    name_privacy BOOLEAN DEFAULT false NOT NULL,
    surname_privacy BOOLEAN DEFAULT false NOT NULL,
    email_privacy BOOLEAN DEFAULT false NOT NULL,
    location_privacy BOOLEAN DEFAULT false NOT NULL,
    description_privacy BOOLEAN DEFAULT false NOT NULL
    );
/*CREATE TABLE IF NOT EXISTS users_privacy{
    id SERIAL PRIMARY KEY,
    UNIQUE user_id INT NOT NULL,
    account BOOLEAN DEFAULT true,
    name BOOLEAN DEFAULT true,
    surname BOOLEAN DEFAULT true,
    email BOOLEAN DEFAULT true,
    location BOOLEAN DEFAULT true,
    description BOOLEAN DEFAULT true,
    FOREIGN KEY (userid) REFERENCES users(id)
    };*/
-- +goose StatementEnd
