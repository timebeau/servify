package handlers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"servify/apps/server/internal/config"
	"servify/apps/server/internal/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// AuthHandler handles authentication endpoints.
type AuthHandler struct {
	DB     *gorm.DB
	Config *config.Config
}

// NewAuthHandler creates a new AuthHandler.
func NewAuthHandler(db *gorm.DB, cfg *config.Config) *AuthHandler {
	return &AuthHandler{DB: db, Config: cfg}
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

	// Check uniqueness
	var count int64
	h.DB.Model(&models.User{}).Where("username = ? OR email = ?", req.Username, req.Email).Count(&count)
	if count > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "用户名或邮箱已存在"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "内部错误"})
		return
	}

	role := req.Role
	if role == "" {
		role = "customer"
	}
	// Only allow admin role if explicitly requested and no users exist yet
	if role == "admin" {
		var total int64
		h.DB.Model(&models.User{}).Count(&total)
		if total > 0 {
			role = "customer" // Downgrade to customer for security
		}
	}

	user := models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hash),
		Name:     req.Name,
		Phone:    req.Phone,
		Role:     role,
		Status:   "active",
	}
	if err := h.DB.Create(&user).Error; err != nil {
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "unique") {
			c.JSON(http.StatusConflict, gin.H{"error": "用户名或邮箱已存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建用户失败"})
		return
	}

	token, err := h.generateToken(user.ID, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成 Token 失败"})
		return
	}

	now := time.Now()
	h.DB.Model(&user).Update("last_login", now)

	c.JSON(http.StatusCreated, tokenResponse{
		Token:     token,
		ExpiresIn: int(h.Config.JWT.ExpiresIn.Seconds()),
		User:      userResponse{ID: user.ID, Username: user.Username, Email: user.Email, Name: user.Name, Role: user.Role, Status: user.Status},
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

	var user models.User
	if err := h.DB.Where("username = ? OR email = ?", req.Username, req.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	if user.Status != "active" {
		c.JSON(http.StatusForbidden, gin.H{"error": "账号已被禁用"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	token, err := h.generateToken(user.ID, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成 Token 失败"})
		return
	}

	now := time.Now()
	h.DB.Model(&user).Update("last_login", now)

	c.JSON(http.StatusOK, tokenResponse{
		Token:     token,
		ExpiresIn: int(h.Config.JWT.ExpiresIn.Seconds()),
		User:      userResponse{ID: user.ID, Username: user.Username, Email: user.Email, Name: user.Name, Role: user.Role, Status: user.Status},
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
	userIDRaw, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}

	var userID uint
	switch v := userIDRaw.(type) {
	case float64:
		userID = uint(v)
	case uint:
		userID = v
	case int:
		userID = uint(v)
	default:
		c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的 Token"})
		return
	}

	var user models.User
	if err := h.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": userResponse{
			ID:     user.ID,
			Username: user.Username,
			Email:    user.Email,
			Name:     user.Name,
			Phone:    user.Phone,
			Avatar:   user.Avatar,
			Role:     user.Role,
			Status:   user.Status,
		},
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
	userIDRaw, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}

	var userID uint
	switch v := userIDRaw.(type) {
	case float64:
		userID = uint(v)
	case uint:
		userID = v
	case int:
		userID = uint(v)
	default:
		c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的 Token"})
		return
	}

	var user models.User
	if err := h.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	token, err := h.generateToken(user.ID, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成 Token 失败"})
		return
	}

	c.JSON(http.StatusOK, tokenResponse{
		Token:     token,
		ExpiresIn: int(h.Config.JWT.ExpiresIn.Seconds()),
		User:      userResponse{ID: user.ID, Username: user.Username, Email: user.Email, Name: user.Name, Role: user.Role, Status: user.Status},
	})
}

func (h *AuthHandler) generateToken(userID uint, role string) (string, error) {
	now := time.Now()
	payload := map[string]interface{}{
		"iat":    now.Unix(),
		"sub":    userID,
		"user_id": userID,
		"roles":  []string{role},
		"exp":    now.Add(h.Config.JWT.ExpiresIn).Unix(),
	}
	return createHS256JWT(payload, h.Config.JWT.Secret)
}

// createHS256JWT builds a compact JWT using HS256 with the given payload.
func createHS256JWT(payload map[string]interface{}, secret string) (string, error) {
	header := map[string]string{"alg": "HS256", "typ": "JWT"}
	headerJSON, _ := json.Marshal(header)
	payloadJSON, _ := json.Marshal(payload)
	enc := func(b []byte) string { return base64.RawURLEncoding.EncodeToString(b) }

	h := enc(headerJSON)
	p := enc(payloadJSON)
	signing := h + "." + p

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(signing))
	sig := mac.Sum(nil)
	s := enc(sig)
	return signing + "." + s, nil
}

// Request/Response types

type registerRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Role     string `json:"role"` // optional, defaults to "customer"
}

type loginRequest struct {
	Username string `json:"username"` // username or email
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
