package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	voicedelivery "servify/apps/server/internal/modules/voice/delivery"
	"servify/apps/server/internal/platform/voiceprotocol"

	"github.com/gin-gonic/gin"
)

type VoiceHandler struct {
	coordinator *voicedelivery.Coordinator
	registry    *voiceprotocol.Registry
}

func NewVoiceHandler(coordinator *voicedelivery.Coordinator, registry *voiceprotocol.Registry) *VoiceHandler {
	return &VoiceHandler{coordinator: coordinator, registry: registry}
}

type startRecordingRequest struct {
	CallID   string `json:"call_id" binding:"required"`
	Provider string `json:"provider"`
}

type stopRecordingRequest struct {
	RecordingID string `json:"recording_id" binding:"required"`
}

type appendTranscriptRequest struct {
	CallID    string `json:"call_id" binding:"required"`
	Content   string `json:"content" binding:"required"`
	Language  string `json:"language"`
	Finalized bool   `json:"finalized"`
}

type protocolEventRequest struct {
	Payload map[string]interface{} `json:"payload" binding:"required"`
}

func (h *VoiceHandler) StartRecording(c *gin.Context) {
	if h.coordinator == nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "voice coordinator unavailable"})
		return
	}
	var req startRecordingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Invalid request format: " + err.Error()})
		return
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), 15*time.Second)
	defer cancel()
	recording, err := h.coordinator.StartRecording(ctx, voicedelivery.StartRecordingCommand{
		CallID:   req.CallID,
		Provider: req.Provider,
	})
	if err != nil {
		respondVoiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": recording})
}

func (h *VoiceHandler) StopRecording(c *gin.Context) {
	if h.coordinator == nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "voice coordinator unavailable"})
		return
	}
	var req stopRecordingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Invalid request format: " + err.Error()})
		return
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), 15*time.Second)
	defer cancel()
	if err := h.coordinator.StopRecording(ctx, voicedelivery.StopRecordingCommand{RecordingID: req.RecordingID}); err != nil {
		respondVoiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *VoiceHandler) GetRecording(c *gin.Context) {
	if h.coordinator == nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "voice coordinator unavailable"})
		return
	}
	recordingID := c.Param("recordingID")
	if recordingID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "recordingID is required"})
		return
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	recording, err := h.coordinator.GetRecording(ctx, recordingID)
	if err != nil {
		respondVoiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": recording})
}

func (h *VoiceHandler) AppendTranscript(c *gin.Context) {
	if h.coordinator == nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "voice coordinator unavailable"})
		return
	}
	var req appendTranscriptRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Invalid request format: " + err.Error()})
		return
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), 15*time.Second)
	defer cancel()
	transcript, err := h.coordinator.AppendTranscript(ctx, voicedelivery.AppendTranscriptCommand{
		CallID:    req.CallID,
		Content:   req.Content,
		Language:  req.Language,
		Finalized: req.Finalized,
	})
	if err != nil {
		respondVoiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": transcript})
}

func (h *VoiceHandler) ListTranscripts(c *gin.Context) {
	if h.coordinator == nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "voice coordinator unavailable"})
		return
	}
	callID := c.Query("call_id")
	if callID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "call_id is required"})
		return
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	items, err := h.coordinator.ListTranscripts(ctx, callID)
	if err != nil {
		respondVoiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": items})
}

func respondVoiceError(c *gin.Context, err error) {
	status := http.StatusInternalServerError
	// Check for provider disabled/unavailable errors by message
	// This avoids importing the application layer's ProviderError type
	msg := err.Error()
	if msg == "voice recording provider is disabled" ||
		msg == "voice transcript provider is disabled" ||
		msg == "voice call provider is disabled" {
		status = http.StatusServiceUnavailable
	}
	c.JSON(status, ErrorResponse{
		Error:   "Voice runtime unavailable",
		Message: err.Error(),
		Code:    status,
	})
}

func (h *VoiceHandler) ListProtocols(c *gin.Context) {
	if h.registry == nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "voice protocol registry unavailable"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    h.registry.SupportedProtocols(),
	})
}

