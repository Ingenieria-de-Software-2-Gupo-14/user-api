package models

type CreateUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Surname  string `json:"surname"`
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

type UserPrivacy struct {
	Account     bool `json:"account"`
	Name        bool `json:"name"`
	Surname     bool `json:"surname"`
	Email       bool `json:"email"`
	Location    bool `json:"location"`
	Description bool `json:"description"`
}
