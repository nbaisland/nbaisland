package main

import (
    "context"
    "fmt"
    "log"
    "time"
    
    "github.com/nbaisland/nbaisland/internal/config"
    "github.com/nbaisland/nbaisland/internal/nba"
    "github.com/nbaisland/nbaisland/internal/repository"
    "github.com/nbaisland/nbaisland/internal/service"
)

func main() {
    log.Println("=== Player Value Update ===")
    
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
    
    playerRepo := &repository.PSQLPlayerRepo{Pool: pool}
    playerMapRepo := &repository.PlayerMapRepo{Pool: pool}
    nbaRepo := nba.NewRepository(pool)
    
    valueService := service.NewValueService(playerRepo, nbaRepo, playerMapRepo)
    
    ctx = context.Background()
    err = valueService.UpdateValueForAllPlayers(ctx)
    if err != nil {
        log.Fatalf("Error updating player values: %v", err)
    }
    
    log.Println("Successfully updated values for all players!")
}