func (h *VoiceHandler) HandleProtocolCallEvent(c *gin.Context) {
	if h.registry == nil || h.coordinator == nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "voice protocol runtime unavailable"})
		return
	}
	adapter, ok := h.registry.Signaling(voiceprotocol.Protocol(c.Param("protocol")))
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "unsupported signaling protocol"})
		return
	}
	var req protocolEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Invalid request format: " + err.Error()})
		return
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), 15*time.Second)
	defer cancel()
	event, err := mapCallEvent(ctx, adapter, c.Param("event"), req.Payload)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}
	if err := h.coordinator.HandleCallEvent(ctx, event); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": event})
}

func (h *VoiceHandler) HandleProtocolMediaEvent(c *gin.Context) {
	if h.registry == nil || h.coordinator == nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "voice protocol runtime unavailable"})
		return
	}
	adapter, ok := h.registry.Media(voiceprotocol.Protocol(c.Param("protocol")))
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "unsupported media protocol"})
		return
	}
	var req protocolEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Invalid request format: " + err.Error()})
		return
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), 15*time.Second)
	defer cancel()
	event, err := mapMediaEvent(ctx, adapter, c.Param("event"), req.Payload)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}
	if err := h.coordinator.HandleMediaEvent(ctx, event); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": event})
}

func RegisterVoiceRoutes(r *gin.RouterGroup, handler *VoiceHandler) {
	if r == nil || handler == nil {
		return
	}
	r.GET("/voice/protocols", handler.ListProtocols)
	r.POST("/voice/protocols/:protocol/call-events/:event", handler.HandleProtocolCallEvent)
	r.POST("/voice/protocols/:protocol/media-events/:event", handler.HandleProtocolMediaEvent)
	r.POST("/voice/recordings/start", handler.StartRecording)
	r.POST("/voice/recordings/stop", handler.StopRecording)
	r.GET("/voice/recordings/:recordingID", handler.GetRecording)
	r.POST("/voice/transcripts", handler.AppendTranscript)
	r.GET("/voice/transcripts", handler.ListTranscripts)
}

func mapCallEvent(ctx context.Context, adapter voiceprotocol.CallSignalingAdapter, name string, payload map[string]interface{}) (voiceprotocol.CallEvent, error) {
	switch voiceprotocol.CallEventKind(name) {
	case voiceprotocol.CallEventInvite:
		return adapter.MapInvite(ctx, payload)
	case voiceprotocol.CallEventAnswer:
		return adapter.MapAnswer(ctx, payload)
	case voiceprotocol.CallEventHold:
		return adapter.MapHold(ctx, payload)
	case voiceprotocol.CallEventResume:
		return adapter.MapResume(ctx, payload)
	case voiceprotocol.CallEventHangup:
		return adapter.MapHangup(ctx, payload)
	case voiceprotocol.CallEventTransfer:
		return adapter.MapTransfer(ctx, payload)
	case voiceprotocol.CallEventDTMF:
		return adapter.MapDTMF(ctx, payload)
	default:
		return voiceprotocol.CallEvent{}, fmt.Errorf("unsupported call event %q", name)
	}
}

func mapMediaEvent(ctx context.Context, adapter voiceprotocol.MediaSessionAdapter, name string, payload map[string]interface{}) (voiceprotocol.MediaEvent, error) {
	switch voiceprotocol.MediaEventKind(name) {
	case voiceprotocol.MediaEventSessionStarted:
		return adapter.MapSessionStarted(ctx, payload)
	case voiceprotocol.MediaEventSessionClosed:
		return adapter.MapSessionClosed(ctx, payload)
	case voiceprotocol.MediaEventTrackMuted:
		return adapter.MapTrackMuted(ctx, payload)
	case voiceprotocol.MediaEventTrackUnmuted:
		return adapter.MapTrackUnmuted(ctx, payload)
	case voiceprotocol.MediaEventRecordingStart:
		return adapter.MapRecordingStarted(ctx, payload)
	case voiceprotocol.MediaEventRecordingStop:
		return adapter.MapRecordingStopped(ctx, payload)
	default:
		return voiceprotocol.MediaEvent{}, fmt.Errorf("unsupported media event %q", name)
	}
}
