package service

import ( 
	"context"
	"fmt"
	"time"
	"github.com/nbaisland/nbaisland/internal/models"
    "github.com/nbaisland/nbaisland/internal/repository"
)
type TransactionError struct {
	Code string
	Msg string
}

func (e *TransactionError) Error() string {
    return fmt.Sprintf("%s: %s", e.Code, e.Msg)
}

type TransactionService struct {
	TransactionRepo repository.TransactionRepository
	PlayerRepo repository.PlayerRepository
	UserRepo repository.UserRepository
}

func NewTransactionService(transactionRepo repository.TransactionRepository, playerRepo repository.PlayerRepository, userRepo repository.UserRepository) *TransactionService {
	return &TransactionService{TransactionRepo: transactionRepo, PlayerRepo: playerRepo, UserRepo: userRepo}
}

func (s *TransactionService) GetAll(ctx context.Context) ([]*models.Transaction, error){
	return s.TransactionRepo.GetAll(ctx)
}

func (s *TransactionService) GetTransactionByID(ctx context.Context, id int64) (*models.Transaction, error){
	t, err := s.TransactionRepo.GetByID(ctx, id)
	return t, err
}

func (s *TransactionService) GetByUserID(ctx context.Context, id int64) ([]*models.Transaction, error){
	t, err := s.TransactionRepo.GetByUserID(ctx, id)
	return t, err
}

func (s *TransactionService) GetByPlayerID(ctx context.Context, id int64) ([]*models.Transaction, error){
	t, err := s.TransactionRepo.GetByPlayerID(ctx, id)
	return t, err
}

func (s *TransactionService) Buy(ctx context.Context, userID int64, playerID int64, quantity int) (error) {
	userDetail, err := s.UserRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	if userDetail == nil {
		return &TransactionError{
			Code: "UserNotFound",
			Msg: "Could not find user",
		}
	}
	playerDetail, err := s.PlayerRepo.GetByID(ctx, playerID)
	if err != nil {
		return err
	}
	if playerDetail == nil {
		return &TransactionError{
			Code: "PlayerNotFound",
			Msg: "Could not find player",
		}
	}
	cost := float64(playerDetail.Value) * float64(quantity)
	if cost > userDetail.Currency {
		return &TransactionError{
			Code: "USER_LACKS_MONEY",
			Msg: fmt.Sprintf("This trade would cost %v, user only has %v", cost, userDetail.Currency),
		}
	}
	if quantity > playerDetail.Capacity {
		return &TransactionError{
			Code: "NO_CAPACITY",
			Msg: fmt.Sprintf("Player only has %v capacity remaining, exceeding %v requested", playerDetail.Capacity, quantity),
		}
	}
	// TODO: Check for capacity or something
	buyT := &models.Transaction{
		UserID:   userID,
        AssetID:  playerID,
        Type:     "BUY",
        Quantity: quantity,
        Price:    playerDetail.Value,
        Timestamp: time.Now(),
	}
	err = s.TransactionRepo.CreateTransaction(ctx, buyT)
	if err != nil {
		return err
	}
	newCurrencyValue := userDetail.Currency - cost
	err = s.UserRepo.UpdateCurrency(ctx, userID, newCurrencyValue)
	if err != nil {
		return err
	}
    err = s.TransactionRepo.RefreshPositionsMV(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (s *TransactionService) Sell(ctx context.Context, userID int64, playerID int64, quantity int) (float64, error) {
    if quantity <= 0 {
        return 0, &TransactionError{
            Code: "QUANTITY_INVALID",
            Msg:  "Must provide a number greater than 0 to sell",
        }
    }
    playerDetail, err := s.PlayerRepo.GetByID(ctx, playerID)
	if err != nil {
		return 0, err
	}
    position, err := s.TransactionRepo.GetPositionsByUserIDAndPlayerID(ctx, userID, playerID)
    if err != nil {
        return 0, err
    }
    if position == nil {
        return 0, &TransactionError{
            Code: "NO_POSITION",
            Msg:  "Could not find position",
        }
    }
    if position.Quantity < quantity {
        return 0, &TransactionError{
            Code: "QUANTITY_EXCEEDS_POSITION",
            Msg:  fmt.Sprintf("Request to sell %v exceeds held position (%v)", quantity, position.Quantity),
        }
    }
    totalValue := float64(quantity) * playerDetail.Value
	sellT := &models.Transaction{
        UserID:   userID,
        AssetID:  playerID,
        Type:     "SELL",
        Quantity: quantity,
        Price:    playerDetail.Value,
        Timestamp: time.Now(),
    }

    err = s.TransactionRepo.CreateTransaction(ctx, sellT)
    if err != nil {
        return 0, err
    }
	userDetail, err := s.UserRepo.GetByID(ctx, userID)
	newCurrencyValue := userDetail.Currency + totalValue
	s.UserRepo.UpdateCurrency(ctx, userID, newCurrencyValue)
	// TODO: Update capacity... 
	s.PlayerRepo.UpdateCapacity(ctx, playerID, quantity)
    err = s.TransactionRepo.RefreshPositionsMV(ctx)
	if err != nil {
		return totalValue, err
	}
    return totalValue, nil
}


func (s *TransactionService) GetPositions(ctx context.Context) ([]*models.Position, error){
	return s.TransactionRepo.GetAllPositions(ctx)
}

func (s *TransactionService) GetPositionsByUserID(ctx context.Context, id int64) ([]*models.Position, error){
	return s.TransactionRepo.GetPositionsByUserID(ctx, id)
}

func (s *TransactionService) GetPositionsByPlayerID(ctx context.Context, id int64) ([]*models.Position, error){
	return s.TransactionRepo.GetPositionsByPlayerID(ctx, id)
}