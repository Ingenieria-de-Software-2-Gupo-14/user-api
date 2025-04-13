package services

import (
	. "ing-soft-2-tp1/internal/database"
	. "ing-soft-2-tp1/internal/models"
)

// CreateUser creates and returns a User Struct
func CreateUser(id int, email string, password string) User {
	user := User{
		Id:          id,
		Username:    email,
		Name:        "",
		Surname:     "",
		Email:       email,
		Password:    password,
		Description: "",
	}
	return user
}

func CreateAdminUser(id int, username string, password string) User {
	admin := CreateUser(id, username, password)
	admin.Admin = true
	return admin
}

// RemoveUserFromDatabase removes user from database
func RemoveUserFromDatabase(db *Database[User], id int) {
	db.DeleteUser(id)
}

func AddUserToDatabase(db *Database[User], user User) {
	db.AddUser(user)
}

func GetUserFromDatabase(db *Database[User], id int) (user User, ok bool) {
	user, ok = db.GetUser(id)
	return user, ok
}

func GetAllUsersFromDatabase(db *Database[User]) (users []User) {
	return db.GetAllUsers()
}
