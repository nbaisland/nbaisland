package nba

import (
    "context"
    "errors"
    "fmt"
    "time"
    
    "github.com/jackc/pgx/v5"
    "github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
    pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
    return &Repository{pool: pool}
}

func (r *Repository) UpsertPlayer(ctx context.Context, player *Player) error {
    query := `
        INSERT INTO nba_players (id, full_name, first_name, last_name, is_active, updated_at)
        VALUES ($1, $2, $3, $4, $5, NOW())
        ON CONFLICT (id) DO UPDATE SET
            full_name = EXCLUDED.full_name,
            first_name = EXCLUDED.first_name,
            last_name = EXCLUDED.last_name,
            is_active = EXCLUDED.is_active,
            updated_at = NOW()
    `
    
    _, err := r.pool.Exec(ctx, query, player.ID, player.FullName, player.FirstName, player.LastName, player.IsActive)
    return err
}

func (r *Repository) GetNBAPlayer(ctx context.Context, playerID int) (*Player, error) {
    query := `
        SELECT id, full_name, first_name, last_name, is_active, created_at, updated_at
        FROM nba_players
        WHERE id = $1
    `
    
    var player Player
    var createdAt, updatedAt time.Time
    
    err := r.pool.QueryRow(ctx, query, playerID).Scan(
        &player.ID,
        &player.FullName,
        &player.FirstName,
        &player.LastName,
        &player.IsActive,
        &createdAt,
        &updatedAt,
    )
    
    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return nil, nil
        }
        return nil, err
    }
    
    return &player, nil
}

func (r *Repository) SaveSeasonStats(ctx context.Context, stats *PlayerSeasonStats) error {
    query := `
        INSERT INTO nba_season_stats 
            (player_id, season, games_played, total_points, total_rebounds, 
             total_assists, total_steals, total_blocks, points_per_game, 
             rebounds_per_game, assists_per_game, steals_per_game, blocks_per_game, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, NOW())
        ON CONFLICT (player_id, season) DO UPDATE SET
            games_played = EXCLUDED.games_played,
            total_points = EXCLUDED.total_points,
            total_rebounds = EXCLUDED.total_rebounds,
            total_assists = EXCLUDED.total_assists,
            total_steals = EXCLUDED.total_steals,
            total_blocks = EXCLUDED.total_blocks,
            points_per_game = EXCLUDED.points_per_game,
            rebounds_per_game = EXCLUDED.rebounds_per_game,
            assists_per_game = EXCLUDED.assists_per_game,
            steals_per_game = EXCLUDED.steals_per_game,
            blocks_per_game = EXCLUDED.blocks_per_game,
            updated_at = NOW()
    `
    
    _, err := r.pool.Exec(ctx, query,
        stats.PlayerID,
        stats.Season,
        stats.GamesPlayed,
        stats.TotalPoints,
        stats.TotalRebounds,
        stats.TotalAssists,
        stats.TotalSteals,
        stats.TotalBlocks,
        stats.PointsPerGame,
        stats.ReboundsPerGame,
        stats.AssistsPerGame,
        stats.StealsPerGame,
        stats.BlocksPerGame,
    )
    
    return err
}

func (r *Repository) SaveWeeklyStats(ctx context.Context, stats *WeeklyStats) error {
    query := `
        INSERT INTO nba_weekly_stats 
            (player_id, season, week_start, week_end, games_played, 
             total_points, total_rebounds, total_assists, total_steals, total_blocks,
             points_per_game, rebounds_per_game, assists_per_game, 
             steals_per_game, blocks_per_game)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
        ON CONFLICT (player_id, week_end) DO UPDATE SET
            games_played = EXCLUDED.games_played,
            total_points = EXCLUDED.total_points,
            total_rebounds = EXCLUDED.total_rebounds,
            total_assists = EXCLUDED.total_assists,
            total_steals = EXCLUDED.total_steals,
            total_blocks = EXCLUDED.total_blocks,
            points_per_game = EXCLUDED.points_per_game,
            rebounds_per_game = EXCLUDED.rebounds_per_game,
            assists_per_game = EXCLUDED.assists_per_game,
            steals_per_game = EXCLUDED.steals_per_game,
            blocks_per_game = EXCLUDED.blocks_per_game
    `
    
    _, err := r.pool.Exec(ctx, query,
        stats.PlayerID,
        stats.Season,
        stats.WeekStart,
        stats.WeekEnd,
        stats.GamesPlayed,
        stats.TotalPoints,
        stats.TotalRebounds,
        stats.TotalAssists,
        stats.TotalSteals,
        stats.TotalBlocks,
        stats.PointsPerGame,
        stats.ReboundsPerGame,
        stats.AssistsPerGame,
        stats.StealsPerGame,
        stats.BlocksPerGame,
    )
    
    return err
}

