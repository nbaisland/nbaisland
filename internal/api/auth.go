package api

import (
	"net/http"
	"github.com/gin-gonic/gin"
    "go.uber.org/zap"

	"github.com/nbaisland/nbaisland/internal/logger"
	"github.com/nbaisland/nbaisland/internal/service"
	"github.com/nbaisland/nbaisland/internal/auth"
	"github.com/nbaisland/nbaisland/internal/models"
)

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token    string `json:"token"`
	UserID   int64    `json:"user_id"`
	Username string `json:"username"`
}

type AuthHandler struct {
	UserService *service.UserService
}
func (h *AuthHandler) Register(c *gin.Context) {
	ctx := c.Request.Context()
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Log.Error("Failed to Register User",
			zap.Error(err),
			zap.String("handler", "Register"),
			zap.Any("Request", req),
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if req.Username == "" {
		logger.Log.Debug("Blank Username",
			zap.String("handler", "Register"),
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username required"})
		return
	}
	if req.Password == "" {
		logger.Log.Debug("Blank Password",
			zap.String("handler", "Register"),
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password required"})
		return
	}
	if req.Email == "" {
		logger.Log.Debug("Blank Email",
			zap.String("handler", "Register"),
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email required"})
		return
	}

	existingUser, _ :=  h.UserService.GetByUsername(ctx, req.Username)
	if existingUser != nil {
		logger.Log.Debug("Blank Email",
			zap.String("handler", "Register"),
			zap.String("Email", req.Email),
		)
		c.JSON(http.StatusConflict, gin.H{"error": "Username already taken"})
		return
	}

	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		logger.Log.Debug("Problem hasging password",
			zap.String("handler", "Register"),
			zap.Error(err),
		)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating user, password issue"})
		return
	}

	user := &models.User{
		Username: req.Username,
		Password: hashedPassword,
		Email:    req.Email,
		Currency: 10000.0,
	}

	if err := h.UserService.CreateUser(ctx, user); err != nil {
		logger.Log.Error("Failed to Create User",
			zap.Error(err),
			zap.String("handler", "Register"),
			zap.Any("User", user),
		)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating user"})
		return
	}

	token, err := auth.GenerateToken(user.ID, user.Username)
	if err != nil {
		logger.Log.Error("Could not create a token",
			zap.Error(err),
			zap.String("handler", "Register"),
		)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
		return
	}

	c.JSON(http.StatusOK, AuthResponse{
		Token:    token,
		UserID:   user.ID,
		Username: user.Username,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	ctx := c.Request.Context()
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Log.Error("Failed to Parse Login message",
			zap.Error(err),
			zap.String("handler", "Register"),
			zap.Any("Request", req),
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	user, err := h.UserService.GetByUsername(ctx, req.Username)
	if err != nil {
		logger.Log.Debug("Could not get username",
			zap.String("handler", "Login"),
			zap.Error(err),
		)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}
	if user == nil {
		logger.Log.Debug("User doesnt exist",
			zap.String("handler", "Login"),
			zap.String("Username", req.Username),
		)
		c.JSON(http.StatusNotFound, gin.H{"error" : "User Does Not exist"})
		return
	}
	
	if !auth.CheckPassword(user.Password, req.Password) {
		logger.Log.Debug("Invalid Password ",
			zap.String("handler", "Login"),
		)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token, err := auth.GenerateToken(user.ID, user.Username)
	if err != nil {
		logger.Log.Error("Could not create a token",
			zap.Error(err),
			zap.String("handler", "Login"),
		)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
		return
	}
	c.JSON(http.StatusOK, AuthResponse{
		Token:    token,
		UserID:   user.ID,
		Username: user.Username,
	})
}

func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	ctx := c.Request.Context()
	claims, ok := c.Get("user")
	if !ok {
		logger.Log.Error("Could not get Current user",
			zap.String("handler", "GetCurrentUser"),
			zap.Any("Claims", claims),
		)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	authClaims, ok := claims.(*auth.Claims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	user, err := h.UserService.GetByID(ctx, authClaims.UserID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}
