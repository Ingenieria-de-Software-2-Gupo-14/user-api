package models

type ResponseUser struct {
	User User `json:"data"`
}

type ResponseUsers struct {
	Users []User `json:"data"`
}
