package main

import (
    "log"
    "context"
    "time"
    "fmt"

    "github.com/gin-gonic/gin"

    "github.com/nthnklssn/nba_island/internal/config"
    "github.com/nthnklssn/nba_island/internal/service"
    "github.com/nthnklssn/nba_island/internal/repository"
    "github.com/nthnklssn/nba_island/internal/api"
)

func main() {
    cfg := config.Load()
	dsn := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable",
        cfg.DBUser, 
        cfg.DBPassword,
        cfg.DBHost,
        cfg.DBPort,
        cfg.DBName,
    )
    ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

    pool, err := repository.NewDB(ctx, dsn)
    if err != nil {
        log.Fatal(err)
    }
    defer pool.Close()
    log.Println("Connected to the database successfully!")

    userRepo := &repository.PSQLUserRepo{Pool: pool}
    UserService := service.NewUserService(userRepo)
    playerRepo := &repository.PSQLPlayerRepo{Pool: pool}
    PlayerService := service.NewPlayerService(playerRepo)
    holdingRepo := &repository.PSQLHoldingRepo{Pool: pool}
    HoldingService := service.NewHoldingService(holdingRepo, playerRepo, userRepo)

    handler := &api.Handler{
        UserService: UserService,
        PlayerService: PlayerService,
        HoldingService: HoldingService,
    }


    r := gin.Default()
    r.GET("/health", handler.CheckHealth)
    r.GET("/ready", handler.CheckReady)
    r.GET("/users", handler.GetUsers)
    r.GET("/users/:id", handler.GetUserByID)
    r.POST("/users", handler.CreateUser)
    r.DELETE("/users/:id", handler.DeleteUser)


    r.GET("/players", handler.GetPlayers)
    r.GET("/players/:id", handler.GetPlayerByID)
    r.POST("/players", handler.CreatePlayer)
    r.DELETE("/players/:id", handler.DeletePlayer)

    r.GET("/holdings", handler.GetHoldings)
    r.POST("/holdings", handler.MakePurchase)
    r.GET("/holdings/:id", handler.GetHoldingByID)
    r.POST("/holdings/:id/sell", handler.SellHolding)
    r.Run("0.0.0.0:8080")
}
