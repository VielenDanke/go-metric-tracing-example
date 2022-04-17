package model

type User struct {
	ID       int
	Username string
	Email    string
}

func NewUser(username, email string) User {
	return User{Username: username, Email: username}
}
