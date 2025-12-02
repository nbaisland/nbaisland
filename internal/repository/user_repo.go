package repository

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nbaisland/nbaisland/internal/models"
)

type UserRepository interface {
    GetByID(ctx context.Context, id int64) (*models.User, error)
    GetAll(ctx context.Context) ([]*models.User, error)
    Create(ctx context.Context, u *models.User) error
    UpdateName(ctx context.Context, id int64, name string) error
    UpdatePassword(ctx context.Context, id int64, password string) error
    UpdateEmail(ctx context.Context, id int64, email string) error
    UpdateCurrency(ctx context.Context, id int64, currency float64) error
    Delete(ctx context.Context, id int64) error
}

type PSQLUserRepo struct {
    Pool *pgxpool.Pool
}

func (r *PSQLUserRepo) GetByID(ctx context.Context, id int64) (*models.User, error) {
	var u = &models.User{}
	err := r.Pool.QueryRow(ctx, "SELECT id, name, email, currency from users where id=$1", id).Scan(
		&u.ID,
		&u.Name,
		&u.Email,
		&u.Currency,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return u, nil
}


func (r *PSQLUserRepo) GetAll(ctx context.Context) ([]*models.User, error){
	var users []*models.User

	rows, err := r.Pool.Query(ctx, "SELECT id, name, email, currency from users")
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	for rows.Next() {
		u := &models.User{}
		err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Currency)

		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func (r *PSQLUserRepo) Create(ctx context.Context, u *models.User) error {
	err := r.Pool.QueryRow(ctx, "INSERT INTO users (name, email, currency, password) VALUES ($1, $2, $3, $4) RETURNING id",
	u.Name, u.Email, u.Currency, u.Password,
	).Scan(&u.ID)
	return err
}

func (r *PSQLUserRepo) UpdateName(ctx context.Context, id int64, name string) error {
	_, err := r.Pool.Exec(ctx, "UPDATE users SET name=$2 where id = $1", id, name)
	return err
}

func (r *PSQLUserRepo) UpdateEmail(ctx context.Context, id int64, email string) error {
	_, err := r.Pool.Exec(ctx, "UPDATE users SET email=$2 where id = $1", id, email)
	return err
}

func (r *PSQLUserRepo) UpdatePassword(ctx context.Context, id int64, password string) error {
	_, err := r.Pool.Exec(ctx, "UPDATE users SET password=$2 where id = $1", id, password)
	return err
}

func (r *PSQLUserRepo) UpdateCurrency(ctx context.Context, id int64, currency float64) error {
	_, err := r.Pool.Exec(ctx, "UPDATE users SET currency=$2 where id = $1", id, currency)
	return err
}

func (r *PSQLUserRepo) Delete(ctx context.Context, id int64) error {
	_, err := r.Pool.Exec(ctx, "DELETE FROM users WHERE id=$1", id)
	return err
}
