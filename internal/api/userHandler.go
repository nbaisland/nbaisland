package api

import (
	"fmt"
	"strconv"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/nbaisland/nbaisland/internal/auth"
	"github.com/nbaisland/nbaisland/internal/service"
)


type UserHandler struct {
	UserService *service.UserService
}

func (h *UserHandler) GetUsers(c *gin.Context) {
	ctx := c.Request.Context()
	users, err := h.UserService.GetAll(ctx)
	if err != nil {
		c.JSON(500, gin.H{"error" : "failed to fetch users"})
		return
	}
	c.JSON(200, users)
}

func (h *UserHandler) GetUserByID(c *gin.Context) {
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


func (h *UserHandler) GetUserByUsername(c *gin.Context) {
	ctx := c.Request.Context()
	username := c.Param("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"Error" : "Provide a username"})
		return
	}
	user, err := h.UserService.GetByUsername(ctx, username)
	if err != nil {
		c.JSON(500, gin.H{"error" : fmt.Sprintf("failed to fetch user for username specified `%v`, %v", username, err)})
		return
	}
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"Error" : "Could not find user"})
		return
	}
	c.JSON(200, user)
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	claimsAny, ok := c.Get("user")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	claims := claimsAny.(*auth.Claims)

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Provide a valid id"})
		return
	}

	if claims.UserID != int64(id) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only delete your own account.. nice try"})
		return
	}

	err = h.UserService.DeleteUser(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Could not delete user: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"deleted_user_id": id})
}
