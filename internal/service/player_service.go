package service

import ( 
	"context"
	"github.com/nthnklssn/nba_island/internal/models"
    "github.com/nthnklssn/nba_island/internal/repository"
)

type PlayerService struct {
	Repo repository.PlayerRepository
}

func NewPlayerService(repo repository.PlayerRepository) *PlayerService {
	return &PlayerService{Repo: repo}
}

func(s *PlayerService) GetAll(ctx context.Context) ([]*models.Player, error) {
	return s.Repo.GetAll(ctx)
}

func(s *PlayerService) GetPlayerByID(ctx context.Context, id int) (*models.Player, error) {
	return s.Repo.GetByID(ctx, id)
}

func(s *PlayerService) CreatePlayer(ctx context.Context, name string, value float64, capacity int) (error){
	p := &models.Player{
		Name: name,
		Value: value,
		Capacity: capacity,
	}
	return s.Repo.Create(ctx, p)
}

func(s *PlayerService) DeletePlayer(ctx context.Context, id int) (error){
	return s.Repo.Delete(ctx, id)
}