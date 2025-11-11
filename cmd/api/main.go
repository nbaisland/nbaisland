package main

import (
    "log"
    "context"
    "time"

    "github.com/gin-gonic/gin"

    "github.com/nthnklssn/sports_island/internal/service"
    "github.com/nthnklssn/sports_island/internal/repository"
    "github.com/nthnklssn/sports_island/internal/api"
)

func main() {
    ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
    pool, err := repository.NewDB(ctx)
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
