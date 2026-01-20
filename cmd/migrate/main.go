package main

import (
    "fmt"
    "log"
    "os"
    
    "github.com/nbaisland/nbaisland/internal/config"
    "github.com/nbaisland/nbaisland/internal/logger"
    "github.com/nbaisland/nbaisland/internal/database"

    "go.uber.org/zap"
)

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: go run cmd/migrate/main.go [up|down|version]")
        os.Exit(1)
    }
    if err := logger.InitLogger("dev"); err != nil {
        log.Fatal("Failed to initialize logger:", err)
    }
    defer logger.Sync()
    cfg := config.Load()
    dsn := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=%v",
        cfg.DBUser,
        cfg.DBPassword,
        cfg.DBHost,
        cfg.DBPort,
        cfg.DBName,
        cfg.DBSSLMODE,
    )
	logger.Log.Debug(dsn)
    
    command := os.Args[1]
    
    switch command {
    case "up":
        if err := database.RunMigrations(dsn); err != nil {
            logger.Log.Fatal("Migration failed:", zap.Error(err))
        }
        fmt.Println("Migrations applied successfully")
        
    case "down":
        if err := database.RollbackMigration(dsn); err != nil {
            logger.Log.Fatal("Rollback failed:", zap.Error(err))
        }
        fmt.Println("Last migration rolled back successfully")
        
    case "version":
        fmt.Println("Check logs for version info")
        
    default:
        fmt.Println("Unknown command. Use: up, down, or version")
        os.Exit(1)
    }
}