package api

import (
	"fmt"
	"net/http"
	"github.com/gin-gonic/gin"
    "go.uber.org/zap"
	"github.com/nbaisland/nbaisland/internal/logger"
	"github.com/nbaisland/nbaisland/internal/service"
)


type HealthHandler struct {
	HealthService *service.HealthService
}

func (h *HealthHandler) CheckHealth(c *gin.Context){
	if err := h.HealthService.Check(c.Request.Context()); err != nil {
		logger.Log.Error("DB is not up",
			zap.Error(err),
			zap.String("handler", "CheckHealth"),
		)
		c.JSON(http.StatusInternalServerError, gin.H{"error" : fmt.Sprintf("Could not confirm db is up: %v", err)})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status" : "ok"})
}

func (h *HealthHandler) CheckReady(c *gin.Context){
	// ctx := c.Request.Context()
	c.JSON(http.StatusOK, gin.H{"status" : "ok"})
}
