package models

import "time"

// RequestHistory logs a single executed HTTP request (sent via the "send" endpoint).
type RequestHistory struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	UserID         uint      `gorm:"index;not null" json:"user_id"`
	Method         string    `gorm:"not null" json:"method"`
	URL            string    `gorm:"not null" json:"url"`
	StatusCode     int       `json:"status_code"`
	ResponseTimeMs int64     `json:"response_time_ms"`
	Error          string    `json:"error,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
}