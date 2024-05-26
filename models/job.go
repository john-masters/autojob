package models

type Job struct {
	ID          int
	UserID      int
	Title       string
	Company     string
	Link        string
	Description string
	CoverLetter string
}
