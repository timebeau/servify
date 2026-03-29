package handlers

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"servify/apps/server/internal/models"
	"servify/apps/server/internal/services"

	"github.com/gin-gonic/gin"
)

// AuthHandler handles authentication endpoints.
type AuthHandler struct {
	service authService
}

type authService interface {
	Register(ctx context.Context, req services.RegisterInput) (*services.AuthResult, error)
	Login(ctx context.Context, req services.LoginInput) (*services.AuthResult, error)
	GetCurrentUser(ctx context.Context, userID uint) (*models.User, error)
	RefreshToken(ctx context.Context, userID uint) (*services.AuthResult, error)
}

// NewAuthHandler creates a new AuthHandler.
func NewAuthHandler(service authService) *AuthHandler {
	return &AuthHandler{service: service}
}

// Register godoc
// @Summary Register a new user
// @Tags auth
// @Accept json
// @Produce json
// @Param body body registerRequest true "Registration data"
// @Success 201 {object} tokenResponse
// @Failure 400 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Router /api/v1/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数无效", "message": err.Error()})
		return
	}

	req.Username = strings.TrimSpace(req.Username)
	req.Email = strings.TrimSpace(req.Email)
	if req.Username == "" || req.Email == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户名、邮箱和密码不能为空"})
		return
	}

	result, err := h.service.Register(c.Request.Context(), services.RegisterInput{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
		Name:     req.Name,
		Phone:    req.Phone,
		Role:     req.Role,
	})
	if err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidAuthInput):
			c.JSON(http.StatusBadRequest, gin.H{"error": "用户名、邮箱和密码不能为空"})
		case errors.Is(err, services.ErrAuthUserAlreadyExists):
			c.JSON(http.StatusConflict, gin.H{"error": "用户名或邮箱已存在"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "创建用户失败"})
		}
		return
	}

	c.JSON(http.StatusCreated, tokenResponse{
		Token:     result.Token,
		ExpiresIn: result.ExpiresIn,
		User:      mapUserResponse(result.User),
	})
}

// Login godoc
// @Summary Login
// @Tags auth
// @Accept json
// @Produce json
// @Param body body loginRequest true "Login credentials"
// @Success 200 {object} tokenResponse
// @Failure 401 {object} map[string]string
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数无效"})
		return
	}

	result, err := h.service.Login(c.Request.Context(), services.LoginInput{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		switch {
		case errors.Is(err, services.ErrAuthInvalidCredentials):
			c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		case errors.Is(err, services.ErrAuthUserDisabled):
			c.JSON(http.StatusForbidden, gin.H{"error": "账号已被禁用"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "登录失败"})
		}
		return
	}

	c.JSON(http.StatusOK, tokenResponse{
		Token:     result.Token,
		ExpiresIn: result.ExpiresIn,
		User:      mapUserResponse(result.User),
	})
}

// GetCurrentUser godoc
// @Summary Get current user info
// @Tags auth
// @Security BearerAuth
// @Produce json
// @Success 200 {object} userResponse
// @Failure 401 {object} map[string]string
// @Router /api/v1/auth/me [get]
func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	userID, ok := authUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的 Token"})
		return
	}

	user, err := h.service.GetCurrentUser(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": mapUserResponse(user),
	})
}

// RefreshToken godoc
// @Summary Refresh JWT token
// @Tags auth
// @Security BearerAuth
// @Produce json
// @Success 200 {object} tokenResponse
// @Router /api/v1/auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	userID, ok := authUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的 Token"})
		return
	}

	result, err := h.service.RefreshToken(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	c.JSON(http.StatusOK, tokenResponse{
		Token:     result.Token,
		ExpiresIn: result.ExpiresIn,
		User:      mapUserResponse(result.User),
	})
}

func authUserID(c *gin.Context) (uint, bool) {
	userIDRaw, exists := c.Get("user_id")
	if !exists {
		return 0, false
	}

	switch v := userIDRaw.(type) {
	case float64:
		return uint(v), true
	case uint:
		return v, true
	case int:
		return uint(v), true
	default:
		return 0, false
	}
}

// Request/Response types
type registerRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Role     string `json:"role"`
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type tokenResponse struct {
	Token     string       `json:"token"`
	ExpiresIn int          `json:"expires_in"`
	User      userResponse `json:"user"`
}

type userResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Avatar   string `json:"avatar"`
	Role     string `json:"role"`
	Status   string `json:"status"`
}

func mapUserResponse(user *models.User) userResponse {
	if user == nil {
		return userResponse{}
	}
	return userResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Name:     user.Name,
		Phone:    user.Phone,
		Avatar:   user.Avatar,
		Role:     user.Role,
		Status:   user.Status,
	}
}
