package models

import "time"

// SavedRequest represents a single saved HTTP request within a collection.
type SavedRequest struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	UserID       uint      `gorm:"index;not null" json:"user_id"`
	CollectionID uint      `gorm:"index;not null" json:"collection_id"`
	Name         string    `gorm:"not null" json:"name"`
	Method       string    `gorm:"not null" json:"method"`
	URL          string    `gorm:"not null" json:"url"`
	Headers      string    `json:"headers"` // stored as a JSON string, e.g. {"Content-Type":"application/json"}
	Body         string    `json:"body"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}