package database

import (
	"database/sql"
	"errors"
	_ "github.com/lib/pq"
	. "ing-soft-2-tp1/internal/models"
)

type Database struct {
	DB *sql.DB
}

// CreateDatabase creates and returns a database
func CreateDatabase(db *sql.DB) *Database {
	return &Database{DB: db}
}

// GetUser returns User corresponding to id and ok bool value, if ok true, the User was in the database, if ok false then the User wasn't in the database
func (db Database) GetUser(id int) (*User, error) {
	row := db.DB.QueryRow("SELECT * FROM users WHERE id = $1", id)

	var user User
	err := row.Scan(&user.Id, &user.Username, &user.Name, &user.Surname, &user.Password, &user.Email, &user.Location, &user.Admin, &user.BlockedUser, &user.ProfilePhoto, &user.Description)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found") //TODO Make custom error
		}
		return nil, err
	}

	return &user, nil
}

// GetAllUsers returns a slices containing all elements of the database, if the database is empty then it returns an empty slice
func (db Database) GetAllUsers() ([]User, error) {
	rows, err := db.DB.Query("SELECT id, username, name, surname,  password,email, location, admin, blocked_user, profile_photo,description FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.Id, &user.Username, &user.Name, &user.Surname, &user.Password, &user.Email, &user.Location, &user.Admin, &user.BlockedUser, &user.ProfilePhoto, &user.Description)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

// DeleteUser deletes a User from the database corresponding to the id
func (db Database) DeleteUser(id int) error {
	_, err := db.DB.Exec("DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return err
	}

	return nil
}

// AddUser adds an elements to the database
func (db Database) AddUser(user *User) error {
	_, err := db.DB.Exec("INSERT INTO users (username, name, surname,  password,email, location, admin, blocked_user, profile_photo,description) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)", &user.Username, &user.Name, &user.Surname, &user.Password, &user.Email, &user.Location, &user.Admin, &user.BlockedUser, &user.ProfilePhoto, &user.Description)
	if err != nil {
		return err
	}

	return nil
}

func (db Database) GetUserByEmailAndPassword(email string, password string) (*User, error) {
	println(email + password)
	row := db.DB.QueryRow("SELECT * FROM users WHERE (email ILIKE $1) AND (password ILKE $2)", email, password)
	var user User
	err := row.Scan(&user.Id, &user.Username, &user.Name, &user.Surname, &user.Password, &user.Email, &user.Location, &user.Admin, &user.BlockedUser, &user.ProfilePhoto, &user.Description)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found") //TODO Make custom error
		}
		return nil, err
	}

	return &user, nil
}

//id, username, name, surname,  password,email, location, admin, blocked_user, profile_photo,description
