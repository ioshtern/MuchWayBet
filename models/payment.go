package models

import "time"

type Payment struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Type      string    `json:"type"` // "deposit", "withdrawal"
	Amount    float64   `json:"amount"`
	Status    string    `json:"status"` // "pending", "completed", "failed"
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
