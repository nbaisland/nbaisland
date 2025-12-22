package service

import ( 
	"context"
	"log"
	"github.com/nbaisland/nbaisland/internal/models"
    "github.com/nbaisland/nbaisland/internal/repository"
)

type UserService struct {
	Repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{Repo: repo}
}

func(s *UserService) GetAll(ctx context.Context) ([]*models.User, error) {
	return s.Repo.GetAll(ctx)
}

func(s *UserService) GetByID(ctx context.Context, id int64) (*models.User, error) {
	return s.Repo.GetByID(ctx, id)
}

func(s *UserService) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	return s.Repo.GetByUsername(ctx, username)
}

func(s *UserService) CreateUser(ctx context.Context, u *models.User) (error) {
	err := s.Repo.Create(ctx, u)
	if err != nil {
		log.Printf("Error: %v", err)
	}
	return err
}

func(s *UserService) DeleteUser(ctx context.Context, id int64) error {
	return s.Repo.Delete(ctx, id)
}