func (r *Repository) GetSeasonStats(ctx context.Context, playerID int64, season string) (*PlayerSeasonStats, error) {
    query := `
        SELECT player_id, season, games_played, total_points, total_rebounds,
               total_assists, total_steals, total_blocks, points_per_game,
               rebounds_per_game, assists_per_game, steals_per_game, blocks_per_game
        FROM nba_season_stats
        WHERE player_id = $1 AND season = $2
    `
    
    var stats PlayerSeasonStats
    err := r.pool.QueryRow(ctx, query, playerID, season).Scan(
        &stats.PlayerID,
        &stats.Season,
        &stats.GamesPlayed,
        &stats.TotalPoints,
        &stats.TotalRebounds,
        &stats.TotalAssists,
        &stats.TotalSteals,
        &stats.TotalBlocks,
        &stats.PointsPerGame,
        &stats.ReboundsPerGame,
        &stats.AssistsPerGame,
        &stats.StealsPerGame,
        &stats.BlocksPerGame,
    )
    
    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return nil, nil
        }
        return nil, err
    }
    
    return &stats, nil
}

func (r *Repository) GetLatestWeeklyStats(ctx context.Context, playerID int64) (*WeeklyStats, error) {
    query := `
        SELECT player_id, season, week_start, week_end, games_played,
               total_points, total_rebounds, total_assists, total_steals, total_blocks,
               points_per_game, rebounds_per_game, assists_per_game,
               steals_per_game, blocks_per_game
        FROM nba_weekly_stats
        WHERE player_id = $1
        ORDER BY week_end DESC
        LIMIT 1
    `
    
    var stats WeeklyStats
    err := r.pool.QueryRow(ctx, query, playerID).Scan(
        &stats.PlayerID,
        &stats.Season,
        &stats.WeekStart,
        &stats.WeekEnd,
        &stats.GamesPlayed,
        &stats.TotalPoints,
        &stats.TotalRebounds,
        &stats.TotalAssists,
        &stats.TotalSteals,
        &stats.TotalBlocks,
        &stats.PointsPerGame,
        &stats.ReboundsPerGame,
        &stats.AssistsPerGame,
        &stats.StealsPerGame,
        &stats.BlocksPerGame,
    )
    
    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return nil, nil
        }
        return nil, err
    }
    
    return &stats, nil
}

