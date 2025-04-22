package domain

import "time"

type Event struct {
	ID        string    `bson:"id"`
	Name      string    `bson:"name"`
	StartTime time.Time `bson:"start_time"`
	Status    string    `bson:"status"`
	WinnerID  *string   `bson:"winner_id,omitempty"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}
