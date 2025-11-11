
package repository

import (
	"context"
	"fmt"
	"log"
	"time"
	"github.com/jackc/pgx/v5/pgxpool"

)


func NewDB(ctx context.Context) (*pgxpool.Pool, error) {
	// TODO: make this non shitty,  probs some env or somet
	dsn := "postgres://devuser:devpass@localhost:5432/mydb?sslmode=disable"
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Fatal(err)
	}


	// FIGURE THIS OUT LATER...
	cfg.MaxConns = 10
	cfg.MinConns = 1
	cfg.MaxConnLifetime = time.Hour

	pool, err := pgxpool.NewWithConfig(ctx, cfg)

	if err != nil {
        return nil, fmt.Errorf("Coudlnt create pool: %v", err)
    }
    log.Println("Connected to db")

	return pool, nil
}