func (r *Repository) GetAllActivePlayers(ctx context.Context) ([]Player, error) {
    query := `
        SELECT id, full_name, first_name, last_name, is_active
        FROM nba_players
        WHERE is_active = true
        ORDER BY full_name
    `
    
    rows, err := r.pool.Query(ctx, query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var players []Player
    for rows.Next() {
        var p Player
        err := rows.Scan(&p.ID, &p.FullName, &p.FirstName, &p.LastName, &p.IsActive)
        if err != nil {
            return nil, err
        }
        players = append(players, p)
    }
    
    if err := rows.Err(); err != nil {
        return nil, err
    }
    
    return players, nil
}

func (r *Repository) BatchSaveSeasonStats(ctx context.Context, allStats []PlayerSeasonStats) error {
    if len(allStats) == 0 {
        return nil
    }
    
    tx, err := r.pool.Begin(ctx)
    if err != nil {
        return fmt.Errorf("failed to begin transaction: %w", err)
    }
    defer tx.Rollback(ctx)
    
    query := `
        INSERT INTO nba_season_stats 
            (player_id, season, games_played, total_points, total_rebounds, 
             total_assists, total_steals, total_blocks, points_per_game, 
             rebounds_per_game, assists_per_game, steals_per_game, blocks_per_game, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, NOW())
        ON CONFLICT (player_id, season) DO UPDATE SET
            games_played = EXCLUDED.games_played,
            total_points = EXCLUDED.total_points,
            total_rebounds = EXCLUDED.total_rebounds,
            total_assists = EXCLUDED.total_assists,
            total_steals = EXCLUDED.total_steals,
            total_blocks = EXCLUDED.total_blocks,
            points_per_game = EXCLUDED.points_per_game,
            rebounds_per_game = EXCLUDED.rebounds_per_game,
            assists_per_game = EXCLUDED.assists_per_game,
            steals_per_game = EXCLUDED.steals_per_game,
            blocks_per_game = EXCLUDED.blocks_per_game,
            updated_at = NOW()
    `
    
    for _, stats := range allStats {
        _, err := tx.Exec(ctx, query,
            stats.PlayerID,
            stats.Season,
            stats.GamesPlayed,
            stats.TotalPoints,
            stats.TotalRebounds,
            stats.TotalAssists,
            stats.TotalSteals,
            stats.TotalBlocks,
            stats.PointsPerGame,
            stats.ReboundsPerGame,
            stats.AssistsPerGame,
            stats.StealsPerGame,
            stats.BlocksPerGame,
        )
        if err != nil {
            return fmt.Errorf("failed to insert stats for player %d: %w", stats.PlayerID, err)
        }
    }
    
    if err := tx.Commit(ctx); err != nil {
        return fmt.Errorf("failed to commit transaction: %w", err)
    }
    
    return nil
}

// BatchSaveWeeklyStats saves multiple weekly stats in a batch
func (r *Repository) BatchSaveWeeklyStats(ctx context.Context, allStats []WeeklyStats) error {
    if len(allStats) == 0 {
        return nil
    }
    
    tx, err := r.pool.Begin(ctx)
    if err != nil {
        return fmt.Errorf("failed to begin transaction: %w", err)
    }
    defer tx.Rollback(ctx)
    
    query := `
        INSERT INTO nba_weekly_stats 
            (player_id, season, week_start, week_end, games_played, 
             total_points, total_rebounds, total_assists, total_steals, total_blocks,
             points_per_game, rebounds_per_game, assists_per_game, 
             steals_per_game, blocks_per_game)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
        ON CONFLICT (player_id, week_end) DO UPDATE SET
            games_played = EXCLUDED.games_played,
            total_points = EXCLUDED.total_points,
            total_rebounds = EXCLUDED.total_rebounds,
            total_assists = EXCLUDED.total_assists,
            total_steals = EXCLUDED.total_steals,
            total_blocks = EXCLUDED.total_blocks,
            points_per_game = EXCLUDED.points_per_game,
            rebounds_per_game = EXCLUDED.rebounds_per_game,
            assists_per_game = EXCLUDED.assists_per_game,
            steals_per_game = EXCLUDED.steals_per_game,
            blocks_per_game = EXCLUDED.blocks_per_game
    `
    
    for _, stats := range allStats {
        _, err := tx.Exec(ctx, query,
            stats.PlayerID,
            stats.Season,
            stats.WeekStart,
            stats.WeekEnd,
            stats.GamesPlayed,
            stats.TotalPoints,
            stats.TotalRebounds,
            stats.TotalAssists,
            stats.TotalSteals,
            stats.TotalBlocks,
            stats.PointsPerGame,
            stats.ReboundsPerGame,
            stats.AssistsPerGame,
            stats.StealsPerGame,
            stats.BlocksPerGame,
        )
        if err != nil {
            return fmt.Errorf("failed to insert weekly stats for player %d: %w", stats.PlayerID, err)
        }
    }
    
    if err := tx.Commit(ctx); err != nil {
        return fmt.Errorf("failed to commit transaction: %w", err)
    }
    
    return nil
}

func (r *Repository) SaveCareerStats(ctx context.Context, stats *PlayerCareerStats) error {
    query := `
        INSERT INTO nba_career_stats 
            (player_id, games_played, points_per_game, rebounds_per_game, 
             assists_per_game, steals_per_game, blocks_per_game, 
             field_goal_pct, three_point_pct, free_throw_pct, minutes_per_game, updated_at,
             points_total, rebounds_total, assists_total, steals_total, blocks_total, minutes_total)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, NOW(), $12, $13, $14, $15, $16, $17)
        ON CONFLICT (player_id) DO UPDATE SET
            games_played = EXCLUDED.games_played,
            points_per_game = EXCLUDED.points_per_game,
            rebounds_per_game = EXCLUDED.rebounds_per_game,
            assists_per_game = EXCLUDED.assists_per_game,
            steals_per_game = EXCLUDED.steals_per_game,
            blocks_per_game = EXCLUDED.blocks_per_game,
            field_goal_pct = EXCLUDED.field_goal_pct,
            three_point_pct = EXCLUDED.three_point_pct,
            free_throw_pct = EXCLUDED.free_throw_pct,
            minutes_per_game = EXCLUDED.minutes_per_game,
            points_total = EXCLUDED.points_total,
            rebounds_total = EXCLUDED.rebounds_total,
            assists_total = EXCLUDED.assists_total,
            steals_total = EXCLUDED.steals_total,
            blocks_total = EXCLUDED.blocks_total,
            minutes_total = EXCLUDED.minutes_total,
            updated_at = NOW()
    `
    
    _, err := r.pool.Exec(ctx, query,
        stats.PlayerID,
        stats.GamesPlayed,
        stats.PointsPerGame,
        stats.ReboundsPerGame,
        stats.AssistsPerGame,
        stats.StealsPerGame,
        stats.BlocksPerGame,
        stats.FieldGoalPct,
        stats.ThreePointPct,
        stats.FreeThrowPct,
        stats.MinutesPerGame,
        stats.PointsTotal,
        stats.ReboundsTotal,
        stats.AssistsTotal,
        stats.StealsTotal,
        stats.BlocksTotal,
        stats.MinutesTotal,
    )
    
    return err
}

func (r *Repository) GetCareerStats(ctx context.Context, playerID int64) (*PlayerCareerStats, error) {
    query := `
        SELECT player_id, games_played, points_per_game, rebounds_per_game,
               assists_per_game, steals_per_game, blocks_per_game,
               field_goal_pct, three_point_pct, free_throw_pct, minutes_per_game,
               points_total, rebounds_total, assists_total, steals_total, blocks_total, minutes_total
        FROM nba_career_stats
        WHERE player_id = $1
    `
    
    var stats PlayerCareerStats
    err := r.pool.QueryRow(ctx, query, playerID).Scan(
        &stats.PlayerID,
        &stats.GamesPlayed,
        &stats.PointsPerGame,
        &stats.ReboundsPerGame,
        &stats.AssistsPerGame,
        &stats.StealsPerGame,
        &stats.BlocksPerGame,
        &stats.FieldGoalPct,
        &stats.ThreePointPct,
        &stats.FreeThrowPct,
        &stats.MinutesPerGame,
        &stats.PointsTotal,
        &stats.ReboundsTotal,
        &stats.AssistsTotal,
        &stats.StealsTotal,
        &stats.BlocksTotal,
        &stats.MinutesTotal,
    )
    
    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return nil, nil
        }
        return nil, err
    }
    
    return &stats, nil
}

