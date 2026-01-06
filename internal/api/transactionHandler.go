package api

import (
	"fmt"
	// "log"
	"strconv"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/nbaisland/nbaisland/internal/service"
)

type TransactionRequest struct {
	PlayerID    int64  `json:"player_id"`
	UserID    int64  `json:"user_id"`
	Quantity    int  `json:"quantity"`
}

type TransactionHandler struct {
	TransactionService *service.TransactionService
}

func (h *TransactionHandler) GetTransactionsOfUser(c *gin.Context){
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
		c.JSON(http.StatusOK, [])
		return
	}
	c.JSON(200, transactions)

}

func (h *TransactionHandler) GetTransactionsOfPlayer(c *gin.Context){
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
		c.JSON(http.StatusOK, [])
		return
	}
	c.JSON(200, transactions)
}

func (h *TransactionHandler) GetTransactionByID(c *gin.Context) {
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

func (h *TransactionHandler) GetTransactions(c *gin.Context) {
	ctx := c.Request.Context()
	transactions, err := h.TransactionService.GetAll(ctx)

	if err != nil {
		c.JSON(500, gin.H{"error" : fmt.Sprintf("failed to fetch transactions: %v", err)})
		return
	}
	c.JSON(200, transactions)
}

func (h *TransactionHandler) BuyTransaction(c *gin.Context) {
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

func (h *TransactionHandler) SellTransaction(c *gin.Context) {
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
		return
   }
   c.JSON(http.StatusOK, gin.H{"proceeds" : proceeds})
}

func (h *TransactionHandler) GetPositions(c *gin.Context) {
	ctx := c.Request.Context()
	positions, err := h.TransactionService.GetPositions(ctx)

	if err != nil {
		c.JSON(500, gin.H{"error" : fmt.Sprintf("failed to fetch positions: %v", err)})
		return
	}
	c.JSON(200, positions)
}

func (h *TransactionHandler) GetPositionsOfPlayer(c *gin.Context){
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
	// if positions == nil {
	// 	c.JSON(http.StatusNotFound, gin.H{"message" : "No positions found for player"})
	// 	return
	// }
	c.JSON(200, positions)
}

func (h *TransactionHandler) GetPositionsOfUser(c *gin.Context){
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
	// if positions == nil {
	// 	c.JSON(http.StatusNotFound, gin.H{"message" : "No positions found for user"})
	// 	return
	// }
	c.JSON(200, positions)

}