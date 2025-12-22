package repository

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5"
	"github.com/nbaisland/nbaisland/internal/models"
)

type TransactionRepository interface {
    GetByID(ctx context.Context, id int64) (*models.Transaction, error)
    GetByUserID(ctx context.Context, id int64) ([]*models.Transaction, error)
    GetByPlayerID(ctx context.Context, id int64) ([]*models.Transaction, error)
    GetAll(ctx context.Context) ([]*models.Transaction, error)
    CreateTransaction(ctx context.Context, u *models.Transaction) error
    Delete(ctx context.Context, id int64) error
	GetPositionsByUserIDAndPlayerID(ctx context.Context, user_id int64, player_id int64) (*models.Position, error)
	GetAllPositions(ctx context.Context) ([]*models.Position, error)
	GetPositionsByUserID(ctx context.Context, id int64) ([]*models.Position, error)
	GetPositionsByPlayerID(ctx context.Context, id int64) ([]*models.Position, error)
	RefreshPositionsMV(ctx context.Context) error 
}

type PSQLTransactionRepo struct {
    Pool *pgxpool.Pool
}

func scanTransaction(row pgx.Row) (*models.Transaction, error) {
    var t models.Transaction
    err := row.Scan(
        &t.ID,
        &t.UserID,
        &t.AssetID,
        &t.Type,
        &t.Quantity,
        &t.Price,
        &t.Timestamp,
    )
    if err != nil {
        return nil, err
    }
    return &t, nil
}

func scanTransactionRows(rows pgx.Rows) ([]*models.Transaction, error) {
    var transactions []*models.Transaction
    for rows.Next() {
        var t models.Transaction
        err := rows.Scan(
			&t.ID,
			&t.UserID,
			&t.AssetID,
			&t.Type,
			&t.Quantity,
			&t.Price,
			&t.Timestamp,
		)
		if err != nil {
            return nil, err
        }
        transactions = append(transactions, &t)
    }
    if rows.Err() != nil {
        return nil, rows.Err()
    }
    return transactions, nil
}


func (r *PSQLTransactionRepo) GetByID(ctx context.Context, id int64) (*models.Transaction, error) {
	row := r.Pool.QueryRow(ctx, "SELECT id, user_id, asset_id, type, quantity, price, timestamp from transactions where id=$1", id)

	t, err := scanTransaction(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return t, nil
}

func (r *PSQLTransactionRepo)  GetByUserID(ctx context.Context, id int64) ([]*models.Transaction, error){
	rows, err := r.Pool.Query(ctx, "SELECT id, user_id, asset_id, type, quantity, price, timestamp from transactions WHERE user_id=$1", id)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	Transactions, err := scanTransactionRows(rows)

	return Transactions, nil
}


func (r *PSQLTransactionRepo) GetByPlayerID(ctx context.Context, id int64) ([]*models.Transaction, error){
	rows, err := r.Pool.Query(ctx, "SELECT id, user_id, asset_id, type, quantity, price, timestamp from transactions WHERE asset_id=$1", id)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	Transactions, err := scanTransactionRows(rows)

	return Transactions, nil
}


func (r *PSQLTransactionRepo) GetAll(ctx context.Context) ([]*models.Transaction, error){
	rows, err := r.Pool.Query(ctx, "SELECT id, user_id, asset_id, type, quantity, price, timestamp from transactions")
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	Transactions, err := scanTransactionRows(rows)

	return Transactions, nil
}



func (r *PSQLTransactionRepo) CreateTransaction(ctx context.Context, t *models.Transaction) error {
	_, err := r.Pool.Exec(ctx, "INSERT INTO transactions (user_id, asset_id, type, quantity, price, timestamp) VALUES ($1, $2, $3, $4, $5, $6)", t.UserID, t.AssetID, t.Type, t.Quantity, t.Price, t.Timestamp)
	return err
}

func (r *PSQLTransactionRepo) Delete(ctx context.Context, id int64) error {
	_, err := r.Pool.Exec(ctx, "DELETE FROM transactions WHERE id=$1", id)
	return err
}

func (r *PSQLTransactionRepo) GetPositionsByUserIDAndPlayerID(ctx context.Context, user_id int64, player_id int64) (*models.Position, error) {
	var p = &models.Position{}
	err := r.Pool.QueryRow(ctx, "SELECT user_id, asset_id, quantity, average_cost from positions_mv where user_id=$1 AND asset_id=$2", user_id, player_id).Scan(
		&p.UserID,
		&p.AssetID,
		&p.Quantity,
		&p.AverageCost,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return p, nil
}

func (r *PSQLTransactionRepo) GetAllPositions(ctx context.Context) ([]*models.Position, error){
	var positions []*models.Position
	rows, err := r.Pool.Query(ctx, "SELECT user_id, asset_id, quantity, average_cost from positions_mv")
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		p := &models.Position{}
		err := rows.Scan(
			&p.UserID,
			&p.AssetID,
			&p.Quantity,
			&p.AverageCost,
		)
		if err != nil {
			return nil, err
		}
		positions = append(positions, p)
	}
	return positions, nil
}

func (r *PSQLTransactionRepo) GetPositionsByUserID(ctx context.Context, id int64) ([]*models.Position, error){
	var positions []*models.Position
	rows, err := r.Pool.Query(ctx, "SELECT user_id, asset_id, quantity, average_cost from positions_mv WHERE user_id=$1", id)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		p := &models.Position{}
		err := rows.Scan(
			&p.UserID,
			&p.AssetID,
			&p.Quantity,
			&p.AverageCost,
		)
		if err != nil {
			return nil, err
		}
		positions = append(positions, p)
	}
	return positions, nil
}

func (r *PSQLTransactionRepo) GetPositionsByPlayerID(ctx context.Context, id int64) ([]*models.Position, error){
	var positions []*models.Position
	rows, err := r.Pool.Query(ctx, "SELECT user_id, asset_id, quantity, average_cost from positions_mv WHERE asset_id=$1", id)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		p := &models.Position{}
		err := rows.Scan(
			&p.UserID,
			&p.AssetID,
			&p.Quantity,
			&p.AverageCost,
		)
		if err != nil {
			return nil, err
		}
		positions = append(positions, p)
	}
	return positions, nil
}

func (r *PSQLTransactionRepo) RefreshPositionsMV(ctx context.Context) error {
    _, err := r.Pool.Exec(ctx, "REFRESH MATERIALIZED VIEW positions_mv")
    if err != nil {
        return err
    }
    return nil
}
