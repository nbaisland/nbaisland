package models

import "time"

type PricePoint struct {
    Price     float64   `json:"price"`
    Timestamp time.Time `json:"timestamp"`
}