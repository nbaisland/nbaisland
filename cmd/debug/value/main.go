package main

import (
    "context"
    "fmt"
    "log"
    "time"
    
    "github.com/nbaisland/nbaisland/internal/config"
    "github.com/nbaisland/nbaisland/internal/logger"
    "github.com/nbaisland/nbaisland/internal/nba"
    "github.com/nbaisland/nbaisland/internal/repository"
    "github.com/nbaisland/nbaisland/internal/service"
    "go.uber.org/zap"
)

func main() {
    if err := logger.InitLogger("dev"); err != nil {
        log.Fatal("Failed to initialize logger:", err)
    }
    defer logger.Sync()
    
    cfg := config.Load()
    dsn := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=%v",
        cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName, cfg.DBSSLMODE)
    
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    pool, err := repository.NewDB(ctx, dsn)
    cancel()
    
    if err != nil {
        logger.Log.Fatal("Failed to connect to DB:", zap.Error(err))
    }
    defer pool.Close()
    
    logger.Log.Debug("Connected to database")
    
    playerRepo := &repository.PSQLPlayerRepo{Pool: pool}
    playerMapRepo := &repository.PlayerMapRepo{Pool: pool}
    nbaRepo := nba.NewRepository(pool)
    
    valueService := service.NewValueService(playerRepo, nbaRepo, playerMapRepo)
    
    ctx = context.Background()
    err = valueService.UpdateValueForAllPlayers(ctx, "2025-26")
    if err != nil {
        logger.Log.Fatal("Error updating player values:", zap.Error(err))
    }
    
    logger.Log.Info("Successfully updated values for all players!")
}