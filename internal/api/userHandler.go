package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/nbaisland/nbaisland/internal/auth"
	"github.com/nbaisland/nbaisland/internal/logger"
	"github.com/nbaisland/nbaisland/internal/service"
)

type UserHandler struct {
	UserService *service.UserService
}

func (h *UserHandler) GetUsers(c *gin.Context) {
	ctx := c.Request.Context()

	users, err := h.UserService.GetAll(ctx)
	if err != nil {
		logger.Log.Error("failed to fetch users", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch users"})
		return
	}

	c.JSON(http.StatusOK, users)
}

func (h *UserHandler) GetUserByID(c *gin.Context) {
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

	user, err := h.UserService.GetByID(ctx, id)
	if err != nil {
		logger.Log.Error("failed to fetch user by id",
			zap.Int64("user_id", id),
			zap.Error(err),
		)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch user"})
		return
	}

	if user == nil {
		logger.Log.Info("user not found",
			zap.Int64("user_id", id),
		)
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) GetUserByUsername(c *gin.Context) {
	ctx := c.Request.Context()

	username := c.Param("username")
	if username == "" {
		logger.Log.Warn("missing username parameter",
			zap.String("route", c.FullPath()),
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Provide a username"})
		return
	}

	user, err := h.UserService.GetByUsername(ctx, username)
	if err != nil {
		logger.Log.Error("failed to fetch user by username",
			zap.String("username", username),
			zap.Error(err),
		)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch user"})
		return
	}

	if user == nil {
		logger.Log.Info("user not found by username",
			zap.String("username", username),
		)
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	claimsAny, ok := c.Get("user")
	if !ok {
		logger.Log.Warn("unauthorized delete user attempt: missing claims",
			zap.String("route", c.FullPath()),
		)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	claims := claimsAny.(*auth.Claims)

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		logger.Log.Warn("invalid user id parameter on delete",
			zap.String("param", idStr),
			zap.Int64("auth_user_id", claims.UserID),
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Provide a valid id"})
		return
	}

	if claims.UserID != id {
		logger.Log.Warn("forbidden delete user attempt",
			zap.Int64("auth_user_id", claims.UserID),
			zap.Int64("target_user_id", id),
		)
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only delete your own account"})
		return
	}

	if err := h.UserService.DeleteUser(c.Request.Context(), id); err != nil {
		logger.Log.Error("failed to delete user",
			zap.Int64("user_id", id),
			zap.Error(err),
		)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not delete user"})
		return
	}

	logger.Log.Info("user deleted",
		zap.Int64("user_id", id),
	)

	c.JSON(http.StatusOK, gin.H{"deleted_user_id": id})
}
