package models

type Player struct {
	ID int `json:"id"`
    Name  string `json:"name" binding:"required"`
    Value float64    `json:"value" binding:"required"`
    Capacity int    `json:"value" binding:"required"`
}