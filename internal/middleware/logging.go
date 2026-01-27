package middleware

import (
	"fmt"
    "net/http"
    "time"
    "go.uber.org/zap"
    "github.com/gin-gonic/gin"
	"github.com/nbaisland/nbaisland/internal/logger"
	"github.com/nbaisland/nbaisland/metrics"
)

type responseWriter struct {
    http.ResponseWriter
    statusCode int
    bytesWritten int
}

func (rw *responseWriter) WriteHeader(code int) {
    rw.statusCode = code
    rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
    n, err := rw.ResponseWriter.Write(b)
    rw.bytesWritten += n
    return n, err
}

func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		var userID int64
		if userIDVal, exists := c.Get("user_id"); exists {
			if id, ok := userIDVal.(int64); ok {
				userID = id
			}
		}

		requestID := ""
		if reqID, exists := c.Get("request_id"); exists {
			if id, ok := reqID.(string); ok {
				requestID = id
			}
		}

		c.Next()

		duration := time.Since(start)

		logger.Log.Info("HTTP Request",
			zap.String("request_id", requestID),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("duration", duration),
			zap.Int64("user_id", userID),
			zap.String("ip", c.ClientIP()),
			zap.String("user_agent", c.Request.UserAgent()),
			zap.Int("size", c.Writer.Size()),
		)
		metrics.HttpRequestsTotal.WithLabelValues(
			c.Request.Method,
			path,
			fmt.Sprintf("%d",c.Writer.Status()),
		).Inc()

		metrics.HttpRequestDuration.WithLabelValues(
			c.Request.Method,
			path,
		).Observe(duration.Seconds())

		if path == "/auth/login" {
			return
		}
		if duration > time.Second {
			logger.Log.Warn("Slow request detected",
				zap.String("path", path),
				zap.Duration("duration", duration),
				zap.String("request_id", requestID),
			)
		}
	}
}
