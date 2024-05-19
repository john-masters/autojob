package models

type History struct {
	ID      int
	UserID  int
	Name    string
	Role    string
	Start   string
	Finish  string
	Current bool
	Duties  string
}
