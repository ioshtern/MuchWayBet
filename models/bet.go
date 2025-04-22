package models

import "time"

type Bet struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	EventID   string    `json:"event_id"`
	Amount    float64   `json:"amount"`
	Odds      float64   `json:"odds"`
	Status    string    `json:"status"`           // "pending", "won", "lost"
	Payout    float64   `json:"payout,omitempty"` // If won
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
