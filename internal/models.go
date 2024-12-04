package internal

import (
	"time"
)

type Posts struct {
	ID        string
	Title     string
	Content   string
	CreatedAt time.Time
}
type User struct {
	ID       string
	Username string
	Email    string
	Password string
	Confirm  string
}
