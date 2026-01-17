package api

import (
	"fmt"
	"strings"
	"strconv"
	"net/http"
    "go.uber.org/zap"

	"github.com/gin-gonic/gin"
	"github.com/nbaisland/nbaisland/internal/logger"
	"github.com/nbaisland/nbaisland/internal/service"
)

type CreatePlayer struct {
	Name    string `json:"name"`
	Value   float64 `json:"value"`
	Capacity  int   `json:"capacity"`
	Slug    string  `json:"slug"`
}

type PlayerHandler struct {
	PlayerService *service.PlayerService
}


func (h *PlayerHandler) GetPlayers(c *gin.Context) {
	ctx := c.Request.Context()
	players, err := h.PlayerService.GetAll(ctx)

	if err != nil {
		logger.Log.Error("Failed to Fetch players",
			zap.Error(err),
			zap.String("handler", "GetPlayers"),
		)
		c.JSON(500, gin.H{"error" : fmt.Sprintf("failed to fetch players : %v", err)})
		return
	}
	c.JSON(200, players)
}

func (h *PlayerHandler) GetPlayersByID(c *gin.Context) {
	ctx := c.Request.Context()
	idsParam := c.Query("ids")
	if idsParam == "" {
		h.GetPlayers(c)
		return
	}
	splitIDS := strings.Split(idsParam, ",")
	ids := make([]int64, 0, len(splitIDS))
	for _, p := range splitIDS {
		id, err := strconv.ParseInt(strings.TrimSpace(p), 10, 64)
		if err != nil {
			logger.Log.Warn("Could not parse ids",
				zap.Error(err),
				zap.String("handler", "GetPlayersByID"),
				zap.String("IDs", idsParam),
			)
			c.JSON(http.StatusBadRequest, gin.H{"error" : fmt.Sprintf("invalid id %v", p)})
			return
		}
		ids = append(ids, id)
	}

	players, err := h.PlayerService.GetPlayersByIDs(ctx, ids)

	if err != nil {
		logger.Log.Error("Failed to Fetch players by id",
			zap.Error(err),
			zap.String("handler", "GetPlayersByID"),
			zap.String("IDs", idsParam),
		)
		c.JSON(500, gin.H{"error" : fmt.Sprintf("failed to fetch players : %v", err)})
		return
	}
	c.JSON(200, players)
}

func (h *PlayerHandler) GetPlayerByID(c *gin.Context) {
	ctx := c.Request.Context()
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		logger.Log.Warn("Invalid ID provided",
			zap.Error(err),
			zap.String("handler", "GetPlayerByID"),
			zap.String("ID", idStr),
		)
		c.JSON(http.StatusBadRequest, gin.H{"Error" : "Provide a valid id"})
		return
	}
	player, err := h.PlayerService.GetPlayerByID(ctx, id)
	if err != nil {
		logger.Log.Error("Failed to Fetch player by id",
			zap.Error(err),
			zap.String("handler", "GetPlayer"),
			zap.Int64("IDs", id),
		)
		c.JSON(http.StatusInternalServerError, gin.H{"Error" : fmt.Sprintf("Could not find player: %v", err)}) 
		return
	}
	c.JSON(200, player)
}

func (h *PlayerHandler) GetPlayerBySlug(c *gin.Context) {
	ctx := c.Request.Context()
	slug := c.Param("slug")
	if slug == "" {
		logger.Log.Warn("Blank Slug Provided",
			zap.String("handler", "GetPlayerBySlug"),
			zap.String("slug", slug),
		)
		c.JSON(http.StatusBadRequest, gin.H{"Error" : "Provide a valid Slug"})
		return
	}
	player, err := h.PlayerService.GetPlayerBySlug(ctx, slug)
	if err != nil {
		logger.Log.Error("Failed to Fetch player by slug",
			zap.Error(err),
			zap.String("handler", "GetPlayerBySlug"),
			zap.String("Slug", slug),
		)
		c.JSON(http.StatusInternalServerError, gin.H{"Error" : fmt.Sprintf("Could not find player: %v", err)}) 
		return
	}
	c.JSON(200, player)
}

func (h *PlayerHandler) CreatePlayer(c *gin.Context){
	ctx := c.Request.Context()
	var req CreatePlayer
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Invalid request: %v", err),
		})
		return
	}
	err := h.PlayerService.CreatePlayer(ctx, req.Name, req.Value, req.Capacity, req.Slug)
	if err != nil {
		logger.Log.Error("Failed to Create player",
			zap.Error(err),
			zap.String("handler", "CreatePlayer"),
			zap.Any("Request", req),
		)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Could not create player: %v", err),
		})
		return
	}
	c.JSON(http.StatusOK, req)
}

func (h *PlayerHandler) DeletePlayer(c *gin.Context){
	ctx := c.Request.Context()
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		logger.Log.Warn("Invalid ID provided to delete player",
			zap.Error(err),
			zap.String("handler", "DeletePlayer"),
			zap.String("id", idStr),
		)
		c.JSON(http.StatusBadRequest, gin.H{
			"error" : "Provide a valid id",
		})
	}
	err = h.PlayerService.DeletePlayer(ctx, id)
	if err != nil {
		logger.Log.Error("Could Not delete Player",
			zap.Error(err),
			zap.String("handler", "DeletePlayer"),
			zap.Int64("id", id),
		)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error" : "Could Not delete player",
		})
		return
	}
	c.JSON(http.StatusOK, id)
}
