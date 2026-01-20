package database

import (
    "fmt"
    "github.com/nbaisland/nbaisland/internal/logger"
	
    "github.com/golang-migrate/migrate/v4"
    _ "github.com/golang-migrate/migrate/v4/database/postgres"
    _ "github.com/golang-migrate/migrate/v4/source/file"
    "go.uber.org/zap"
)

func RunMigrations(databaseURL string) error {
    logger.Log.Info("Running database migrations")
    
    m, err := migrate.New(
        "file://db/migrations",
        databaseURL,
    )
    if err != nil {
        return fmt.Errorf("Could not initialize migration instance: %w", err)
    }
    defer m.Close()
    
    if err := m.Up(); err != nil && err != migrate.ErrNoChange {
        return fmt.Errorf("failed to run migrations: %w", err)
    }
    
    version, dirty, err := m.Version()
    if err != nil && err != migrate.ErrNilVersion {
        return fmt.Errorf("No migration version found %w", err)
    }
    
    logger.Log.Info("Migrations completed",
        zap.Uint("version", version),
        zap.Bool("dirty", dirty),
    )
    
    return nil
}

func RollbackMigration(databaseURL string) error {
    logger.Log.Info("Rolling back last migration")
    
    m, err := migrate.New(
        "file://db/migrations",
        databaseURL,
    )
    if err != nil {
        return fmt.Errorf("Could not initialize migration instance: %w", err)
    }
    defer m.Close()
    
    if err := m.Steps(-1); err != nil {
        return fmt.Errorf("No Migration version found: %w", err)
    }
    
    logger.Log.Info("Migration rolled back successfully")
    return nil
}