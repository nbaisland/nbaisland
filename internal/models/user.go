package models

type User struct {
	ID int64 `json:"id"`
    Username string `json:"username" binding:"required"`
    Name  string `json:"name" binding:"required"`
    Email string `json:"email" binding:"required,email"`
    Password string `json:"password"`
    Currency float64	`json:"currency"`
}