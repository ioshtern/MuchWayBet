package domain

import "time"

type Bet struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	EventID   string    `json:"event_id"`
	Amount    float64   `json:"amount"`
	Odds      float64   `json:"odds"`
	Status    string    `json:"status"`
	Payout    float64   `json:"payout,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
