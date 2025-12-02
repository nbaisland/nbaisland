package models

import "time"

type Transaction struct {
    ID   int64 `json:"id"`
    UserID   int64 `json:"user_id" binding:"required"`
    AssetID int64 `json:"player_id" binding:"required"`
	Type string `json:"type" binding:"required`
    Quantity float64 `json:"quantity" binding:"required"`
	Price float64  `json:"price" binding:"required"`
	Timestamp time.Time `json:"timestamp" binding:"required"`

}