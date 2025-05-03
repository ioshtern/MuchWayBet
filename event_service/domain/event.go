package domain

import "time"

type Event struct {
	ID        string
	Name      string
	StartTime time.Time
	Status    string
	WinnerID  *string
	CreatedAt time.Time
	UpdatedAt time.Time
}
