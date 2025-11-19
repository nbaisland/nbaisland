package api

import (
	"fmt"
	"strconv"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/nbaisland/nbaisland/internal/service"
)

type CreateUserRequest struct {
	Name    string `json:"name"`
	Email    string `json:"email"`
	Password    string `json:"password"`
}

type CreatedUserResponse struct {
	ID    int `json:"id`
	Name    string `json:"name`
	Email    string `json:"email`
}

type CreatePlayer struct {
	Name    string `json:"name"`
	Value   float64 `json:"value"`
	Capacity  int   `json:"capacity"`
}

type CreateHoldingRequest struct {
	PlayerID    int  `json:"player_id"`
	UserID    int  `json:"user_id"`
	Quantity    float64  `json:"quantity"`
}

type SellHoldingRequest struct {
	Quantity    float64  `json:"quantity"`     
}
type Handler struct {
	UserService *service.UserService
	PlayerService *service.PlayerService
	HoldingService *service.HoldingService
	HealthService *service.HealthService
}

func (h *Handler) GetUsers(c *gin.Context) {
	ctx := c.Request.Context()
	users, err := h.UserService.GetAll(ctx)
	if err != nil {
		c.JSON(500, gin.H{"error" : "failed to fetch users"})
		return
	}
	c.JSON(200, users)
}

func (h *Handler) CheckHealth(c *gin.Context){
	if err := h.HealthService.Check(c.Request.Context()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error" : fmt.Sprintf("Could not confirm db is up: %v", err)})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status" : "ok"})
}

func (h *Handler) CheckReady(c *gin.Context){
	// ctx := c.Request.Context()
	c.JSON(http.StatusOK, gin.H{"status" : "ok"})
}

func (h *Handler) GetUserByID(c *gin.Context) {
	ctx := c.Request.Context()
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error" : "Provide a valid id"})
		return
	}
	user, err := h.UserService.GetByID(ctx, id)
	if err != nil {
		c.JSON(500, gin.H{"error" : fmt.Sprintf("failed to fetch user for id specified `%v`", id)})
		return
	}
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"Error" : "Could not find user"})
		return
	}
	c.JSON(200, user)
}

func (h* Handler) CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Invalid request: %v", err),
		})
		return
	}
	user, err := h.UserService.CreateUser(c.Request.Context(), req.Name, req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error" : fmt.Sprintf("Failed to create user: %v", err),
		})
		return
	}
	res := CreatedUserResponse{
		ID: user.ID,
		Name: user.Name,
		Email: user.Email,
	}
	c.JSON(http.StatusOK, res)
}

func (h* Handler) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error" : "Provide a valid id",
		})
	}
	err = h.UserService.DeleteUser(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error" : fmt.Sprintf("Could not delete user: %v", err),
		})
		return
	}
	c.JSON(http.StatusOK, id)
}

func (h* Handler) GetPlayers(c *gin.Context) {
	ctx := c.Request.Context()
	players, err := h.PlayerService.GetAll(ctx)

	if err != nil {
		c.JSON(500, gin.H{"error" : "failed to fetch players"})
		return
	}
	c.JSON(200, players)
}

func (h* Handler) GetPlayerByID(c *gin.Context) {
	ctx := c.Request.Context()
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error" : "Provide a valid id"})
	}
	player, err := h.PlayerService.GetPlayerByID(ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error" : fmt.Sprintf("Could not find player: %v", err)}) 
		return
	}
	if player == nil {
		c.JSON(http.StatusNotFound, gin.H{"Error" : "Could not find player"})
		return
	}
	c.JSON(200, player)
}

func (h* Handler) CreatePlayer(c *gin.Context){
	ctx := c.Request.Context()
	var req CreatePlayer
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Invalid request: %v", err),
		})
		return
	}
	err := h.PlayerService.CreatePlayer(ctx, req.Name, req.Value, req.Capacity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Could not create player: %v", err),
		})
		return
	}
	c.JSON(http.StatusOK, req)
}

func (h* Handler) DeletePlayer(c *gin.Context){
	ctx := c.Request.Context()
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error" : "Provide a valid id",
		})
	}
	err = h.PlayerService.DeletePlayer(ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error" : "Could Not delete player",
		})
		return
	}
	c.JSON(http.StatusOK, id)
}

func (h* Handler) GetHoldingByID(c *gin.Context) {
	ctx := c.Request.Context()
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error" : "Provide a valid id"})
		return
	}
	holding, err := h.HoldingService.GetHoldingByID(ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error" : fmt.Sprintf("Could not find holding: %v", err)}) 
		return
	}
	if holding == nil {
		c.JSON(http.StatusNotFound, gin.H{"error" : "Could not find holding"})
		return
	}

	c.JSON(200, holding)
}

func (h* Handler) GetHoldings(c *gin.Context) {
	ctx := c.Request.Context()
	holdings, err := h.HoldingService.GetAll(ctx)

	if err != nil {
		c.JSON(500, gin.H{"error" : fmt.Sprintf("failed to fetch holdings: %v", err)})
		return
	}
	c.JSON(200, holdings)
}

func (h* Handler) MakePurchase(c *gin.Context) {
	ctx := c.Request.Context()
	var req CreateHoldingRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Invalid request: %v", err),
		})
		return
	}
	err := h.HoldingService.MakePurchase(ctx, req.PlayerID, req.UserID, req.Quantity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error" : fmt.Sprintf("Could not make purchase: %v", err)})
		return
	}
	c.JSON(http.StatusOK, req)
}

func (h *Handler) SellHolding(c *gin.Context) {
	ctx := c.Request.Context()
	idStr := c.Param("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid holding id"})
        return
    }
	var req SellHoldingRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Invalid request: %v", err),
		})
		return
	}
   proceeds, err := h.HoldingService.SellHolding(ctx, id, req.Quantity)
   if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error" : fmt.Sprintf("Could not process trade: %v", err)})
   }
   c.JSON(http.StatusOK, gin.H{"proceeds" : proceeds})
}