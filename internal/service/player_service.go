package service

import ( 
	"context"
	"github.com/nbaisland/nbaisland/internal/models"
    "github.com/nbaisland/nbaisland/internal/repository"
    "github.com/nbaisland/nbaisland/internal/utils"
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

func(s *PlayerService) GetPlayersByIDs(ctx context.Context, player_ids []int64) ([]*models.Player, error) {
	return s.Repo.GetByIDs(ctx, player_ids)
}

func(s *PlayerService) GetPlayerByID(ctx context.Context, id int64) (*models.Player, error) {
	return s.Repo.GetByID(ctx, id)
}

func(s *PlayerService) GetPlayerBySlug(ctx context.Context, slug string) (*models.Player, error) {
	return s.Repo.GetBySlug(ctx, slug)
}

func(s *PlayerService) CreatePlayer(ctx context.Context, name string, value float64, capacity int, slug string) (error){
	if slug == "" {
        slug = utils.ToSlug(name)
    }
	p := &models.Player{
		Name: name,
		Value: value,
		Capacity: capacity,
		Slug: slug,
	}
	return s.Repo.Create(ctx, p)
}

func(s *PlayerService) DeletePlayer(ctx context.Context, id int64) (error){
	return s.Repo.Delete(ctx, id)
}