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
    log.Println("NBA Initial Setup")
    
    cfg := config.Load()
    dsn := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=%v",
        cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName, cfg.DBSSLMODE)
    
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    pool, err := repository.NewDB(ctx, dsn)
    if err != nil {
        log.Fatal(err)
    }
    defer pool.Close()
    
    nbaClient := nba.NewClient()
    nbaRepo := nba.NewRepository(pool)
    nbaService := nba.NewNBAService(nbaClient, nbaRepo, pool)

    ctx = context.Background()
    
    log.Println("Seeding players with 10+ games...")
    if err := nbaService.SeedTopPlayers(ctx, "2025-26", 10); err != nil {
        log.Fatalf("Seed failed: %v", err)
    }
    
    log.Println("Loading initial season stats...")
    if err := nbaService.UpdateAllSeasonStats(ctx, "2025-26"); err != nil {
        log.Fatalf("Season stats failed: %v", err)
    }
    log.Println("Loading career stats...")
    if err := nbaService.UpdateAllCareerStats(ctx); err != nil {
        log.Printf("Warning: Career stats failed: %v", err)
    }

    log.Println("Setup complete!")
}