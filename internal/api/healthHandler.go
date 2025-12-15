package api

import (
	"fmt"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/nbaisland/nbaisland/internal/service"
)


type HealthHandler struct {
	HealthService *service.HealthService
}

func (h *HealthHandler) CheckHealth(c *gin.Context){
	if err := h.HealthService.Check(c.Request.Context()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error" : fmt.Sprintf("Could not confirm db is up: %v", err)})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status" : "ok"})
}

func (h *HealthHandler) CheckReady(c *gin.Context){
	// ctx := c.Request.Context()
	c.JSON(http.StatusOK, gin.H{"status" : "ok"})
}
