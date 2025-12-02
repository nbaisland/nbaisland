// internal/models/user.go
package models

type User struct {
	ID int64 `json:"id"`
    Name  string `json:"name" binding:"required"`
    Email string `json:"email" binding:"required,email"`
    Password string `json:"password"`
    Currency float64	`json:"currency"`
}