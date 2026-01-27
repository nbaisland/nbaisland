package main

import (
    "context"
    "fmt"
    "log"
    "os"
    "os/signal"
    "strings"
    "syscall"
    "time"

    "github.com/gin-contrib/cors"
    "github.com/gin-gonic/gin"
    "github.com/prometheus/client_golang/prometheus/promhttp"
    "go.uber.org/zap"

    "github.com/nbaisland/nbaisland/internal/api"
    "github.com/nbaisland/nbaisland/internal/config"
    "github.com/nbaisland/nbaisland/internal/database"
    "github.com/nbaisland/nbaisland/internal/logger"
    "github.com/nbaisland/nbaisland/internal/middleware"
    "github.com/nbaisland/nbaisland/internal/nba"
    "github.com/nbaisland/nbaisland/internal/repository"
    "github.com/nbaisland/nbaisland/internal/scheduler"
    "github.com/nbaisland/nbaisland/internal/service"
)

func main() {
    cfg := config.Load()
    if err := os.MkdirAll("logs", 0755); err != nil {
        log.Fatal("Failed to create logs directory, permissions?:", err)
    }

    if err := logger.InitLogger(cfg.ENV); err != nil {
        log.Fatal("Failed to initialize logger:", err)
    }
    defer logger.Sync()

    logger.Log.Info("Starting application",
        zap.String("env", cfg.ENV),
        zap.String("version", "0.1.0"),
    )

    if cfg.ENV == "production" {
	    gin.SetMode(gin.ReleaseMode)
    } else {
	    gin.SetMode(gin.DebugMode)
    }
    dsn := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=%v",
        cfg.DBUser,
        cfg.DBPassword,
        cfg.DBHost,
        cfg.DBPort,
        cfg.DBName,
        cfg.DBSSLMODE,
    )

    logger.Log.Info("Running DB Migrations")
    if err:= database.RunMigrations(dsn); err != nil{
        logger.Log.Fatal("Failed to run migrations",
            zap.Error(err),
        )
    }
    ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

    pool, err := repository.NewDB(ctx, dsn)
    if err != nil {
        logger.Log.Fatal("Failed to connect to database", zap.Error(err))
    }
    defer pool.Close()
    logger.Log.Info("Connected to the database successfully")

    nbaClient := nba.NewClient()
    nbaRepo := nba.NewRepository(pool)
    nbaService := nba.NewNBAService(nbaClient, nbaRepo, pool)

    userRepo := &repository.PSQLUserRepo{Pool: pool}
    UserService := service.NewUserService(userRepo)

    playerRepo := &repository.PSQLPlayerRepo{Pool: pool}
    PlayerService := service.NewPlayerService(playerRepo)

    transactionRepo := &repository.PSQLTransactionRepo{Pool: pool}
    TransactionService := service.NewTransactionService(transactionRepo, playerRepo, userRepo)

    priceHistoryRepo := &repository.PSQLPlayerPriceRepo{Pool: pool}
    PriceService := service.NewPriceHistoryService(priceHistoryRepo)

    playerMapRepo := &repository.PlayerMapRepo{Pool: pool}

    valueService := service.NewValueService(playerRepo, nbaRepo, playerMapRepo)
    HealthService := service.NewHealthService(pool)

    AuthHandler := &api.AuthHandler{UserService: UserService}
    userHandler := &api.UserHandler{UserService: UserService}
    playerHandler := &api.PlayerHandler{PlayerService: PlayerService}
    transactionHandler := &api.TransactionHandler{TransactionService: TransactionService}
    healthHandler := &api.HealthHandler{HealthService: HealthService}
    priceHistoryHandler := &api.PriceHistoryHandler{PriceHistoryService: PriceService}

    // #TODO: NBA Handler (admin only features).. scores etc

    sched := scheduler.New()

    sched.AddWeekly("Weekly Dividend", 4, 0, func(ctx context.Context) error {
        logger.Log.Info("Running scheduled weekly NBA stats update")
        return nbaService.UpdateAllWeeklyStats(ctx, "2025-26")
    })

    sched.AddNightly("Season Stats", 2, 0, func(ctx context.Context) error {
        logger.Log.Info("Running scheduled season NBA stats update")
        return nbaService.UpdateAllSeasonStats(ctx, "2025-26")
    })

    sched.AddNightly("Daily Update", 2, 40, func(ctx context.Context) error {
        logger.Log.Info("Daily Value Update")
        return valueService.UpdateValueForAllPlayers(ctx, "2025-26")
    })

    appCtx, appCancel := context.WithCancel(context.Background())
    defer appCancel()

    sched.Start(appCtx)

    r := gin.New()

    r.Use(gin.Recovery())
    r.Use(middleware.RequestIDMiddleware()) 
    r.Use(middleware.LoggingMiddleware())
    allowedOrigins := []string{"http://localhost:3000"}
    if cfg.CORSOrigin != "" {
        envOrigins := strings.Split(cfg.CORSOrigin, ",")
        for _, origin := range envOrigins {
            allowedOrigins = append(allowedOrigins, strings.TrimSpace(origin))
        }
    }
    logger.Log.Info("CORS configuration", zap.Strings("allowed_origins", allowedOrigins))
    r.Use(cors.New(cors.Config{
        AllowOrigins:     allowedOrigins,
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Content-Type", "Authorization"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
        MaxAge:           12 * time.Hour,
    }))

    r.GET("/health", healthHandler.CheckHealth)
    r.GET("/ready", healthHandler.CheckReady)
    r.GET("/metrics",  gin.WrapH(promhttp.Handler()))
    r.POST("/auth/register", AuthHandler.Register)
    r.POST("/auth/login", AuthHandler.Login)

    api := r.Group("/api")
    api.Use(middleware.AuthMiddleware())
    {
        api.GET("/users", userHandler.GetUsers)
        api.GET("/users/:id", userHandler.GetUserByID)
        api.GET("/users/username/:username", userHandler.GetUserByUsername)

        api.GET("/players", playerHandler.GetPlayersByID)
        api.GET("/players/:id", playerHandler.GetPlayerByID)
        api.GET("/players/name/:slug", playerHandler.GetPlayerBySlug)
        api.GET("/players/:id/price-history", priceHistoryHandler.GetPlayerPriceHistory)
        api.GET("/auth/me", AuthHandler.GetCurrentUser)
        // api.POST("/auth/logout", AuthHandler.Logout)

        api.DELETE("/users/:id", userHandler.DeleteUser)

        api.POST("/players", playerHandler.CreatePlayer)
        api.DELETE("/players/:id", playerHandler.DeletePlayer)

        api.GET("/transactions", transactionHandler.GetTransactions)
        api.POST("/transactions/buy", transactionHandler.BuyTransaction)
        api.POST("/transactions/sell", transactionHandler.SellTransaction)
        api.GET("/transactions/:id", transactionHandler.GetTransactionByID)

        api.GET("/positions", transactionHandler.GetPositions)

        api.GET("/users/:id/transactions", transactionHandler.GetTransactionsOfUser)
        api.GET("/users/:id/positions", transactionHandler.GetPositionsOfUser)
        api.GET("/players/:id/transactions", transactionHandler.GetTransactionsOfPlayer)
        api.GET("/players/:id/positions", transactionHandler.GetPositionsOfPlayer)
    }

    go func() {
        if err := r.Run(":8080"); err != nil {
            logger.Log.Fatal("Server failed to start", zap.Error(err))
        }
    }()
    logger.Log.Info("Server started", zap.Int("port", 8080))
    
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit

    logger.Log.Info("Shutting down server")

    appCancel()
}