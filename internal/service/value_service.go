package service

import ( 
	"context"
	// "fmt"
    "github.com/nbaisland/nbaisland/internal/repository"
    "github.com/nbaisland/nbaisland/internal/nba"
)


type ValueService struct {
	PlayerRepo repository.PlayerRepository
	NBARepo *nba.Repository
	playerMapRepo repository.PlayerIDMapRepository
}

func NewValueService(playerRepo repository.PlayerRepository, nbaRepo *nba.Repository, playerMapRepo repository.PlayerIDMapRepository) *ValueService {
	return &ValueService{PlayerRepo: playerRepo, NBARepo: nbaRepo, playerMapRepo: playerMapRepo}
}

// General thinking here is that value should be added to players who:
// Win, put up good stats, play minutes , Play games (rock solid workhorses should be more valuable than glass cannons)
// Players like lebron james should still be worth a lot due to career
// Should create a faster method for doing multiple !! TODO: Important
func(s *ValueService) CalculateValueBasedOnStats(ctx context.Context, playerID int64, season string) (float64, error) {
	nbaID, err := s.playerMapRepo.GetNBAPlayerByAppID(ctx, playerID)
	if err != nil {
		return 0, err
	}
	seasonStats, err := s.NBARepo.GetSeasonStats(ctx, nbaID, season)
	if err != nil {
		return 0, err
	}
	careerStats, err := s.NBARepo.GetCareerStats(ctx, nbaID)
	if err != nil {
		return 0, err
	}
	capacityRemaining, err := s.PlayerRepo.GetCapacityByID(ctx, playerID)
	if err != nil {
		return 0, err
	}
	// TODO: Add getting player tags, other information like champsionships, mvps, defense, all nba into value equation 
	// awardsValue := Championships * 2
	var seasonValue float64
	if seasonStats != nil && seasonStats.GamesPlayed > 10 {
		seasonValue = seasonStats.PointsPerGame + seasonStats.AssistsPerGame * 2 + seasonStats.ReboundsPerGame * 2 + seasonStats.StealsPerGame * 3 + seasonStats.BlocksPerGame * 3
	}else {
		seasonValue = 0.0
	}
	var careerValue float64
	if careerStats != nil {
		careerValue = careerStats.PointsTotal * 0.001 + careerStats.ReboundsTotal * 0.002 + careerStats.AssistsTotal * 0.002 + careerStats.StealsTotal * 0.0025 + careerStats.BlocksTotal * 0.0025 + careerStats.MinutesTotal * 0.00001

	} else {
		careerValue = 0.0
	}

	seasonMult := 1.0
	careerMult := 1.0

	totalVal := seasonValue * seasonMult + careerValue * careerMult // + awardsValue..

	totalCapcity := 10 // TODO : Should get this from config
	demandScaling := 0.4 // TODO : This too
	demand := float64((totalCapcity  - capacityRemaining) / totalCapcity)
	demandMult := 1 + smoothStep(demand) * demandScaling
	returnedValue := totalVal * demandMult
	return returnedValue, nil

}

func smoothStep(x float64) float64 {
	return x*x*(3-2*x)
}

func(s *ValueService) UpdatePlayerValue(ctx context.Context, playerID int64) error {
	season := "2025-26" // TODO: Fix this jank with config
	value, err := s.CalculateValueBasedOnStats(ctx, playerID, season)
	if err != nil {
		return err
	}
	return s.PlayerRepo.UpdateValue(ctx, playerID, value)
}


func(s *ValueService) UpdateValueForAllPlayers(ctx context.Context) error {
	season := "2025-26" // TODO: Fix this jank with config
	allIDs, err := s.PlayerRepo.GetAllIDs(ctx)
	if err != nil {
		return err
	}
	updates := make(map[int64]float64, len(allIDs))
	for _, id := range allIDs {
		v, err := s.CalculateValueBasedOnStats(ctx, id, season)
		if err != nil {
			return err
		}
		updates[id] = v
	}

	return s.PlayerRepo.UpdateAllValues(ctx, updates)
}

