package internal

import "time"

type Posts struct {
	ID        string
	Title     string
	Content   string
	CreatedAt time.Time
}