func (r *Repository) BatchSaveCareerStats(ctx context.Context, allStats []PlayerCareerStats) error {
    if len(allStats) == 0 {
        return nil
    }
    
    tx, err := r.pool.Begin(ctx)
    if err != nil {
        return fmt.Errorf("failed to begin transaction: %w", err)
    }
    defer tx.Rollback(ctx)
    
    query := `
        INSERT INTO nba_career_stats 
            (player_id, games_played, points_per_game, rebounds_per_game, 
             assists_per_game, steals_per_game, blocks_per_game, 
             field_goal_pct, three_point_pct, free_throw_pct, minutes_per_game, updated_at,
             points_total, rebounds_total, assists_total, steals_total, blocks_total, minutes_total)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, NOW(), $12, $13, $14, $15, $16, $17)
        ON CONFLICT (player_id) DO UPDATE SET
            games_played = EXCLUDED.games_played,
            points_per_game = EXCLUDED.points_per_game,
            rebounds_per_game = EXCLUDED.rebounds_per_game,
            assists_per_game = EXCLUDED.assists_per_game,
            steals_per_game = EXCLUDED.steals_per_game,
            blocks_per_game = EXCLUDED.blocks_per_game,
            field_goal_pct = EXCLUDED.field_goal_pct,
            three_point_pct = EXCLUDED.three_point_pct,
            free_throw_pct = EXCLUDED.free_throw_pct,
            minutes_per_game = EXCLUDED.minutes_per_game,
            points_total = EXCLUDED.points_total,
            rebounds_total = EXCLUDED.rebounds_total,
            assists_total = EXCLUDED.assists_total,
            steals_total = EXCLUDED.steals_total,
            blocks_total = EXCLUDED.blocks_total,
            minutes_total = EXCLUDED.minutes_total,
            updated_at = NOW()
    `
    
    for _, stats := range allStats {
        _, err := tx.Exec(ctx, query,
            stats.PlayerID,
            stats.GamesPlayed,
            stats.PointsPerGame,
            stats.ReboundsPerGame,
            stats.AssistsPerGame,
            stats.StealsPerGame,
            stats.BlocksPerGame,
            stats.FieldGoalPct,
            stats.ThreePointPct,
            stats.FreeThrowPct,
            stats.MinutesPerGame,
            stats.PointsTotal,
            stats.ReboundsTotal,
            stats.AssistsTotal,
            stats.StealsTotal,
            stats.BlocksTotal,
            stats.MinutesTotal,
        )
        if err != nil {
            return fmt.Errorf("failed to insert career stats for player %d: %w", stats.PlayerID, err)
        }
    }
    
    if err := tx.Commit(ctx); err != nil {
        return fmt.Errorf("failed to commit transaction: %w", err)
    }
    
    return nil
}