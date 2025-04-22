package domain

import "time"

type Payment struct {
	ID        string    `bson:"id"`
	UserID    string    `bson:"user_id"`
	Type      string    `bson:"type"` // "deposit", "withdrawal"
	Amount    float64   `bson:"amount"`
	Status    string    `bson:"status"` // "pending", "completed", "failed"
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}
