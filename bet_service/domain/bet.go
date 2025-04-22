package domain

import "time"

type Bet struct {
	ID        string    `bson:"id"`
	UserID    string    `bson:"user_id"`
	EventID   string    `bson:"event_id"`
	Amount    float64   `bson:"amount"`
	Odds      float64   `bson:"odds"`
	Status    string    `bson:"status"`           // "pending", "won", "lost"
	Payout    float64   `bson:"payout,omitempty"` // If won
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}
