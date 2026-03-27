package handlers

import (
	"net/http"
	"strconv"

	"servify/apps/server/internal/platform/usersecurity"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type UserSecurityHandler struct {
	service *usersecurity.Service
	logger  *logrus.Logger
}

func NewUserSecurityHandler(service *usersecurity.Service, logger *logrus.Logger) *UserSecurityHandler {
	if logger == nil {
		logger = logrus.StandardLogger()
	}
	return &UserSecurityHandler{service: service, logger: logger}
}

// RevokeTokens 强制失效用户已有 token
// @Summary 强制失效用户 token
// @Description 提升 token_version 并刷新 token_valid_after，使该用户旧 token 失效
// @Tags 安全管理
// @Accept json
// @Produce json
// @Param id path int true "用户ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/security/users/{id}/revoke-tokens [post]
func (h *UserSecurityHandler) RevokeTokens(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid user ID",
			Message: "ID must be a valid number",
		})
		return
	}

	if h.service == nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "User security service unavailable",
			Message: "user security service not configured",
		})
		return
	}

	user, err := h.service.GetUser(c.Request.Context(), uint(id))
	if err != nil {
		h.logger.Errorf("Failed to load user %d: %v", id, err)
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "User not found",
			Message: err.Error(),
		})
		return
	}

	version, err := h.service.RevokeTokens(c.Request.Context(), uint(id))
	if err != nil {
		h.logger.Errorf("Failed to revoke user %d tokens: %v", id, err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to revoke user tokens",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "User tokens revoked successfully",
		"user_id":       id,
		"role":          user.Role,
		"token_version": version,
	})
}

// GetUserSecurity 获取用户安全状态
// @Summary 获取用户安全状态
// @Description 返回用户当前状态、token_version、token_valid_after 和最近登录时间
// @Tags 安全管理
// @Produce json
// @Param id path int true "用户ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/security/users/{id} [get]
func (h *UserSecurityHandler) GetUserSecurity(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid user ID",
			Message: "ID must be a valid number",
		})
		return
	}

	if h.service == nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "User security service unavailable",
			Message: "user security service not configured",
		})
		return
	}

	user, err := h.service.GetUser(c.Request.Context(), uint(id))
	if err != nil {
		h.logger.Errorf("Failed to load user %d security state: %v", id, err)
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "User not found",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id":           user.ID,
		"role":              user.Role,
		"status":            user.Status,
		"token_version":     user.TokenVersion,
		"token_valid_after": user.TokenValidAfter,
		"last_login":        user.LastLogin,
	})
}

func RegisterUserSecurityRoutes(r *gin.RouterGroup, handler *UserSecurityHandler) {
	if handler == nil {
		return
	}
	security := r.Group("/security")
	{
		security.GET("/users/:id", handler.GetUserSecurity)
		security.POST("/users/:id/revoke-tokens", handler.RevokeTokens)
	}
}
