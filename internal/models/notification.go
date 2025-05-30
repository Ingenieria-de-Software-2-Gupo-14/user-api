package models

import "time"

type NotifyRequest struct {
	Users             []int  `json:"users"  validate:"required,dive,gt=0"`
	NotificationTitle string `json:"notification_title" db:"notification_title" validate:"required,min=1,max=225"`
	NotificationText  string `json:"notification_text" db:"notification_text" validate:"required,min=1,max=225"`
}

type NotificationToken struct {
	NotificationToken string    `json:"notification_token" db:"token" validate:"required,min=1,max=225"`
	CreatedTime       time.Time `json:"created_time" db:"created_time"`
}

type NotificationTokens struct {
	NotificationTokens []NotificationToken `json:"notifications"`
}

type NotificationSetUpRequest struct {
	Token string `json:"token"  validate:"required, max=225"`
}
