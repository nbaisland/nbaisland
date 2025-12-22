package models

type Position struct {
    UserID   int64 `json:"user_id" binding:"required"`
    AssetID int64 `json:"player_id" binding:"required"`
    Quantity int `json:"quantity" binding:"required"`
    AverageCost float64 `json:"average_cost" binding:"required"`
}