package nba

import (
    "context"
    "fmt"
    "log"
    
    "github.com/nbaisland/nbaisland/internal/utils"
    "github.com/jackc/pgx/v5/pgxpool"
)

type NBAService struct {
    client *Client
    repo   *Repository
    pool   *pgxpool.Pool
}

func NewNBAService(client *Client, repo *Repository, pool *pgxpool.Pool) *NBAService {
    return &NBAService{
        client: client,
        repo:   repo,
        pool:   pool,
    }
}

func (s *NBAService) UpdateAllSeasonStats(ctx context.Context, season string) error {
    log.Printf("Starting season stats update for %s...", season)
    
    players, err := s.client.GetActivePlayers()
    if err != nil {
        return fmt.Errorf("failed to get active players: %w", err)
    }
    
    log.Printf("Found %d active players", len(players))
    
    for i, player := range players {
        if i%100 == 0 {
            log.Printf("Saving players: %d/%d", i, len(players))
        }
        
        err := s.repo.UpsertPlayer(ctx, &Player{
            ID:        int64(player.ID),
            FullName:  player.FullName,
            FirstName: player.FirstName,
            LastName:  player.LastName,
            IsActive:  true,
        })
        if err != nil {
            log.Printf("Warning: Failed to save player %s (ID: %d): %v", player.FullName, player.ID, err)
        }
    }
    
    log.Println("Fetching season stats for all players...")
    
    allStats, err := s.client.GetAllPlayersSeasonStats(ctx, season)
    if err != nil {
        return fmt.Errorf("failed to get season stats: %w", err)
    }
    
    log.Printf("Retrieved stats for %d players, saving to database...", len(allStats))
    
    batchSize := 50
    for i := 0; i < len(allStats); i += batchSize {
        end := i + batchSize
        if end > len(allStats) {
            end = len(allStats)
        }
        
        batch := allStats[i:end]
        if err := s.repo.BatchSaveSeasonStats(ctx, batch); err != nil {
            log.Printf("Warning: Failed to save batch starting at index %d: %v", i, err)
            continue
        }
        
        log.Printf("Saved batch: %d/%d players", end, len(allStats))
    }
    
    log.Printf("Season stats update completed! Saved %d players", len(allStats))
    return nil
}

func (s *NBAService) UpdateAllWeeklyStats(ctx context.Context, season string) error {
    log.Printf("Starting weekly stats update for %s...", season)
    
    players, err := s.client.GetActivePlayers()
    if err != nil {
        return fmt.Errorf("failed to get active players: %w", err)
    }
    
    log.Printf("Found %d active players, fetching weekly stats...", len(players))
    
    allStats, err := s.client.GetAllPlayersWeeklyStats(ctx, season)
    if err != nil {
        return fmt.Errorf("failed to get weekly stats: %w", err)
    }
    
    log.Printf("Retrieved weekly stats for %d players (players who played this week)", len(allStats))
    
    batchSize := 50
    for i := 0; i < len(allStats); i += batchSize {
        end := i + batchSize
        if end > len(allStats) {
            end = len(allStats)
        }
        
        batch := allStats[i:end]
        if err := s.repo.BatchSaveWeeklyStats(ctx, batch); err != nil {
            log.Printf("Warning: Failed to save weekly batch starting at index %d: %v", i, err)
            continue
        }
        
        log.Printf("Saved batch: %d/%d players", end, len(allStats))
    }
    
    log.Printf("Weekly stats update completed! Saved %d players", len(allStats))
    return nil
}

