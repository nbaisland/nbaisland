package repository

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nbaisland/nbaisland/internal/models"
)

type PlayerRepository interface {
    GetByID(ctx context.Context, id int64) (*models.Player, error)
	GetValueByID(ctx context.Context, id int64) (float64, error)
    GetAll(ctx context.Context) ([]*models.Player, error)
	GetByIDs(ctx context.Context, ids []int64) ([]*models.Player, error)
    Create(ctx context.Context, u *models.Player) error
    Update(ctx context.Context, u *models.Player) error
    Delete(ctx context.Context, id int64) error
}

type PSQLPlayerRepo struct {
    Pool *pgxpool.Pool
}

func (r *PSQLPlayerRepo) GetByID(ctx context.Context, id int64) (*models.Player, error) {
	var p = &models.Player{}
	err := r.Pool.QueryRow(ctx, "SELECT id, name, value, capacity from players where id=$1", id).Scan(
		&p.ID,
		&p.Name,
		&p.Value,
		&p.Capacity,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return p, nil
}

func (r *PSQLPlayerRepo) GetByIDs(ctx context.Context, ids []int64) ([]*models.Player, error) {
	var players []*models.Player
	rows, err := r.Pool.Query(ctx, "SELECT id, name, value, capacity from players where id=ANY($1)", ids)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		p := &models.Player{}
		err := rows.Scan(&p.ID,
			&p.Name,
			&p.Value,
			&p.Capacity,
		)

		if err != nil {
			return nil, err
		}
		players = append(players, p)
	}
	return players, nil
}

func (r *PSQLPlayerRepo) GetValueByID(ctx context.Context, id int64) (float64, error) {
	var value float64
	err := r.Pool.QueryRow(ctx, "SELECT value from players where id=$1", id).Scan(&value)

	if errors.Is(err, pgx.ErrNoRows) {
		return 0, nil
	}

	if err != nil {
		return 0, err
	}

	return value, nil
}

func (r *PSQLPlayerRepo) GetAll(ctx context.Context) ([]*models.Player, error){
	var players []*models.Player

	rows, err := r.Pool.Query(ctx, "SELECT id, name, value, capacity from players")
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		p := &models.Player{}
		err := rows.Scan(&p.ID,
			&p.Name,
			&p.Value,
			&p.Capacity,
		)

		if err != nil {
			return nil, err
		}
		players = append(players, p)
	}
	return players, nil
}

func (r *PSQLPlayerRepo) Create(ctx context.Context, p *models.Player) error {
	_, err := r.Pool.Exec(ctx, "INSERT INTO players (name, value, capacity) VALUES ($1, $2, $3)", p.Name, p.Value, p.Capacity)
	return err
}

func (r *PSQLPlayerRepo) Update(ctx context.Context, p *models.Player) error {
	_, err := r.Pool.Exec(ctx, "UPDATE players SET name=$2, value=$3, capacity=$4 where id = $1", p.ID, p.Name, p.Value, p.Capacity)
	return err
}

func (r *PSQLPlayerRepo) UpdateValue(ctx context.Context, id int64, v float64) error {
	_, err := r.Pool.Exec(ctx, "UPDATE players SET value=$1 WHERE id=$2", v, id)
	return err
}

func (r *PSQLPlayerRepo) UpdateCapacity(ctx context.Context, id int64, c float64) error {
	_, err := r.Pool.Exec(ctx, "UPDATE players SET capacity=$1 WHERE id=$2", c, id)
	return err
}

func (r *PSQLPlayerRepo) Delete(ctx context.Context, id int64) error {
	_, err := r.Pool.Exec(ctx, "DELETE FROM players WHERE id=$1", id)
	return err
}
