package models

type User struct {
	ID         int
	FirstName  string
	LastName   string
	Email      string
	SearchTerm string
	Password   string
	IsMember   bool
}
