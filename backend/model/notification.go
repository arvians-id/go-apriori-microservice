package model

import "time"

type Notification struct {
	IdNotification int64     `json:"id_notification"`
	UserId         int64     `json:"user_id"`
	Title          string    `json:"title"`
	Description    *string   `json:"description"`
	URL            *string   `json:"url"`
	IsRead         bool      `json:"is_read"`
	CreatedAt      time.Time `json:"created_at"`
}
