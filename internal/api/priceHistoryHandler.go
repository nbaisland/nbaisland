package api

import (
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"github.com/nbaisland/nbaisland/internal/logger"
	"github.com/nbaisland/nbaisland/internal/service"
)


type PriceHistoryHandler struct {
	PriceHistoryService *service.PriceHistoryService
}

func (h *PriceHistoryHandler) GetPlayerPriceHistory(c *gin.Context) {
	// players/:id/price-history?range=7d
	idStr := c.Param("id")
	playerID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid player id"})
		return
	}

	timeRange := c.Query("range")
	if timeRange == "" {
		timeRange = "30d"
	}

	history, err := h.PriceHistoryService.GetPlayerPriceHistory(
		c.Request.Context(),
		playerID,
		timeRange,
	)
	if err != nil {
		logger.Log.Error("Failed to Fetch Price History",
			zap.Error(err),
			zap.String("handler", "GetPlayerPriceHistory"),
		)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch price history"})
		return
	}

	c.JSON(http.StatusOK, history)
}
