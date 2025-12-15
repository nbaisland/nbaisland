package nba

import (
    "context"
    "fmt"
    "time"
    
    "github.com/n-ae/nba-api-go/pkg/stats"
    "github.com/n-ae/nba-api-go/pkg/stats/endpoints"
    "github.com/n-ae/nba-api-go/pkg/stats/parameters"
    "github.com/n-ae/nba-api-go/pkg/stats/static"
)


type Client struct {
    statsClient *stats.Client
}

func NewClient() *Client {
    return &Client{
        statsClient: stats.NewDefaultClient(),
    }
}

func (c *Client) GetActivePlayers() ([]static.Player, error) {
    activePlayers, err := static.GetActivePlayers()
    if err != nil {
        return nil, fmt.Errorf("Could not get active players: %w", err)
    }
    
    if activePlayers == nil {
        return nil, fmt.Errorf("No active players found")
    }
    
    return activePlayers, nil
}

func (c *Client) GetPlayerGameLog(ctx context.Context, playerID int, season string) ([]endpoints.GameLog, error) {
    playerString := fmt.Sprintf("%d", playerID)
    
    req := endpoints.PlayerGameLogRequest{
        PlayerID:   playerString,
        Season:     parameters.Season(season),
        SeasonType: "Regular Season",
        LeagueID:   parameters.LeagueIDNBA,
    }
    
    resp, err := endpoints.PlayerGameLog(ctx, c.statsClient, req)
    if err != nil {
        return nil, fmt.Errorf("failed to get game log for player %d: %w", playerID, err)
    }
    gameLogData := resp.Data.PlayerGameLog
    return gameLogData, nil
}

func (c *Client) GetPlayerGameLogDateRange(ctx context.Context, playerID int, season string, dateFrom, dateTo time.Time) ([]endpoints.GameLog, error) {
    playerString := fmt.Sprintf("%d", playerID)
    
    dateFromStr := dateFrom.Format("01/02/2006")
    dateToStr := dateTo.Format("01/02/2006")
    
    req := endpoints.PlayerGameLogRequest{
        PlayerID:   playerString,
        Season:     parameters.Season(season),
        SeasonType: "Regular Season",
        LeagueID:   parameters.LeagueIDNBA,
        DateFrom:   dateFromStr,
        DateTo:     dateToStr,
    }
    
    resp, err := endpoints.PlayerGameLog(ctx, c.statsClient, req)
    if err != nil {
        return nil, fmt.Errorf("failed to get game log for player %d: %w", playerID, err)
    }
    gameLogData := resp.Data.PlayerGameLog
    
    return gameLogData, nil
}

func AggregateSeasonStats(playerID int, playerName string, season string, games []endpoints.GameLog) *PlayerSeasonStats {
    if len(games) == 0 {
        return &PlayerSeasonStats{
            PlayerID:   playerID,
            PlayerName: playerName,
            Season:     season,
        }
    }
    
    var totalPts, totalReb, totalAst, totalStl, totalBlk int
    
    for _, game := range games {
        totalPts += game.PTS
        totalReb += game.REB
        totalAst += game.AST
        totalStl += game.STL
        totalBlk += game.BLK
    }
    
    gamesPlayed := len(games)
    
    return &PlayerSeasonStats{
        PlayerID:        playerID,
        PlayerName:      playerName,
        Season:          season,
        GamesPlayed:     gamesPlayed,
        TotalPoints:     totalPts,
        TotalRebounds:   totalReb,
        TotalAssists:    totalAst,
        TotalSteals:     totalStl,
        TotalBlocks:     totalBlk,
        PointsPerGame:   float64(totalPts) / float64(gamesPlayed),
        ReboundsPerGame: float64(totalReb) / float64(gamesPlayed),
        AssistsPerGame:  float64(totalAst) / float64(gamesPlayed),
        StealsPerGame:   float64(totalStl) / float64(gamesPlayed),
        BlocksPerGame:   float64(totalBlk) / float64(gamesPlayed),
    }
}

func AggregateWeeklyStats(playerID int, playerName string, season string, weekStart, weekEnd time.Time, games []endpoints.GameLog) *WeeklyStats {
    if len(games) == 0 {
        return &WeeklyStats{
            PlayerID:   playerID,
            PlayerName: playerName,
            Season:     season,
            WeekStart:  weekStart,
            WeekEnd:    weekEnd,
        }
    }
    
    var totalPts, totalReb, totalAst, totalStl, totalBlk int
    
    for _, game := range games {
        totalPts += game.PTS
        totalReb += game.REB
        totalAst += game.AST
        totalStl += game.STL
        totalBlk += game.BLK
    }
    
    gamesPlayed := len(games)
    
    return &WeeklyStats{
        PlayerID:        playerID,
        PlayerName:      playerName,
        Season:          season,
        WeekStart:       weekStart,
        WeekEnd:         weekEnd,
        GamesPlayed:     gamesPlayed,
        TotalPoints:     totalPts,
        TotalRebounds:   totalReb,
        TotalAssists:    totalAst,
        TotalSteals:     totalStl,
        TotalBlocks:     totalBlk,
        PointsPerGame:   float64(totalPts) / float64(gamesPlayed),
        ReboundsPerGame: float64(totalReb) / float64(gamesPlayed),
        AssistsPerGame:  float64(totalAst) / float64(gamesPlayed),
        StealsPerGame:   float64(totalStl) / float64(gamesPlayed),
        BlocksPerGame:   float64(totalBlk) / float64(gamesPlayed),
    }
}

