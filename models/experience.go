package models

import "time"

type Experience struct {
	ID      int
	UserID  int
	Name    string
	Role    string
	Start   time.Time
	Finish  time.Time
	Current bool
	Duties  string
}
