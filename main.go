package main

import (
    "net/http"
    "log"
    "context"

    "github.com/gin-gonic/gin"
    "github.com/jackc/pgx/v5"
)

type User struct {
    Name  string `json:"name" binding:"required"`
    Email string `json:"email" binding:"required,email"`
    Currency int	`json:"currency"`
}

type Player struct {
    Name  string `json:"name" binding:"required"`
    Value int    `json:"value" binding:"required"`
}

type Holding struct {
    UserID   int `json:"user_id" binding:"required"`
    PlayerID int `json:"player_id" binding:"required"`
    Quantity int `json:"quantity" binding:"required"`
}

type PurchaseRequest struct {
    UserID int `json:"user_id" binding:"required"`
    PlayerID int `json:"player_id" binding:"required"`
}
const userStartingMoney int = 100

func main() {
    dsn := "postgres://devuser:devpass@localhost:5432/mydb?sslmode=disable"
    conn, err := pgx.Connect(context.Background(), dsn)
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close(context.Background())
    log.Println("Connected to the database successfully!")
    r := gin.Default()

    // Simple routes
    r.GET("/", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{"message": "Hello, Gin!"})
    })

    r.GET("/users", func(c *gin.Context) {
        rows, err := conn.Query(context.Background(), "SELECT id, name, email FROM users")
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        defer rows.Close()
        users := []map[string]interface{}{}
        for rows.Next() {
            var id int
            var name, email string
            rows.Scan(&id, &name, &email)
            users = append(users, map[string]interface{}{
                "id":    id,
                "name":  name,
                "email": email,
            })
        }
        c.JSON(http.StatusOK, users)
    })

    r.POST("/users", func(c *gin.Context) {
        var newUser User

        if err := c.ShouldBindJSON(&newUser); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
        newUser.Currency = userStartingMoney
        _, err := conn.Exec(c.Request.Context(),
            "INSERT INTO users (name, email, currency) VALUES ($1, $2, $3)",
            newUser.Name, newUser.Email, newUser.Currency,
        )
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        c.JSON(http.StatusCreated, gin.H{
            "message": "User created",
            "user":    newUser,
        })
    })

    r.GET("/holdings", func(c *gin.Context) {
        rows, err := conn.Query(context.Background(),
            "SELECT user_id, player_id, quantity FROM user_investments")
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        defer rows.Close()
        holdings := []map[string]interface{}{}
        for rows.Next() {
            var userID, playerID, quantity int
            rows.Scan(&userID, &playerID, &quantity)
            holdings = append(holdings, map[string]interface{}{
                "user_id":   userID,
                "player_id": playerID,
                "quantity":  quantity,
            })
        }
        c.JSON(http.StatusOK, holdings)
    })

    r.POST("/players/buy", func(c *gin.Context) {
        var newPurchase Holding
        if err := c.ShouldBindJSON(&newPurchase); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        // Calculate if there is "room" on the island

        // Caclulate if the user can purchase


    })

    r.POST("/players/buy", func(c *gin.Context) {
        var holdingToSell Holding
        if err := c.ShouldBindJSON(&newPurchase); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        // Calculate if there is "room" on the island

        // Caclulate if the user can purchase


    })

    r.GET("/players", func( c *gin.Context){
        rows.err := conn.Query(context.Background(),
        "SELECT id, name, value, capacity")
    })
    


    r.Run("0.0.0.0:8080")
}
