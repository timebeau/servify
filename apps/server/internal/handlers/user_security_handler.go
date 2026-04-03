package handlers

import (
	"net/http"
	"strconv"

	auditplatform "servify/apps/server/internal/platform/audit"
	"servify/apps/server/internal/platform/usersecurity"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type UserSecurityHandler struct {
	service *usersecurity.Service
	logger  *logrus.Logger
}

type batchRevokeTokensRequest struct {
	UserIDs []uint `json:"user_ids" binding:"required,min=1"`
}

type userSecurityPreviewItem struct {
	UserID           uint   `json:"user_id"`
	Username         string `json:"username"`
	Name             string `json:"name"`
	Role             string `json:"role"`
	Status           string `json:"status"`
	TokenVersion     int    `json:"token_version"`
	NextTokenVersion int    `json:"next_token_version"`
	LastLogin        any    `json:"last_login"`
	TokenValidAfter  any    `json:"token_valid_after"`
}

type revokeSessionRequest struct {
	SessionID string `json:"session_id" binding:"required"`
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
	auditplatform.SetBefore(c, gin.H{
		"user_id":           user.ID,
		"role":              user.Role,
		"status":            user.Status,
		"token_version":     user.TokenVersion,
		"token_valid_after": user.TokenValidAfter,
		"last_login":        user.LastLogin,
	})

	version, err := h.service.RevokeTokens(c.Request.Context(), uint(id))
	if err != nil {
		h.logger.Errorf("Failed to revoke user %d tokens: %v", id, err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to revoke user tokens",
			Message: err.Error(),
		})
		return
	}
	updated, err := h.service.GetUser(c.Request.Context(), uint(id))
	if err == nil && updated != nil {
		auditplatform.SetAfter(c, gin.H{
			"user_id":           updated.ID,
			"role":              updated.Role,
			"status":            updated.Status,
			"token_version":     updated.TokenVersion,
			"token_valid_after": updated.TokenValidAfter,
			"last_login":        updated.LastLogin,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "User tokens revoked successfully",
		"user_id":       id,
		"role":          user.Role,
		"token_version": version,
	})
}

// BatchRevokeTokens 批量强制失效用户已有 token
// @Summary 批量强制失效用户 token
// @Description 批量提升 token_version 并刷新 token_valid_after，使这些用户旧 token 失效
// @Tags 安全管理
// @Accept json
// @Produce json
// @Param payload body object true "用户 ID 列表"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/security/users/revoke-tokens [post]
func (h *UserSecurityHandler) BatchRevokeTokens(c *gin.Context) {
	if h.service == nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "User security service unavailable",
			Message: "user security service not configured",
		})
		return
	}

	var req batchRevokeTokensRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Message: err.Error(),
		})
		return
	}

	before := make([]gin.H, 0, len(req.UserIDs))
	for _, userID := range req.UserIDs {
		user, err := h.service.GetUser(c.Request.Context(), userID)
		if err != nil {
			h.logger.Errorf("Failed to load user %d: %v", userID, err)
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   "User not found",
				Message: err.Error(),
			})
			return
		}
		before = append(before, gin.H{
			"user_id":           user.ID,
			"role":              user.Role,
			"status":            user.Status,
			"token_version":     user.TokenVersion,
			"token_valid_after": user.TokenValidAfter,
			"last_login":        user.LastLogin,
		})
	}
	auditplatform.SetBefore(c, before)

	versions, err := h.service.BatchRevokeTokens(c.Request.Context(), req.UserIDs)
	if err != nil {
		h.logger.Errorf("Failed to batch revoke user tokens: %v", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to batch revoke user tokens",
			Message: err.Error(),
		})
		return
	}

	after := make([]gin.H, 0, len(req.UserIDs))
	results := make([]gin.H, 0, len(req.UserIDs))
	for _, userID := range req.UserIDs {
		user, err := h.service.GetUser(c.Request.Context(), userID)
		if err == nil && user != nil {
			after = append(after, gin.H{
				"user_id":           user.ID,
				"role":              user.Role,
				"status":            user.Status,
				"token_version":     user.TokenVersion,
				"token_valid_after": user.TokenValidAfter,
				"last_login":        user.LastLogin,
			})
		}
		results = append(results, gin.H{
			"user_id":       userID,
			"token_version": versions[userID],
		})
	}
	auditplatform.SetAfter(c, after)

	c.JSON(http.StatusOK, gin.H{
		"message": "User tokens revoked successfully",
		"count":   len(results),
		"items":   results,
	})
}