func (c *Client) GetPlayerSeasonStats(ctx context.Context, playerID int, playerName string, season string) (*PlayerSeasonStats, error) {
    games, err := c.GetPlayerGameLog(ctx, playerID, season)
    if err != nil {
        return nil, err
    }
    
    return AggregateSeasonStats(playerID, playerName, season, games), nil
}

func (c *Client) GetPlayerWeeklyStats(ctx context.Context, playerID int, playerName string, season string) (*WeeklyStats, error) {
    weekEnd := time.Now()
    weekStart := weekEnd.AddDate(0, 0, -7)
    
    games, err := c.GetPlayerGameLogDateRange(ctx, playerID, season, weekStart, weekEnd)
    if err != nil {
        return nil, err
    }
    
    return AggregateWeeklyStats(playerID, playerName, season, weekStart, weekEnd, games), nil
}

func (c *Client) GetAllPlayersSeasonStats(ctx context.Context, season string) ([]PlayerSeasonStats, error) {
    players, err := c.GetActivePlayers()
    if err != nil {
        return nil, err
    }
    
    fmt.Printf("Fetching season stats for %d players...\n", len(players))
    
    var allStats []PlayerSeasonStats
    
    for i, player := range players {
        if i%50 == 0 {
            fmt.Printf("Progress: %d/%d players\n", i, len(players))
        }
        
        stats, err := c.GetPlayerSeasonStats(ctx, player.ID, player.FullName, season)
        if err != nil {
            fmt.Printf("Warning: Failed for player %s (ID: %d): %v\n", player.FullName, player.ID, err)
            continue
        }
        
        if stats.GamesPlayed > 0 {
            allStats = append(allStats, *stats)
        }
        
        // Rate limit delay
        time.Sleep(100 * time.Millisecond)
    }
    
    fmt.Printf("Successfully retrieved stats for %d players\n", len(allStats))
    return allStats, nil
}

func (c *Client) GetAllPlayersWeeklyStats(ctx context.Context, season string) ([]WeeklyStats, error) {
    players, err := c.GetActivePlayers()
    if err != nil {
        return nil, err
    }
    
    fmt.Printf("Fetching weekly stats for %d players...\n", len(players))
    
    var allStats []WeeklyStats
    
    for i, player := range players {
        if i%50 == 0 {
            fmt.Printf("Progress: %d/%d players\n", i, len(players))
        }
        
        stats, err := c.GetPlayerWeeklyStats(ctx, player.ID, player.FullName, season)
        if err != nil {
            fmt.Printf("Warning: Failed for player %s (ID: %d): %v\n", player.FullName, player.ID, err)
            continue
        }
        
        if stats.GamesPlayed > 0 {
            allStats = append(allStats, *stats)
        }
        
        // rate limit delay
        time.Sleep(100 * time.Millisecond)
    }
    
    fmt.Printf("Successfully retrieved weekly stats for %d players\n", len(allStats))
    return allStats, nil
}

func (c *Client) GetCustomDateRangeStats(ctx context.Context, playerID int, playerName string, season string, dateFrom, dateTo time.Time) (*WeeklyStats, error) {
    games, err := c.GetPlayerGameLogDateRange(ctx, playerID, season, dateFrom, dateTo)
    if err != nil {
        return nil, err
    }
    
    return AggregateWeeklyStats(playerID, playerName, season, dateFrom, dateTo, games), nil
}


func PrintSeasonStats(stats PlayerSeasonStats) {
    fmt.Printf("\n=== %s (ID: %d) - %s Season ===\n", stats.PlayerName, stats.PlayerID, stats.Season)
    fmt.Printf("Games Played: %d\n", stats.GamesPlayed)
    fmt.Printf("PPG: %.1f (%d total)\n", stats.PointsPerGame, stats.TotalPoints)
    fmt.Printf("RPG: %.1f (%d total)\n", stats.ReboundsPerGame, stats.TotalRebounds)
    fmt.Printf("APG: %.1f (%d total)\n", stats.AssistsPerGame, stats.TotalAssists)
    fmt.Printf("SPG: %.1f (%d total)\n", stats.StealsPerGame, stats.TotalSteals)
    fmt.Printf("BPG: %.1f (%d total)\n", stats.BlocksPerGame, stats.TotalBlocks)
}

func PrintWeeklyStats(stats WeeklyStats) {
    fmt.Printf("\n=== %s (ID: %d) - Week of %s ===\n", 
        stats.PlayerName, stats.PlayerID, stats.WeekStart.Format("Jan 2"))
    fmt.Printf("Games Played: %d\n", stats.GamesPlayed)
    fmt.Printf("PPG: %.1f (%d total)\n", stats.PointsPerGame, stats.TotalPoints)
    fmt.Printf("RPG: %.1f (%d total)\n", stats.ReboundsPerGame, stats.TotalRebounds)
    fmt.Printf("APG: %.1f (%d total)\n", stats.AssistsPerGame, stats.TotalAssists)
    fmt.Printf("SPG: %.1f (%d total)\n", stats.StealsPerGame, stats.TotalSteals)
    fmt.Printf("BPG: %.1f (%d total)\n", stats.BlocksPerGame, stats.TotalBlocks)
}