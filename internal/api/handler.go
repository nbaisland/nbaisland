package api

import (
	"fmt"
	"strings"
	"strconv"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/nbaisland/nbaisland/internal/service"
)

type CreateUserRequest struct {
	Name    string `json:"name"`
	UserName    string `json:"userName"`
	Email    string `json:"email"`
	Password    string `json:"password"`
}

type CreatedUserResponse struct {
	ID    int64 `json:"id`
	Name    string `json:"name`
	UserName    string `json:"userName"`
	Email    string `json:"email`
}

type CreatePlayer struct {
	Name    string `json:"name"`
	Value   float64 `json:"value"`
	Capacity  int   `json:"capacity"`
	Slug    string  `json:"slug"`
}

type TransactionRequest struct {
	PlayerID    int64  `json:"player_id"`
	UserID    int64  `json:"user_id"`
	Quantity    float64  `json:"quantity"`
}

type Handler struct {
	UserService *service.UserService
	PlayerService *service.PlayerService
	TransactionService *service.TransactionService
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
	id, err := strconv.ParseInt(idStr, 10, 64)
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


func (h *Handler) GetUserByUserName(c *gin.Context) {
	ctx := c.Request.Context()
	userName := c.Param("username")
	if userName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"Error" : "Provide a username"})
		return
	}
	user, err := h.UserService.GetByUserName(ctx, userName)
	if err != nil {
		c.JSON(500, gin.H{"error" : fmt.Sprintf("failed to fetch user for username specified `%v`, %v", userName, err)})
		return
	}
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"Error" : "Could not find user"})
		return
	}
	c.JSON(200, user)
}


func (h* Handler) GetTransactionsOfUser(c *gin.Context){
	ctx := c.Request.Context()
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error" : "Provide a valid id"})
		return
	}
	transactions, err := h.TransactionService.GetByUserID(ctx, id)
	if err != nil {
		c.JSON(500, gin.H{"error" : fmt.Sprintf("failed to fetch transactions for user for id specified `%v`", id)})
		return
	}
	if transactions == nil {
		c.JSON(http.StatusNotFound, gin.H{"message" : "No Transactions found for user"})
		return
	}
	c.JSON(200, transactions)

}

func (h* Handler) CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Invalid request: %v", err),
		})
		return
	}
	user, err := h.UserService.CreateUser(c.Request.Context(), req.Name, req.UserName, req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error" : fmt.Sprintf("Failed to create user: %v", err),
		})
		return
	}
	res := CreatedUserResponse{
		ID: user.ID,
		UserName: user.UserName,
		Name: user.Name,
		Email: user.Email,
	}
	c.JSON(http.StatusOK, res)
}

func (h* Handler) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
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
		c.JSON(500, gin.H{"error" : fmt.Sprintf("failed to fetch players : %v", err)})
		return
	}
	c.JSON(200, players)
}

func (h* Handler) GetPlayersByID(c *gin.Context) {
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
			c.JSON(http.StatusBadRequest, gin.H{"error" : fmt.Sprintf("invalid id %v", p)})
			return
		}
		ids = append(ids, id)
	}

	players, err := h.PlayerService.GetPlayersByIDs(ctx, ids)

	if err != nil {
		c.JSON(500, gin.H{"error" : fmt.Sprintf("failed to fetch players : %v", err)})
		return
	}
	c.JSON(200, players)
}

