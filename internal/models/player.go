package models

type Player struct {
	ID int64 `json:"id"`
    Name  string `json:"name" binding:"required"`
    Value float64    `json:"value" binding:"required"`
    Capacity int    `json:"capacity" binding:"required"`
}