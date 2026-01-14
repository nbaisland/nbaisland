package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nbaisland/nbaisland/internal/models"
)

type PlayerPriceRepository interface {
	GetPlayerPriceHistory(ctx context.Context, playerID int64, timeRange string) ([]models.PricePoint, error)
	GetAllPlayersPriceHistory(ctx context.Context, timeRange string) (map[int64][]models.PricePoint, error)
	RecordPlayerPrice(ctx context.Context, playerID int64, price float64) error
}

type PSQLPlayerPriceRepo struct {
	Pool *pgxpool.Pool
}

func timeRangeInterval(timeRange string) string {
	switch timeRange {
		case "7d":
			return "7 days"
		case "30d":
			return "30 days"
		case "90d":
			return "90 days"
		case "1y":
			return "1 year"
		case "all":
			return "100 years"
		default:
			return "30 days"
	}
}

func (r *PSQLPlayerPriceRepo) GetPlayerPriceHistory(ctx context.Context, playerID int64, timeRange string) ([]models.PricePoint, error) {
	interval := timeRangeInterval(timeRange)

	query := fmt.Sprintf(`
		SELECT price, timestamp FROM player_price_history WHERE player_id = $1
		AND timestamp >= NOW() - INTERVAL '%s'
		ORDER BY timestamp ASC`, interval)

	rows, err := r.Pool.Query(ctx, query, playerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []models.PricePoint

	for rows.Next() {
		var p models.PricePoint
		if err := rows.Scan(&p.Price, &p.Timestamp); err != nil {
			return nil, err
		}
		history = append(history, p)
	}

	return history, rows.Err()
}

func (r *PSQLPlayerPriceRepo) GetAllPlayersPriceHistory(ctx context.Context, timeRange string) (map[int64][]models.PricePoint, error) {

	interval := timeRangeInterval(timeRange)

	query := fmt.Sprintf(`
		SELECT player_id, price, timestamp FROM player_price_history WHERE timestamp >= NOW() - INTERVAL '%s'
		ORDER BY player_id, timestamp ASC`, interval)

	rows, err := r.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[int64][]models.PricePoint)

	for rows.Next() {
		var playerID int64
		var p models.PricePoint

		if err := rows.Scan(&playerID, &p.Price, &p.Timestamp); err != nil {
			return nil, err
		}

		result[playerID] = append(result[playerID], p)
	}

	return result, rows.Err()
}

func (r *PSQLPlayerPriceRepo) RecordPlayerPrice(ctx context.Context, playerID int64, price float64) error {
	_, err := r.Pool.Exec(ctx, `
		INSERT INTO player_price_history (player_id, price, timestamp)
		VALUES ($1, $2, NOW())`, playerID, price)
	return err
}