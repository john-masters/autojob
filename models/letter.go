package models

import "time"

type Letter struct {
	ID        int
	UserID    int
	Letter    string
	CreatedAt time.Time
}
