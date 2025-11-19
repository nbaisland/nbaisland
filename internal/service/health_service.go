package service

import ( 
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

type HealthService struct {
	Pool *pgxpool.Pool
}

func NewHealthService(pool *pgxpool.Pool) *HealthService {
	return &HealthService{Pool: pool}
}

func(s *HealthService) Check(ctx context.Context) (error) {
	return s.Pool.Ping(ctx)
}