// QueryUsersSecurity 批量查询用户安全状态预览
// @Summary 批量查询用户安全状态预览
// @Description 返回多个用户当前安全状态，并给出下一次 revoke 后的 token_version 预览
// @Tags 安全管理
// @Accept json
// @Produce json
// @Param payload body object true "用户 ID 列表"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/security/users/query [post]
func (h *UserSecurityHandler) QueryUsersSecurity(c *gin.Context) {
	if h.service == nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "User security service unavailable",
			Message: "user security service not configured",
		})
		return
	}

	var req batchRevokeTokensRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Message: err.Error(),
		})
		return
	}

	users, err := h.service.GetUsers(c.Request.Context(), req.UserIDs)
	if err != nil {
		h.logger.Errorf("Failed to query user security state: %v", err)
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "User not found",
			Message: err.Error(),
		})
		return
	}

	items := make([]userSecurityPreviewItem, 0, len(users))
	for _, user := range users {
		items = append(items, userSecurityPreviewItem{
			UserID:           user.ID,
			Username:         user.Username,
			Name:             user.Name,
			Role:             user.Role,
			Status:           user.Status,
			TokenVersion:     user.TokenVersion,
			NextTokenVersion: user.TokenVersion + 1,
			LastLogin:        user.LastLogin,
			TokenValidAfter:  user.TokenValidAfter,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"count": len(items),
		"items": items,
	})
}

// ListUserSessions 获取用户 auth session 列表
// @Summary 获取用户 auth session 列表
// @Description 返回用户当前 auth session 状态，用于更细粒度 session 失效操作
// @Tags 安全管理
// @Produce json
// @Param id path int true "用户ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/security/users/{id}/sessions [get]
func (h *UserSecurityHandler) ListUserSessions(c *gin.Context) {
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

	if _, err := h.service.GetUser(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "User not found",
			Message: err.Error(),
		})
		return
	}

	sessions, err := h.service.ListUserSessions(c.Request.Context(), uint(id))
	if err != nil {
		h.logger.Errorf("Failed to list sessions for user %d: %v", id, err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to list user sessions",
			Message: err.Error(),
		})
		return
	}

	items := make([]gin.H, 0, len(sessions))
	for _, session := range sessions {
		items = append(items, gin.H{
			"session_id":        session.ID,
			"status":            session.Status,
			"token_version":     session.TokenVersion,
			"last_refreshed_at": session.LastRefreshedAt,
			"revoked_at":        session.RevokedAt,
			"created_at":        session.CreatedAt,
			"updated_at":        session.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id": uint(id),
		"count":   len(items),
		"items":   items,
	})
}

// RevokeSession 失效单个 auth session
// @Summary 失效单个 auth session
// @Description 吊销单个登录/刷新 session，使其后续 token 校验失败
// @Tags 安全管理
// @Accept json
// @Produce json
// @Param id path int true "用户ID"
// @Param payload body object true "session_id"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/security/users/{id}/sessions/revoke [post]
func (h *UserSecurityHandler) RevokeSession(c *gin.Context) {
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

	var req revokeSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Message: err.Error(),
		})
		return
	}

	sessions, err := h.service.ListUserSessions(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "User not found",
			Message: err.Error(),
		})
		return
	}
	for _, session := range sessions {
		if session.ID == req.SessionID {
			auditplatform.SetBefore(c, gin.H{
				"user_id":           id,
				"session_id":        session.ID,
				"status":            session.Status,
				"token_version":     session.TokenVersion,
				"last_refreshed_at": session.LastRefreshedAt,
				"revoked_at":        session.RevokedAt,
			})
			break
		}
	}

	session, err := h.service.RevokeSession(c.Request.Context(), uint(id), req.SessionID)
	if err != nil {
		h.logger.Errorf("Failed to revoke session %s for user %d: %v", req.SessionID, id, err)
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "Session not found",
			Message: err.Error(),
		})
		return
	}
	auditplatform.SetAfter(c, gin.H{
		"user_id":           id,
		"session_id":        session.ID,
		"status":            session.Status,
		"token_version":     session.TokenVersion,
		"last_refreshed_at": session.LastRefreshedAt,
		"revoked_at":        session.RevokedAt,
	})

	c.JSON(http.StatusOK, gin.H{
		"message":       "User session revoked successfully",
		"user_id":       uint(id),
		"session_id":    session.ID,
		"token_version": session.TokenVersion,
		"status":        session.Status,
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
		security.GET("/users/:id/sessions", handler.ListUserSessions)
		security.POST("/users/:id/sessions/revoke", handler.RevokeSession)
		security.POST("/users/query", handler.QueryUsersSecurity)
		security.POST("/users/revoke-tokens", handler.BatchRevokeTokens)
		security.POST("/users/:id/revoke-tokens", handler.RevokeTokens)
	}
}
