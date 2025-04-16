package models

type CreateUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type User struct {
	Id           int    `json:"id"`
	Username     string `json:"username"`
	Name         string `json:"name"`
	Surname      string `json:"surname"`
	Password     string `json:"password"`
	Email        string `json:"email"`
	Location     string `json:"location"`
	Admin        bool   `json:"admin"`
	BlockedUser  bool   `json:"blocked-user"`
	ProfilePhoto int    `json:"profile-photo"` //Foreign key to mongodb with profile pic
	Description  string `json:"description"`
}
