package service

import ( 
	"context"
	"github.com/nbaisland/nbaisland/internal/models"
    "github.com/nbaisland/nbaisland/internal/repository"
)

type PriceHistoryService struct {
	Repo repository.PlayerPriceRepository
}

func NewPriceHistoryService(repo repository.PlayerPriceRepository) *PriceHistoryService {
	return &PriceHistoryService{Repo: repo}
}

func(s *PriceHistoryService) GetPlayerPriceHistory(ctx context.Context, playerID int64, timeRange string) ([]models.PricePoint, error) {
	return s.Repo.GetPlayerPriceHistory(ctx, playerID, timeRange)
}
