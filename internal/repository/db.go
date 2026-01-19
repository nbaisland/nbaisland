
package repository

import (
	"context"
	"fmt"
	"time"
	"go.uber.org/zap"

	"github.com/nbaisland/nbaisland/internal/logger"
	"github.com/jackc/pgx/v5/pgxpool"

)


func NewDB(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		logger.Log.Fatal("Failed to Connect to db",
			zap.Error(err),
			zap.String("handler", "NewDB"),
		)
	}


	cfg.MaxConns = 10
	cfg.MinConns = 1
	cfg.MaxConnLifetime = time.Hour

	pool, err := pgxpool.NewWithConfig(ctx, cfg)

	if err != nil {
		logger.Log.Error("Failed to Create Pool",
			zap.Error(err),
			zap.String("handler", "NewDB"),
		)
        return nil, fmt.Errorf("Coudlnt create pool: %v", err)
    }

	return pool, nil
}