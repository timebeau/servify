package handlers

import (
	"errors"
	"net/http"
	"time"

	conversationdelivery "servify/apps/server/internal/modules/conversation/delivery"
	realtimeplatform "servify/apps/server/internal/platform/realtime"

	"github.com/gin-gonic/gin"
)

type ConversationWorkspaceHandler struct {
	service  conversationdelivery.HandlerService
	realtime realtimeplatform.RealtimeGateway
}

func NewConversationWorkspaceHandler(service conversationdelivery.HandlerService, realtime realtimeplatform.RealtimeGateway) *ConversationWorkspaceHandler {
	return &ConversationWorkspaceHandler{service: service, realtime: realtime}
}

func (h *ConversationWorkspaceHandler) ListMessages(c *gin.Context) {
	if h.service == nil {
		c.JSON(http.StatusServiceUnavailable, ErrorResponse{
			Error:   "Conversation service unavailable",
			Message: "conversation service is not configured",
		})
		return
	}

	sessionID := c.Param("id")
	items, err := h.service.ListMessages(c.Request.Context(), sessionID, 100)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, conversationdelivery.ErrConversationNotFound) {
			status = http.StatusNotFound
		}
		c.JSON(status, ErrorResponse{
			Error:   "Failed to load conversation messages",
			Message: err.Error(),
		})
		return
	}

	// Repository returns latest-first; management UI needs chronological order.
	for left, right := 0, len(items)-1; left < right; left, right = left+1, right-1 {
		items[left], items[right] = items[right], items[left]
	}

	c.JSON(http.StatusOK, gin.H{
		"data": items,
	})
}

func (h *ConversationWorkspaceHandler) SendMessage(c *gin.Context) {
	if h.service == nil {
		c.JSON(http.StatusServiceUnavailable, ErrorResponse{
			Error:   "Conversation service unavailable",
			Message: "conversation service is not configured",
		})
		return
	}

	sessionID := c.Param("id")
	var req struct {
		Content string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Message: err.Error(),
		})
		return
	}

	item, err := h.service.SendAgentMessage(c.Request.Context(), sessionID, req.Content)
	if err != nil {
		status := http.StatusInternalServerError
		errLabel := "Failed to send message"
		if errors.Is(err, conversationdelivery.ErrConversationNotFound) {
			status = http.StatusNotFound
			errLabel = "Conversation not found"
		}
		c.JSON(status, ErrorResponse{
			Error:   errLabel,
			Message: err.Error(),
		})
		return
	}

	if h.realtime != nil {
		h.realtime.SendToSession(sessionID, realtimeplatform.Message{
			Type: "agent-message",
			Data: map[string]interface{}{
				"content": item.Content,
				"sender":  item.Sender,
			},
			SessionID: sessionID,
			Timestamp: time.Now(),
		})
	}

	c.JSON(http.StatusCreated, SuccessResponse{
		Message: "Message sent successfully",
		Data:    item,
	})
}

func RegisterConversationWorkspaceRoutes(r *gin.RouterGroup, handler *ConversationWorkspaceHandler) {
	omni := r.Group("/omni")
	{
		omni.GET("/sessions/:id/messages", handler.ListMessages)
		omni.POST("/sessions/:id/messages", handler.SendMessage)
	}
}
