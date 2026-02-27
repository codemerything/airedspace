package main

type User struct {
	UserID   int    `json:"id" db:"id"`
	Username string `json:"username" db:"username"`
	Name     string `json:"name" db:"name"`
	LastName string `json:"lastname" db:"lastname"`
	Email    string `json:"email" db:"email"`
	Password string `json:"password" db:"password"`
}
