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