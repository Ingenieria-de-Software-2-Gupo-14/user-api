package models

import "time"

type LoginAttempt struct {
	ID         int       `json:"id"`
	UserID     int       `json:"user_id"`
	IPAddress  string    `json:"ip_address"`
	UserAgent  string    `json:"user_agent"`
	Successful bool      `json:"successful"`
	CreatedAt  time.Time `json:"created_at"`
}
