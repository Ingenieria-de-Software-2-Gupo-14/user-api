package models

type CreateUserRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type User struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}
