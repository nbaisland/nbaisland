package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	// "github.com/nbaisland/nbaisland/internal/models"
)

type PlayerIDMapRepository interface {
	GetNBAPlayerByAppID(ctx context.Context, playerID int64) (int64, error)
	GetAppPlayerByNBAID(ctx context.Context, nbaID int64) (int64, error)
}

type PlayerMapRepo struct {
    Pool *pgxpool.Pool
}

func (r *PlayerMapRepo) GetNBAPlayerByAppID(ctx context.Context, playerID int64) (int64, error) {
	var nbaPlayerID int64

	err := r.Pool.QueryRow(ctx,
		"SELECT nba_player_id FROM player_nba_mapping WHERE player_id=$1",
		playerID,
	).Scan(&nbaPlayerID)

	if errors.Is(err, pgx.ErrNoRows) {
		return 0, nil
	}

	if err != nil {
		return 0, fmt.Errorf(
			"GetNBAPlayerByAppID failed (player_id=%d): %w",
			playerID,
			err,
		)
	}

	return nbaPlayerID, nil
}


func (r *PlayerMapRepo) GetAppPlayerByNBAID(ctx context.Context, nbaID int64) (int64, error) {
	var appPlayerID int64

	err := r.Pool.QueryRow(ctx,
		"SELECT player_id FROM player_nba_mapping WHERE nba_player_id=$1",
		nbaID,
	).Scan(&appPlayerID)

	if errors.Is(err, pgx.ErrNoRows) {
		return 0, nil
	}

	if err != nil {
		return 0, fmt.Errorf(
			"GetAppPlayerByNBAID failed (player_id=%d): %w",
			appPlayerID,
			err,
		)
	}

	return appPlayerID, nil
}
