package models

import "time"

type Letter struct {
	ID        int
	UserID    int
	Content   string
	CreatedAt time.Time
}
