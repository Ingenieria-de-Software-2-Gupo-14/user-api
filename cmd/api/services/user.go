package services

import . "ing-soft-2-tp1/cmd/api/models"

// CreateUser creates and returns a User Struct
func CreateUser(id int, title string, description string) User {
	user := User{
		Id:          id,
		Title:       title,
		Description: description}
	return user
}
