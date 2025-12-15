
package repository

import (
	"context"
	"fmt"
	"log"
	"time"
	"github.com/jackc/pgx/v5/pgxpool"

)


func NewDB(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Fatal(err)
	}


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