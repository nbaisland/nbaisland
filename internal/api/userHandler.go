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


func (h *UserHandler) GetUserByUserName(c *gin.Context) {
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

func (h *UserHandler) CreateUser(c *gin.Context) {
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

func (h *UserHandler) DeleteUser(c *gin.Context) {
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
