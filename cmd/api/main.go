package main

import (
    "log"
    "context"
    "time"
    "fmt"

    "github.com/gin-gonic/gin"
    "github.com/gin-contrib/cors"
    "github.com/nbaisland/nbaisland/internal/config"
    "github.com/nbaisland/nbaisland/internal/service"
    "github.com/nbaisland/nbaisland/internal/repository"
    "github.com/nbaisland/nbaisland/internal/api"
)

func main() {
    cfg := config.Load()
	dsn := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=%v",
        cfg.DBUser, 
        cfg.DBPassword,
        cfg.DBHost,
        cfg.DBPort,
        cfg.DBName,
        cfg.DBSSLMODE,
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
    HealthService := service.NewHealthService(pool)

    handler := &api.Handler{
        UserService: UserService,
        PlayerService: PlayerService,
        HoldingService: HoldingService,
        HealthService : HealthService,
    }


    r := gin.Default()

    r.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"http://127.0.0.1:5173"},
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Content-Type", "Authorization"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
        MaxAge:           12 * time.Hour,
    }))

    r.GET("/health", handler.CheckHealth)
    r.GET("/ready", handler.CheckReady)
    r.GET("/users", handler.GetUsers)
    r.GET("/users/:id", handler.GetUserByID)
    r.GET("/users/:id/holdings", handler.GetHoldingsOfUser)
    r.POST("/users", handler.CreateUser)
    r.DELETE("/users/:id", handler.DeleteUser)


    r.GET("/players", handler.GetPlayersByID)
    r.GET("/players/:id", handler.GetPlayerByID)
    r.GET("/players/:id/holdings", handler.GetHoldingsOfPlayer)
    r.POST("/players", handler.CreatePlayer)
    r.DELETE("/players/:id", handler.DeletePlayer)

    r.GET("/holdings", handler.GetHoldings)
    r.POST("/holdings", handler.MakePurchase)
    r.GET("/holdings/:id", handler.GetHoldingByID)
    r.POST("/holdings/:id/sell", handler.SellHolding)
    r.Run("0.0.0.0:8080")
}
