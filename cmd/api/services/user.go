package services

import . "ing-soft-2-tp1/cmd/api/models"

// CreateUser creates and returns a User Struct
func CreateUser(id int, username string, password string) User {
	user := User{
		Id:          id,
		Username:    username,
		Name:        "",
		Surname:     "",
		Email:       "",
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
