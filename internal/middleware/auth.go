package middleware

import (
    "net/http"
    "strings"

    "github.com/gin-gonic/gin"
    "github.com/nbaisland/nbaisland/internal/auth"
    "github.com/nbaisland/nbaisland/internal/logger"
    "go.uber.org/zap"
)

func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            logger.Log.Warn("Missing authorization header",
                zap.String("path", c.Request.URL.Path),
                zap.String("ip", c.ClientIP()),
            )
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
            c.Abort()
            return
        }

        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || parts[0] != "Bearer" {
            logger.Log.Warn("Invalid authorization header format",
                zap.String("path", c.Request.URL.Path),
                zap.String("ip", c.ClientIP()),
            )
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
            c.Abort()
            return
        }

        tokenString := parts[1]

        claims, err := auth.ValidateToken(tokenString)
        if err != nil {
            logger.Log.Warn("Invalid or expired token",
                zap.Error(err),
                zap.String("path", c.Request.URL.Path),
                zap.String("ip", c.ClientIP()),
            )
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
            c.Abort()
            return
        }

        c.Set("user_id", claims.UserID)
        c.Set("username", claims.Username)
        // c.Set("is_admin", claims.IsAdmin)

        c.Next()
    }
}