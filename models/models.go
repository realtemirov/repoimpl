package models

type User struct {
	ID        string
	FirstName string
	LastName  string
	UserName  string
	Address   string
	Email     string
	Phone     string
	Password  string
}

type Post struct {
	ID      string
	Title   string
	Content string
	UserID  string
}
