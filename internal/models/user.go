package models

import "time"

type User struct {
	Id           int       `json:"id" db:"id"`
	Name         string    `json:"name" db:"name"`
	Surname      string    `json:"surname" db:"surname"`
	Password     string    `json:"-" db:"password"`
	Email        string    `json:"email" db:"email"`
	Location     string    `json:"location,omitempty" db:"location"`
	Role         string    `json:"role" db:"role"`
	Verified     bool      `json:"verified" db:"verified"`
	ProfilePhoto *string   `json:"profile_photo,omitempty" db:"profile_photo"`
	Description  string    `json:"description,omitempty" db:"description"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
	Blocked      bool      `json:"blocked"` // No direct db tag, calculated with JOIN
}

func (u *User) Update(user *User) {
	if user.Name != "" {
		u.Name = user.Name
	}
	if user.Surname != "" {
		u.Surname = user.Surname
	}
	if user.Location != "" {
		u.Location = user.Location
	}
	if user.Description != "" {
		u.Description = user.Description
	}
}

type BlockedUser struct {
	Id            int        `json:"id" db:"id"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
	BlockedUntil  *time.Time `json:"blocked_until,omitempty" db:"blocked_until"` // Pointer to allow nulls (permanent block)
	Reason        string     `json:"reason" db:"reason"`
	BlockerId     *int       `json:"blocker_id,omitempty" db:"blocker_id"` // Pointer to allow nulls (system block)
	BlockedUserId int        `json:"blocked_user_id" db:"blocked_user_id"`
}

type UserVerification struct {
	Id              int       `json:"id" db:"id"`
	UserId          int       `json:"user_id" db:"user_id"`
	UserEmail       string    `json:"user_email" db:"user_email"`
	VerificationPin string    `json:"pin" db:"verification_pin"`
	PinExpiration   time.Time `json:"pin_expiration" db:"pin_expiration"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
}

type EmailVerifiaction struct {
	VerificationPin string `json:"pin"`
}

type CreateUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Name     string `json:"name" binding:"required,min=3,max=60"`
	Surname  string `json:"surname" binding:"required,min=3,max=60"`
	Role     string `json:"role" binding:"required,oneof=student teacher admin"`
	Verified bool   `json:"-"` //
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type PasswordModifyRequest struct {
	Password string `json:"password" binding:"required,min=8,max=60"`
}

type LocationModifyRequest struct {
	Location string `json:"location" binding:"required"`
}
