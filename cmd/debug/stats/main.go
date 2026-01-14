package main

import (
    "context"
    "fmt"
    "log"
    "time"
    
    "github.com/nbaisland/nbaisland/internal/config"
    "github.com/nbaisland/nbaisland/internal/nba"
    "github.com/nbaisland/nbaisland/internal/repository"
)

func main() {
    log.Println("=== Player Stats Update ===")
    
    cfg := config.Load()
    dsn := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=%v",
        cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName, cfg.DBSSLMODE)
    
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    pool, err := repository.NewDB(ctx, dsn)
    cancel()
    
    if err != nil {
        log.Fatalf("Failed to connect to DB: %v", err)
    }
    defer pool.Close()
    
    log.Println("Connected to database")
    
    nbaClient := nba.NewClient()
    nbaRepo := nba.NewRepository(pool)
    nbaService := nba.NewNBAService(nbaClient, nbaRepo, pool)

    ctx = context.Background()
    
    err = nbaService.UpdateAllSeasonStats(ctx, "2025-26")
    if err != nil {
        log.Fatalf("Error updating player season stats: %v", err)
    }
    
    log.Println("Successfully updated values for all players!")
}