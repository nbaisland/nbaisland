package models

type Position struct {
    UserID   int64 `json:"user_id" binding:"required"`
    AssetID int64 `json:"player_id" binding:"required"`
    Quantity float64 `json:"quantity" binding:"required"`
}