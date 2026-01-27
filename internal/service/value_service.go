package service

import (
    "context"
    "sync"
    "github.com/nbaisland/nbaisland/internal/logger"
    "github.com/nbaisland/nbaisland/internal/nba"
    "github.com/nbaisland/nbaisland/internal/repository"
    "go.uber.org/zap"
)

type ValueWeights struct {
    SeasonPPG    float64
    SeasonAPG    float64
    SeasonRPG    float64
    SeasonSPG    float64
    SeasonBPG    float64
    
    CareerPoints   float64
    CareerRebounds float64
    CareerAssists  float64
    CareerSteals   float64
    CareerBlocks   float64
    CareerMinutes  float64
    
    SeasonMult float64
    CareerMult float64
    
    DemandScaling  float64
    TotalCapacity  int
    MinGamesPlayed int
}

func DefaultValueWeights() ValueWeights {
    return ValueWeights{
        SeasonPPG: 1.0,
        SeasonAPG: 2.0,
        SeasonRPG: 2.0,
        SeasonSPG: 3.0,
        SeasonBPG: 3.0,
        
        CareerPoints:   0.001,
        CareerRebounds: 0.002,
        CareerAssists:  0.002,
        CareerSteals:   0.0025,
        CareerBlocks:   0.0025,
        CareerMinutes:  0.00001,
        
        SeasonMult: 1.0,
        CareerMult: 1.0,
        
        DemandScaling:  0.4,
        TotalCapacity:  10,
        MinGamesPlayed: 10,
    }
}

type ValueService struct {
    PlayerRepo    repository.PlayerRepository
    NBARepo       *nba.Repository
    PlayerMapRepo repository.PlayerIDMapRepository
    Weights       ValueWeights
}

func NewValueService(playerRepo repository.PlayerRepository, nbaRepo *nba.Repository, playerMapRepo repository.PlayerIDMapRepository) *ValueService {
    return &ValueService{
        PlayerRepo:    playerRepo,
        NBARepo:       nbaRepo,
        PlayerMapRepo: playerMapRepo,
        Weights:       DefaultValueWeights(),
    }
}

func (s *ValueService) CalculateValueBasedOnStats(ctx context.Context, playerID int64, season string) (float64, error) {
    nbaID, err := s.PlayerMapRepo.GetNBAPlayerByAppID(ctx, playerID)
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
    
    var seasonValue float64
    if seasonStats != nil && seasonStats.GamesPlayed > s.Weights.MinGamesPlayed {
        seasonValue = (seasonStats.PointsPerGame * s.Weights.SeasonPPG) +
                     (seasonStats.AssistsPerGame * s.Weights.SeasonAPG) +
                     (seasonStats.ReboundsPerGame * s.Weights.SeasonRPG) +
                     (seasonStats.StealsPerGame * s.Weights.SeasonSPG) +
                     (seasonStats.BlocksPerGame * s.Weights.SeasonBPG)
    }
    
    var careerValue float64
    if careerStats != nil {
        careerValue = (careerStats.PointsTotal * s.Weights.CareerPoints) +
                     (careerStats.ReboundsTotal * s.Weights.CareerRebounds) +
                     (careerStats.AssistsTotal * s.Weights.CareerAssists) +
                     (careerStats.StealsTotal * s.Weights.CareerSteals) +
                     (careerStats.BlocksTotal * s.Weights.CareerBlocks) +
                     (careerStats.MinutesTotal * s.Weights.CareerMinutes)
    }
    
    totalVal := (seasonValue * s.Weights.SeasonMult) + (careerValue * s.Weights.CareerMult)
    
    demand := float64(s.Weights.TotalCapacity-capacityRemaining) / float64(s.Weights.TotalCapacity)
    demandMult := 1 + smoothStep(demand)*s.Weights.DemandScaling
    
    returnedValue := totalVal * demandMult
    
    if returnedValue < 10.0 {
        returnedValue = 10.0
    }
    
    return returnedValue, nil
}

func smoothStep(x float64) float64 {
    return x * x * (3 - 2*x)
}

func (s *ValueService) UpdatePlayerValue(ctx context.Context, playerID int64, season string) error {
    value, err := s.CalculateValueBasedOnStats(ctx, playerID, season)
    if err != nil {
        return err
    }
    return s.PlayerRepo.UpdateValue(ctx, playerID, value)
}

func (s *ValueService) UpdateValueForAllPlayers(ctx context.Context, season string) error {
    logger.Log.Info("Starting value update for all players", zap.String("season", season))
    
    allIDs, err := s.PlayerRepo.GetAllIDs(ctx)
    if err != nil {
        return err
    }
    
    logger.Log.Info("Calculating values", zap.Int("player_count", len(allIDs)))
    
    type result struct {
        id    int64
        value float64
        err   error
    }
    
    workers := 10
    jobs := make(chan int64, len(allIDs))
    results := make(chan result, len(allIDs))
    
    var wg sync.WaitGroup
    for i := 0; i < workers; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for id := range jobs {
                value, err := s.CalculateValueBasedOnStats(ctx, id, season)
                results <- result{id: id, value: value, err: err}
            }
        }()
    }
    
    for _, id := range allIDs {
        jobs <- id
    }
    close(jobs)
    
    go func() {
        wg.Wait()
        close(results)
    }()
    
    updates := make(map[int64]float64)
    failedCount := 0
    for r := range results {
        if r.err != nil {
            logger.Log.Warn("Failed to calculate value for player",
                zap.Int64("player_id", r.id),
                zap.Error(r.err),
            )
            failedCount++
            continue
        }
        updates[r.id] = r.value
    }
    
    logger.Log.Info("Value calculation complete",
        zap.Int("successful", len(updates)),
        zap.Int("failed", failedCount),
    )
    
    if err := s.PlayerRepo.UpdateAllValues(ctx, updates); err != nil {
        return err
    }
    
    logger.Log.Info("Player values updated successfully",
        zap.Int("updated", len(updates)),
    )
    
    return nil
}

func (s *ValueService) UpdateValueForPlayers(ctx context.Context, playerIDs []int64, season string) error {
    logger.Log.Info("Updating values for specific players",
        zap.Int("count", len(playerIDs)),
        zap.String("season", season),
    )
    
    updates := make(map[int64]float64, len(playerIDs))
    
    for _, id := range playerIDs {
        value, err := s.CalculateValueBasedOnStats(ctx, id, season)
        if err != nil {
            logger.Log.Warn("Failed to calculate value",
                zap.Int64("player_id", id),
                zap.Error(err),
            )
            continue
        }
        updates[id] = value
    }
    
    if len(updates) == 0 {
        return nil
    }
    
    return s.PlayerRepo.UpdateAllValues(ctx, updates)
}
