package service

import ( 
	"context"
	"fmt"
	"github.com/nthnklssn/sports_island/internal/models"
    "github.com/nthnklssn/sports_island/internal/repository"
)
type HoldingError struct {
	Code string
	Msg string
}

func (e *HoldingError) Error() string {
    return fmt.Sprintf("%s: %s", e.Code, e.Msg)
}

type HoldingService struct {
	HoldingRepo repository.HoldingRepository
	PlayerRepo repository.PlayerRepository
	UserRepo repository.UserRepository
}

func NewHoldingService(holdingRepo repository.HoldingRepository, playerRepo repository.PlayerRepository, userRepo repository.UserRepository) *HoldingService {
	return &HoldingService{HoldingRepo: holdingRepo, PlayerRepo: playerRepo, UserRepo: userRepo}
}

func (s *HoldingService) GetAll(ctx context.Context) ([]*models.Holding, error){
	return s.HoldingRepo.GetAll(ctx)
}

func (s *HoldingService) GetHoldingByID(ctx context.Context, id int) (*models.Holding, error){
	h, err := s.HoldingRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return h, err
}

func (s *HoldingService) MakePurchase(ctx context.Context, player_id int, user_id int, quantity float64) (error) {
	userDetail, err := s.UserRepo.GetByID(ctx, user_id)
	if err != nil {
		return err
	}
	playerDetail, err := s.PlayerRepo.GetByID(ctx, player_id)
	if err != nil {
		return err
	}
	cost := float64(playerDetail.Value) * quantity
	if cost > userDetail.Currency {
		return &HoldingError{Code: "USER_LACKS_MONEY", Msg: fmt.Sprintf("This trade would cost %v, user only has %v", cost, userDetail.Currency)}
	}
	// TODO: Check for capacity or something
	h := &models.Holding{
		UserID : user_id,
		PlayerID : player_id,
		Quantity : quantity,
		BoughtFor : playerDetail.Value,
	}
	return s.HoldingRepo.Create(ctx, h)
}

func(s *HoldingService) SellHolding(ctx context.Context, holdingId int, quantity float64) (float64, error) {
	if quantity <= 0 {
		return 0, &HoldingError{Code : "QUANTITY_INVALID", Msg : "Must provide a number greater than 0 to sell"}
	}
	holding, err := s.HoldingRepo.GetByID(ctx, holdingId)
	if err != nil {
		return 0, err
	}
	if holding == nil {
		return 0, &HoldingError{Code : "HOLDING_NOT_FOUND", Msg : "Could not find holding"}

	}
	if holding.Quantity < quantity {
		return 0, &HoldingError{Code : "QUANTITY_EXCEEDS_HOLDING", Msg : fmt.Sprintf("Request to sell %v of holding exceeds held capacity (%v)", quantity, holding.Quantity)}
	}

	sell_price, err := s.PlayerRepo.GetValueByID(ctx, holding.PlayerID)
	if err != nil {
		return 0, err
	}
	err = s.HoldingRepo.Sell(ctx, holdingId, sell_price)
	if err != nil {
		return 0, err
	}
	proceeds := quantity * sell_price

	remaining := holding.Quantity - quantity

	if remaining > 1e-9 {
		holding.Quantity = remaining
		err = s.HoldingRepo.Create(ctx, holding)
		if err != nil {
			// If we couldnt create a new player (not sure why that is)
			// Just sell entire holding and return money to user
			// TODO: Fix to make this less jank idk low priority for now
			return proceeds + (remaining * sell_price), err
		}
		return proceeds, err
	} else {
		return proceeds, err
	}
}