func (s *NBAService) UpdatePlayerSeasonStats(ctx context.Context, playerID int64, season string) error {
    players, err := s.client.GetActivePlayers()
    if err != nil {
        return fmt.Errorf("failed to get players: %w", err)
    }
    
    var playerName string
    for _, p := range players {
        if int64(p.ID) == playerID {
            playerName = p.FullName
            
            err := s.repo.UpsertPlayer(ctx, &Player{
                ID:        int64(p.ID),
                FullName:  p.FullName,
                FirstName: p.FirstName,
                LastName:  p.LastName,
                IsActive:  true,
            })
            if err != nil {
                log.Printf("Warning: Failed to save player: %v", err)
            }
            break
        }
    }
    
    if playerName == "" {
        return fmt.Errorf("player %d not found", playerID)
    }
    
    stats, err := s.client.GetPlayerSeasonStats(ctx, playerID, playerName, season)
    if err != nil {
        return fmt.Errorf("failed to get player stats: %w", err)
    }
    
    if err := s.repo.SaveSeasonStats(ctx, stats); err != nil {
        return fmt.Errorf("failed to save stats: %w", err)
    }
    
    log.Printf("Updated season stats for %s", playerName)
    return nil
}

func (s *NBAService) UpdatePlayerWeeklyStats(ctx context.Context, playerID int64, season string) error {
    players, err := s.client.GetActivePlayers()
    if err != nil {
        return fmt.Errorf("failed to get players: %w", err)
    }
    
    var playerName string
    for _, p := range players {
        if int64(p.ID) == playerID {
            playerName = p.FullName
            break
        }
    }
    
    if playerName == "" {
        return fmt.Errorf("player %d not found", playerID)
    }
    
    stats, err := s.client.GetPlayerWeeklyStats(ctx, playerID, playerName, season)
    if err != nil {
        return fmt.Errorf("failed to get weekly stats: %w", err)
    }
    
    if err := s.repo.SaveWeeklyStats(ctx, stats); err != nil {
        return fmt.Errorf("failed to save weekly stats: %w", err)
    }
    
    log.Printf("Updated weekly stats for %s", playerName)
    return nil
}

func (s *NBAService) GetPlayerStats(ctx context.Context, playerID int64, season string) (*PlayerStatsResponse, error) {
    seasonStats, err := s.repo.GetSeasonStats(ctx, playerID, season)
    if err != nil {
        return nil, fmt.Errorf("failed to get season stats: %w", err)
    }
    
    weeklyStats, err := s.repo.GetLatestWeeklyStats(ctx, playerID)
    if err != nil {
        return nil, fmt.Errorf("failed to get weekly stats: %w", err)
    }

    careerStats, err := s.repo.GetCareerStats(ctx, playerID)
    
    return &PlayerStatsResponse{
        SeasonStats: seasonStats,
        WeeklyStats: weeklyStats,
        CareerStats: careerStats,
    }, nil
}


func (s *NBAService) SeedTopPlayers(ctx context.Context, season string, minGamesPlayed int) error {
    log.Printf("Seeding players with at least %d games played...", minGamesPlayed)
    
    allStats, err := s.client.GetAllPlayersSeasonStats(ctx, season)
    if err != nil {
        return fmt.Errorf("failed to get season stats: %w", err)
    }
    
    seededCount := 0
    
    for _, stats := range allStats {
        if stats.GamesPlayed < minGamesPlayed {
            continue
        }
        
        err := s.repo.UpsertPlayer(ctx, &Player{
            ID:        stats.PlayerID,
            FullName:  stats.PlayerName,
            IsActive:  true,
        })
        if err != nil {
            log.Printf("Warning: Failed to save player %d: %v", stats.PlayerID, err)
            continue
        }
        
        var exists bool
        err = s.pool.QueryRow(ctx, `
            SELECT EXISTS(SELECT 1 FROM player_nba_mapping WHERE nba_player_id = $1)
        `, stats.PlayerID).Scan(&exists)
        
        if err != nil || exists {
            continue
        }
        
        // TODO: Update value better
        value := stats.PointsPerGame * 10
        capacity := 10
        slug := utils.ToSlug(stats.PlayerName)
        
        var appPlayerID int
        err = s.pool.QueryRow(ctx, `
            INSERT INTO players (name, value, capacity, slug)
            VALUES ($1, $2, $3, $4)
            RETURNING id
        `, stats.PlayerName, value, capacity, slug).Scan(&appPlayerID)
        
        if err != nil {
            log.Printf("Warning: Failed to create app player for %s: %v", stats.PlayerName, err)
            continue
        }
        
        _, err = s.pool.Exec(ctx, `
            INSERT INTO player_nba_mapping (player_id, nba_player_id)
            VALUES ($1, $2)
        `, appPlayerID, stats.PlayerID)
        
        if err != nil {
            log.Printf("Warning: Failed to create mapping for %s: %v", stats.PlayerName, err)
            continue
        }
        
        seededCount++
        if seededCount%50 == 0 {
            log.Printf("Seeded %d players so far...", seededCount)
        }
    }
    
    log.Printf("Seeding complete! Created %d app players", seededCount)
    return nil
}

