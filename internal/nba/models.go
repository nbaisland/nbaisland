package nba

import "time"

type Player struct {
    ID         int
    FullName   string
    FirstName  string
    LastName   string
    IsActive   bool
    CreatedAt  time.Time
    UpdatedAt  time.Time
}

type PlayerSeasonStats struct {
    PlayerID        int
    PlayerName      string
    Season          string
    GamesPlayed     int
    TotalPoints     int
    TotalRebounds   int
    TotalAssists    int
    TotalSteals     int
    TotalBlocks     int
    PointsPerGame   float64
    ReboundsPerGame float64
    AssistsPerGame  float64
    StealsPerGame   float64
    BlocksPerGame   float64
}

type WeeklyStats struct {
    PlayerID        int
    PlayerName      string
    WeekStart       time.Time
    WeekEnd         time.Time
    Season          string
    GamesPlayed     int
    TotalPoints     int
    TotalRebounds   int
    TotalAssists    int
    TotalSteals     int
    TotalBlocks     int
    PointsPerGame   float64
    ReboundsPerGame float64
    AssistsPerGame  float64
    StealsPerGame   float64
    BlocksPerGame   float64
}

type PlayerCareerStats struct {
    PlayerID        int
    PlayerName      string
    GamesPlayed     int
    PointsPerGame   float64
    ReboundsPerGame float64
    AssistsPerGame  float64
    StealsPerGame   float64
    BlocksPerGame   float64
    FieldGoalPct    float64
    ThreePointPct   float64
    FreeThrowPct    float64
    MinutesPerGame  float64
    PointsTotal   float64
    ReboundsTotal float64
    AssistsTotal  float64
    StealsTotal   float64
    BlocksTotal   float64
    MinutesTotal  float64
}