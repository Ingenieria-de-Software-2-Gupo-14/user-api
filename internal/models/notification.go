package models

import "time"

type NotifyRequest struct {
	Users             []int  `json:"users"  validate:"required,dive,gt=0"`
	NotificationTitle string `json:"notification_title" db:"notification_title" validate:"required,min=1,max=225"`
	NotificationText  string `json:"notification_text" db:"notification_text" validate:"required,min=1,max=225"`
	NotificationType  string `json:"notification_type"  validate:"required, oneof=exam_notification homework_notification social_notification"`
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

type NotificationPreferenceRequest struct {
	NotificationType       string `json:"notification_type"  validate:"required, oneof=exam_notification homework_notification social_notification"`
	NotificationPreference bool   `json:"notification_preference"  validate:"required"`
}

type NotificationPreference struct {
	ExamNotification     bool `json:"exam_notification"  validate:"required"`
	HomeworkNotification bool `json:"homework_notification"  validate:"required"`
	SocialNotification   bool `json:"social_notification"  validate:"required"`
}
