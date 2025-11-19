package repository

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5"
	"github.com/nbaisland/nbaisland/internal/models"
)

type HoldingRepository interface {
    GetByID(ctx context.Context, id int) (*models.Holding, error)
    GetAll(ctx context.Context) ([]*models.Holding, error)
    Create(ctx context.Context, u *models.Holding) error
    Sell(ctx context.Context, id int, sell_price float64) error
    Delete(ctx context.Context, id int) error
}

type PSQLHoldingRepo struct {
    Pool *pgxpool.Pool
}

func (r *PSQLHoldingRepo) GetByID(ctx context.Context, id int) (*models.Holding, error) {
	var h = &models.Holding{}
	err := r.Pool.QueryRow(ctx, "SELECT id, user_id, player_id, bought_for, buy_date, quantity, sold_for, sell_date, active from Holdings where id=$1", id).Scan(
		&h.ID,
		&h.UserID,
		&h.PlayerID,
		&h.BoughtFor,
		&h.BuyDate,
		&h.Quantity,
		&h.SoldFor,
		&h.SellDate,
		&h.Active,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return h, nil
}


func (r *PSQLHoldingRepo) GetAll(ctx context.Context) ([]*models.Holding, error){
	var Holdings []*models.Holding

	rows, err := r.Pool.Query(ctx, "SELECT id, user_id, player_id, bought_for, buy_date, quantity, sold_for, sell_date, active from Holdings")
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		h := &models.Holding{}
		err := rows.Scan(&h.ID,
			&h.UserID,
			&h.PlayerID,
			&h.BoughtFor,
			&h.BuyDate,
			&h.Quantity,
			&h.SoldFor,
			&h.SellDate,
			&h.Active,
		)

		if err != nil {
			return nil, err
		}
		Holdings = append(Holdings, h)
	}
	
	return Holdings, nil
}

func (r *PSQLHoldingRepo) Create(ctx context.Context, h *models.Holding) error {
	_, err := r.Pool.Exec(ctx, "INSERT INTO Holdings (user_id, player_id, bought_for, buy_date, quantity) VALUES ($1, $2, $3, NOW(), $4)", h.UserID, h.PlayerID, h.BoughtFor, h.Quantity)
	return err
}

func (r *PSQLHoldingRepo) Sell(ctx context.Context, id int, sell_price float64) error {
	// When you sell you are setting 'active' to false, no other update method should be neccesary..
	_, err := r.Pool.Exec(ctx, "UPDATE Holdings SET active=False, sold_for=$2, sell_date=NOW() where id=$1", id, sell_price)
	return err
}

func (r *PSQLHoldingRepo) Delete(ctx context.Context, id int) error {
	_, err := r.Pool.Exec(ctx, "DELETE FROM Holdings WHERE id=$1", id)
	return err
}
