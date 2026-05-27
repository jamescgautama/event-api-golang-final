package models

import "time"

type Event struct {
	ID             string `json:"id"`
	Title          string `json:"title"`
	TotalQuota     int    `json:"total_quota"`
	AvailableQuota int    `json:"available_quota"`
}

type Booking struct {
	ID        string    `json:"id"`
	EventID   string    `json:"event_id"`
	UserID    string    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

type User struct {
	ID string `json:"id"`
}
