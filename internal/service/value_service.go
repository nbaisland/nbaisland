package service

import ( 
	"context"

	"github.com/nbaisland/nbaisland/internal/models"
    "github.com/nbaisland/nbaisland/internal/repository"
    "github.com/nbaisland/nbaisland/internal/nba"
)


type ValueService struct {
	PlayerRepo repository.PlayerRepository
	NBARepo nba.Repository
	playerMap repository.PlayerIDMapRepository
}

func NewValueService(playerRepo repository.PlayerRepository, nbaRepo nba.Repository) *ValueService {
	return &ValueService{PlayerRepo: playerRepo, NBARepo: nbaRepo}
}

func(s *ValueService) UpdateValue(ctx context.Context, playerID int64) error {
	s.Repo
}

// General thinking here is that value should be added to players who:
// Win, put up good stats, play minutes , Play games (rock solid workhorses should be more valuable than glass cannons)
// 
// 
func(s *ValueService) CalculateValueBasedOnStats(ctx context.Context, playerID int64, season string) (float64, error) {
	nbaID, err := s.playerMap.GetNBAPlayerByAppID(ctx, playerID)
	if err != nil {
		return 0, err
	}


	seasonStats := s.NBARepo.GetSeasonStats(ctx context.Context, nbaID, season)
	if err != nil {
		return 0, err
	}
	careerStats := s.NBARepo.GetCareerStats(ctx context.Context, nbaID)
	if err != nil {
		return 0, err
	}
	capacityRemaining, err := s.PlayerRepo.GetCapacityByID(playerID)
	if err != nil {
		return 0, err
	}
	// TODO: Add getting player tags, other information like champsionships, mvps, defense, all nba into value equation 
	// awardsValue := Championships * 2
	if seasonStats.GamesPlayed > 10 {
		seasonValue := seasonStats.PointsPerGame + seasonStats.AssistsPerGame * 2 + seasonStats.ReboundsPerGame * 2 + seasonStats.StealsPerGame * 3 + seasonStats.BlocksPerGame * 3
	}
	else {
		seasonValue := 0
	}
	careerValue := careerStats.PointsTotal * 0.001 + careerStats.ReboundsTotal * 0.002 + careerStats.AssistsTotal * 0.002 + careerStats.StealsTotal * 0.0025 + careerStats.BlocksTotal * 0.0025 + careerStats.MinutesTotal * 0.00001
	seasonMult := 1
	careerMult := 1

	totalCapcity := 10 // TODO : Should get this from config
	capacityMult := 0.25
	totalVal := seasonValue * seasonMult + careerValue * careerMult
	demand := ((totalCapcity  - capacityRemaining) * capacity) ^ 2
	demandMult := max(1.0, demand)
	return 

}

// func(s *ValueService) UpdateValueForAllPlayers(ctx) error {

// }

