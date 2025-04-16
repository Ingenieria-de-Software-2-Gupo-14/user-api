package services

import (
	. "ing-soft-2-tp1/internal/database"
	. "ing-soft-2-tp1/internal/models"
)

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
func RemoveUserFromDatabase(db *Database, id int) {
	db.DeleteUser(id)
}

func AddUserToDatabase(db *Database, user *User) {
	db.AddUser(user)
}

func GetUserFromDatabase(db *Database, id int) (user *User, ok error) {
	user, ok = db.GetUser(id)
	return user, ok
}

func GetUserFromDatabaseByEmailAndPassword(db *Database, email string, password string) (user *User, ok error) {
	user, ok = db.GetUserByEmailAndPassword(email, password)
	return user, ok
}

func GetAllUsersFromDatabase(db *Database) (users []User, err error) {
	return db.GetAllUsers()
}
