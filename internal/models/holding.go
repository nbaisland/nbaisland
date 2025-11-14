package models

import "time"

type Holding struct {
    ID   int `json:"id"`
    UserID   int `json:"user_id" binding:"required"`
    PlayerID int `json:"player_id" binding:"required"`
    Quantity float64 `json:"quantity" binding:"required"`
	BuyDate time.Time `json:"buy_date"`
	BoughtFor float64 `json:"buy_price"`
	SoldFor *float64 `json:"sell_price"` // Pointer needed to avoid nul issues
	SellDate *time.Time `json:"sell_date"` // Pointer needed to avoid nul issues
	Active bool `json:"active"`

}