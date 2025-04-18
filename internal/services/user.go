package services

import (
	. "ing-soft-2-tp1/internal/models"
)

type Database interface {
	GetUser(id int) (*User, error)
	GetAllUsers() ([]User, error)
	DeleteUser(id int) error
	AddUser(user *User) (int, error)
	GetUserByEmailAndPassword(email string, password string) (*User, error)
	ContainsUserByEmail(email string) bool
	ModifyUser(user *User) error
	ClearDb() error
}

// CreateUser creates and returns a User Struct
func CreateUser(id int, email string, password string) User {
	user := User{
		Id:       id,
		Username: email,
		Email:    email,
		Password: password,
	}
	return user
}

func CreateAdminUser(id int, username string, password string) User {
	admin := CreateUser(id, username, password)
	admin.Admin = true
	return admin
}

// RemoveUserFromDatabase removes user from database
func RemoveUserFromDatabase(db Database, id int) {
	db.DeleteUser(id)
}

func AddUserToDatabase(db Database, user *User) (int, error) {
	return db.AddUser(user)
}

func GetUserFromDatabase(db Database, id int) (user *User, ok error) {
	user, ok = db.GetUser(id)
	return user, ok
}

func GetUserFromDatabaseByEmailAndPassword(db Database, email string, password string) (user *User, ok error) {
	user, ok = db.GetUserByEmailAndPassword(email, password)
	return user, ok
}

func GetAllUsersFromDatabase(db Database) (users []User, err error) {
	return db.GetAllUsers()
}

func ContainsUserByEmail(db Database, email string) bool {
	return db.ContainsUserByEmail(email)
}

func ModifyUser(db Database, user *User) error {
	return db.ModifyUser(user)

}

func ClearDb(db Database) error {
	return db.ClearDb()
}