func (s *NBAService) GetAppPlayerWithStats(ctx context.Context, appPlayerID int, season string) (*AppPlayerWithStats, error) {
    query := `
        SELECT 
            p.id, p.name, p.value, p.capacity, p.slug,
            n.id, n.full_name,
            COALESCE(s.games_played, 0), 
            COALESCE(s.points_per_game, 0), 
            COALESCE(s.rebounds_per_game, 0),
            COALESCE(s.assists_per_game, 0), 
            COALESCE(s.steals_per_game, 0), 
            COALESCE(s.blocks_per_game, 0)
        FROM players p
        JOIN player_nba_mapping m ON p.id = m.player_id
        JOIN nba_players n ON m.nba_player_id = n.id
        LEFT JOIN nba_season_stats s ON n.id = s.player_id AND s.season = $2
        WHERE p.id = $1
    `
    
    var result AppPlayerWithStats
    err := s.pool.QueryRow(ctx, query, appPlayerID, season).Scan(
        &result.AppPlayerID,
        &result.Name,
        &result.Value,
        &result.Capacity,
        &result.Slug,
        &result.NBAPlayerID,
        &result.NBAPlayerName,
        &result.GamesPlayed,
        &result.PointsPerGame,
        &result.ReboundsPerGame,
        &result.AssistsPerGame,
        &result.StealsPerGame,
        &result.BlocksPerGame,
    )
    
    if err != nil {
        return nil, fmt.Errorf("failed to get app player with stats: %w", err)
    }
    
    return &result, nil
}

func (s *NBAService) UpdateAllCareerStats(ctx context.Context) error {
    log.Printf("Starting career stats update...")
    
    allStats, err := s.client.GetAllPlayersCareerStats(ctx)
    if err != nil {
        return fmt.Errorf("failed to get career stats: %w", err)
    }
    
    log.Printf("Retrieved career stats for %d players, saving to database...", len(allStats))
    
    batchSize := 50
    for i := 0; i < len(allStats); i += batchSize {
        end := i + batchSize
        if end > len(allStats) {
            end = len(allStats)
        }
        
        batch := allStats[i:end]
        if err := s.repo.BatchSaveCareerStats(ctx, batch); err != nil {
            log.Printf("Warning: Failed to save batch starting at index %d: %v", i, err)
            continue
        }
        
        log.Printf("Saved batch: %d/%d players", end, len(allStats))
    }
    
    log.Printf("Career stats update completed! Saved %d players", len(allStats))
    return nil
}

func (s *NBAService) GetPlayerCareerStats(ctx context.Context, playerID int64) (*PlayerCareerStats, error) {
    return s.repo.GetCareerStats(ctx, playerID)
}

type PlayerStatsResponse struct {
    SeasonStats *PlayerSeasonStats
    WeeklyStats *WeeklyStats
    CareerStats *PlayerCareerStats
}

type AppPlayerWithStats struct {
    AppPlayerID int
    Name        string
    Value       float64
    Capacity    int
    Slug        string
    
    NBAPlayerID   int
    NBAPlayerName string
    
    GamesPlayed     int
    PointsPerGame   float64
    ReboundsPerGame float64
    AssistsPerGame  float64
    StealsPerGame   float64
    BlocksPerGame   float64
}