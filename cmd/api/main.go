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
    "github.com/nbaisland/nbaisland/internal/nba"
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

    nbaClient := nba.NewClient()
    nbaRepo := nba.NewRepository(pool)
    nbaService := nba.NewService(nbaClient, nbaRepo, pool)

    userRepo := &repository.PSQLUserRepo{Pool: pool}
    UserService := service.NewUserService(userRepo)
    playerRepo := &repository.PSQLPlayerRepo{Pool: pool}
    PlayerService := service.NewPlayerService(playerRepo)
    transactionRepo := &repository.PSQLTransactionRepo{Pool: pool}
    TransactionService := service.NewTransactionService(transactionRepo, playerRepo, userRepo)
    HealthService := service.NewHealthService(pool)

    userHandler := &api.Handler{
        UserService: UserService,
    }

    playerHandler := &api.Handler{
        PlayerService: PlayerService,
    }

    transactionHandler := &api.Handler{
        TransactionService: TransactionService,
    }

    healthHandler := &api.Handler{
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

    r.GET("/health", healthHandler.CheckHealth)
    r.GET("/ready", healthHandler.CheckReady)

    r.GET("/users", userHandler.GetUsers)
    r.GET("/users/:id", userHandler.GetUserByID)
    r.GET("/users/username/:username", userHandler.GetUserByUserName)
    r.POST("/users", userHandler.CreateUser)
    r.DELETE("/users/:id", userHandler.DeleteUser)



    r.GET("/players", playerHandler.GetPlayersByID)
    r.GET("/players/:id", playerHandler.GetPlayerByID)
    r.GET("/players/name/:slug", playerHandler.GetPlayerBySlug)
    r.POST("/players", playerHandler.CreatePlayer)
    r.DELETE("/players/:id", playerHandler.DeletePlayer)

    r.GET("/transactions", transactionHandler.GetTransactions)
    r.POST("/transactions/buy", transactionHandler.BuyTransaction)
    r.GET("/transactions/:id", transactionHandler.GetTransactionByID)
    r.POST("/transactions/sell", transactionHandler.SellTransaction)
    r.GET("/positions", transactionHandler.GetPositions)
    r.GET("/users/:id/transactions", transactionHandler.GetTransactionsOfUser)
    r.GET("/users/:id/positions", transactionHandler.GetPositionsOfUser)
    r.GET("/players/:id/transactions", transactionHandler.GetTransactionsOfPlayer)
    r.GET("/players/:id/positions", transactionHandler.GetPositionsOfPlayer)
    r.Run("0.0.0.0:8080")
}
