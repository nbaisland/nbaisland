package api

import (
	"strconv"
	"net/http"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/nbaisland/nbaisland/internal/logger"
	"github.com/nbaisland/nbaisland/internal/service"
)

type TransactionRequest struct {
	PlayerID int64 `json:"player_id"`
	UserID   int64 `json:"user_id"`
	Quantity int   `json:"quantity"`
}

type TransactionHandler struct {
	TransactionService *service.TransactionService
}

func (h *TransactionHandler) GetTransactionsOfUser(c *gin.Context) {
	ctx := c.Request.Context()

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		logger.Log.Warn("invalid user id parameter",
			zap.String("param", idStr),
			zap.String("route", c.FullPath()),
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Provide a valid id"})
		return
	}

	transactions, err := h.TransactionService.GetByUserID(ctx, id)
	if err != nil {
		logger.Log.Error("failed to fetch transactions for user",
			zap.Int64("user_id", id),
			zap.Error(err),
		)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch transactions"})
		return
	}

	if transactions == nil {
		logger.Log.Debug("no transactions found for user",
			zap.Int64("user_id", id),
		)
		c.JSON(http.StatusOK, []map[string]interface{}{})
		return
	}

	c.JSON(http.StatusOK, transactions)
}

func (h *TransactionHandler) GetTransactionsOfPlayer(c *gin.Context) {
	ctx := c.Request.Context()

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		logger.Log.Warn("invalid player id parameter",
			zap.String("param", idStr),
			zap.String("route", c.FullPath()),
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Provide a valid id"})
		return
	}

	transactions, err := h.TransactionService.GetByPlayerID(ctx, id)
	if err != nil {
		logger.Log.Error("failed to fetch transactions for player",
			zap.Int64("player_id", id),
			zap.Error(err),
		)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch transactions"})
		return
	}

	if transactions == nil {
		logger.Log.Debug("no transactions found for player",
			zap.Int64("player_id", id),
		)
		c.JSON(http.StatusOK, []map[string]interface{}{})
		return
	}

	c.JSON(http.StatusOK, transactions)
}

func (h *TransactionHandler) GetTransactionByID(c *gin.Context) {
	ctx := c.Request.Context()
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		logger.Log.Warn("invalid transaction id parameter",
			zap.String("param", idStr),
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Provide a valid id"})
		return
	}

	transaction, err := h.TransactionService.GetTransactionByID(ctx, id)
	if err != nil {
		logger.Log.Error("failed to fetch transaction by id",
			zap.Int64("transaction_id", id),
			zap.Error(err),
		)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch transaction"})
		return
	}

	if transaction == nil {
		logger.Log.Debug("transaction not found",
			zap.Int64("transaction_id", id),
		)
		c.JSON(http.StatusNotFound, gin.H{"error": "Could not find transaction"})
		return
	}

	c.JSON(http.StatusOK, transaction)
}

func (h *TransactionHandler) GetTransactions(c *gin.Context) {
	ctx := c.Request.Context()
	transactions, err := h.TransactionService.GetAll(ctx)
	if err != nil {
		logger.Log.Error("failed to fetch all transactions", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch transactions"})
		return
	}

	c.JSON(http.StatusOK, transactions)
}

func (h *TransactionHandler) BuyTransaction(c *gin.Context) {
	ctx := c.Request.Context()
	var req TransactionRequest
	if err := c.BindJSON(&req); err != nil {
		logger.Log.Warn("invalid buy transaction request body",
			zap.Error(err),
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if err := h.TransactionService.Buy(ctx, req.UserID, req.PlayerID, req.Quantity); err != nil {
		logger.Log.Error("failed to execute buy transaction",
			zap.Int64("user_id", req.UserID),
			zap.Int64("player_id", req.PlayerID),
			zap.Int("quantity", req.Quantity),
			zap.Error(err),
		)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not make purchase"})
		return
	}

	c.JSON(http.StatusOK, req)
}

func (h *TransactionHandler) SellTransaction(c *gin.Context) {
	ctx := c.Request.Context()
	var req TransactionRequest
	if err := c.BindJSON(&req); err != nil {
		logger.Log.Warn("invalid sell transaction request body",
			zap.Error(err),
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	proceeds, err := h.TransactionService.Sell(ctx, req.UserID, req.PlayerID, req.Quantity)
	if err != nil {
		logger.Log.Error("failed to execute sell transaction",
			zap.Int64("user_id", req.UserID),
			zap.Int64("player_id", req.PlayerID),
			zap.Int("quantity", req.Quantity),
			zap.Error(err),
		)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not process trade"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"proceeds": proceeds})
}

func (h *TransactionHandler) GetPositions(c *gin.Context) {
	ctx := c.Request.Context()
	positions, err := h.TransactionService.GetPositions(ctx)
	if err != nil {
		logger.Log.Error("failed to fetch positions", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch positions"})
		return
	}

	c.JSON(http.StatusOK, positions)
}

func (h *TransactionHandler) GetPositionsOfPlayer(c *gin.Context) {
	ctx := c.Request.Context()
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		logger.Log.Warn("invalid player id parameter",
			zap.String("param", idStr),
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Provide a valid id"})
		return
	}

	positions, err := h.TransactionService.GetPositionsByPlayerID(ctx, id)
	if err != nil {
		logger.Log.Error("failed to fetch positions for player",
			zap.Int64("player_id", id),
			zap.Error(err),
		)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch positions"})
		return
	}

	c.JSON(http.StatusOK, positions)
}

func (h *TransactionHandler) GetPositionsOfUser(c *gin.Context) {
	ctx := c.Request.Context()
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		logger.Log.Warn("invalid user id parameter",
			zap.String("param", idStr),
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Provide a valid id"})
		return
	}

	positions, err := h.TransactionService.GetPositionsByUserID(ctx, id)
	if err != nil {
		logger.Log.Error("failed to fetch positions for user",
			zap.Int64("user_id", id),
			zap.Error(err),
		)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch positions"})
		return
	}

	c.JSON(http.StatusOK, positions)
}
