package models

import "time"

type CreateUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=60"`
	Name     string `json:"name" binding:"required,min=3,max=60"`
	Surname  string `json:"surname" binding:"required,min=3,max=60"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type PasswordModifyRequest struct {
	Password string `json:"password" binding:"required,min=8,max=60"`
}

type LocationModifyRequest struct {
	Location string `json:"location" binding:"required"`
}

type User struct {
	Id           int       `json:"id" db:"id"`
	Name         string    `json:"name" db:"name"`
	Surname      string    `json:"surname" db:"surname"`
	Password     string    `json:"-" db:"password"`
	Email        string    `json:"email" db:"email"`
	Location     string    `json:"location,omitempty" db:"location"`
	Admin        bool      `json:"admin" db:"admin"`
	ProfilePhoto *string   `json:"profile_photo,omitempty" db:"profile_photo"`
	Description  string    `json:"description,omitempty" db:"description"`
	Phone        *string   `json:"phone,omitempty" db:"phone"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
	Blocked      bool      `json:"blocked"` // Ya no tiene etiqueta db directa, se calcula con JOIN
}

type BlockedUser struct {
	Id            int        `json:"id" db:"id"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
	BlockedUntil  *time.Time `json:"blocked_until,omitempty" db:"blocked_until"` // Puntero para permitir nulos (bloqueo permanente)
	Reason        string     `json:"reason" db:"reason"`
	BlockerId     *int       `json:"blocker_id,omitempty" db:"blocker_id"` // Puntero para permitir nulos (bloqueo del sistema)
	BlockedUserId int        `json:"blocked_user_id" db:"blocked_user_id"`
}

type UserVerification struct {
	Id              int       `json:"id" db:"id"`
	Email           string    `json:"email" db:"email"`
	Name            string    `json:"name" db:"name"`
	Surname         string    `json:"surname" db:"surname"`
	Password        string    `json:"password" db:"password"`
	VerificationPin string    `json:"pin" db:"verification_pin"`
	PinExpiration   time.Time `json:"pin_expiration" db:"pin_expiration"`
}
