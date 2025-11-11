package service

import ( 
	"context"

	"golang.org/x/crypto/bcrypt"

	"github.com/nthnklssn/sports_island/internal/models"
    "github.com/nthnklssn/sports_island/internal/repository"
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

func(s *UserService) GetByID(ctx context.Context, id int) (*models.User, error) {
	return s.Repo.GetByID(ctx, id)
}

func(s *UserService) CreateUser(ctx context.Context, name string, email string, password string) (*models.User, error) {
	// TODO: Validate password length etc ..
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	hashedStr := string(hashedPassword)
	const startingCurrency float64 = 100
	if err != nil {
		return nil, err
	}
	u := &models.User{
		Name: name,
		Email: email,
		Password: hashedStr,
		Currency: startingCurrency,
	}
	err = s.Repo.Create(ctx, u)
	if err != nil {
		return nil, err
	}
	return u, nil

}

func(s *UserService) DeleteUser(ctx context.Context, id int) error {
	return s.Repo.Delete(ctx, id)
}