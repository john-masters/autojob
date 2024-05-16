package models

type Experience struct {
	ID               int
	UserID           int
	Name             string
	Role             string
	StartDate        string
	FinishDate       string
	Current          bool
	Responsibilities string
}