func (h* Handler) GetPlayerByID(c *gin.Context) {
	ctx := c.Request.Context()
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error" : "Provide a valid id"})
		return
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

func (h* Handler) GetPlayerBySlug(c *gin.Context) {
	ctx := c.Request.Context()
	slug := c.Param("slug")
	if slug == "" {
		c.JSON(http.StatusBadRequest, gin.H{"Error" : "Provide a valid Slug"})
		return
	}
	player, err := h.PlayerService.GetPlayerBySlug(ctx, slug)
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

func (h* Handler) GetTransactionsOfPlayer(c *gin.Context){
	ctx := c.Request.Context()
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error" : "Provide a valid id"})
		return
	}
	transactions, err := h.TransactionService.GetByPlayerID(ctx, id)
	if err != nil {
		c.JSON(500, gin.H{"error" : fmt.Sprintf("failed to fetch transactions for player for id specified `%v`, %v", id, err)})
		return
	}
	if transactions == nil {
		c.JSON(http.StatusNotFound, gin.H{"message" : "No Transactions found for player"})
		return
	}
	c.JSON(200, transactions)
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
	err := h.PlayerService.CreatePlayer(ctx, req.Name, req.Value, req.Capacity, req.Slug)
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
	id, err := strconv.ParseInt(idStr, 10, 64)
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

func (h* Handler) GetTransactionByID(c *gin.Context) {
	ctx := c.Request.Context()
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error" : "Provide a valid id"})
		return
	}
	transaction, err := h.TransactionService.GetTransactionByID(ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error" : fmt.Sprintf("Could not find transaction: %v", err)}) 
		return
	}
	if transaction == nil {
		c.JSON(http.StatusNotFound, gin.H{"error" : "Could not find transaction"})
		return
	}

	c.JSON(200, transaction)
}

func (h* Handler) GetTransactions(c *gin.Context) {
	ctx := c.Request.Context()
	transactions, err := h.TransactionService.GetAll(ctx)

	if err != nil {
		c.JSON(500, gin.H{"error" : fmt.Sprintf("failed to fetch transactions: %v", err)})
		return
	}
	c.JSON(200, transactions)
}

func (h* Handler) BuyTransaction(c *gin.Context) {
	ctx := c.Request.Context()
	var req TransactionRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Invalid request: %v", err),
		})
		return
	}
	err := h.TransactionService.Buy(ctx, req.UserID, req.PlayerID, req.Quantity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error" : fmt.Sprintf("Could not make purchase: %v", err)})
		return
	}
	c.JSON(http.StatusOK, req)
}

func (h *Handler) SellTransaction(c *gin.Context) {
	ctx := c.Request.Context()
	var req TransactionRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Invalid request: %v", err),
		})
		return
	}
   proceeds, err := h.TransactionService.Sell(ctx, req.UserID, req.PlayerID, req.Quantity)
   if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error" : fmt.Sprintf("Could not process trade: %v", err)})
   }
   c.JSON(http.StatusOK, gin.H{"proceeds" : proceeds})
}

// func (h* Handler) PreviewTransaction(c *gin.Context) {
// 	// would be cool to do a transaction preview
// }

func (h* Handler) GetPositions(c *gin.Context) {
	ctx := c.Request.Context()
	positions, err := h.TransactionService.GetPositions(ctx)

	if err != nil {
		c.JSON(500, gin.H{"error" : fmt.Sprintf("failed to fetch positions: %v", err)})
		return
	}
	c.JSON(200, positions)
}

func (h* Handler) GetPositionsOfPlayer(c *gin.Context){
	ctx := c.Request.Context()
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error" : "Provide a valid id"})
		return
	}
	positions, err := h.TransactionService.GetPositionsByPlayerID(ctx, id)
	if err != nil {
		c.JSON(500, gin.H{"error" : fmt.Sprintf("failed to fetch positions for player for id specified `%v`, %v", id, err)})
		return
	}
	if positions == nil {
		c.JSON(http.StatusNotFound, gin.H{"message" : "No positions found for player"})
		return
	}
	c.JSON(200, positions)
}

func (h* Handler) GetPositionsOfUser(c *gin.Context){
	ctx := c.Request.Context()
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error" : "Provide a valid id"})
		return
	}
	positions, err := h.TransactionService.GetPositionsByUserID(ctx, id)
	if err != nil {
		c.JSON(500, gin.H{"error" : fmt.Sprintf("failed to fetch positions for user for id specified `%v`, err", id, err)})
		return
	}
	if positions == nil {
		c.JSON(http.StatusNotFound, gin.H{"message" : "No positions found for user"})
		return
	}
	c.JSON(200, positions)

}