package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	voiceapp "servify/apps/server/internal/modules/voice/application"
	voicedelivery "servify/apps/server/internal/modules/voice/delivery"
	voiceinfra "servify/apps/server/internal/modules/voice/infra"
	voiceprovidermock "servify/apps/server/internal/modules/voice/provider/mock"
	"servify/apps/server/internal/platform/eventbus"
	"servify/apps/server/internal/platform/pstnprovider"
	"servify/apps/server/internal/platform/sip"
	"servify/apps/server/internal/platform/sipws"
	"servify/apps/server/internal/platform/voiceprotocol"

	"github.com/gin-gonic/gin"
)

func newVoiceCoordinatorForTest() *voicedelivery.Coordinator {
	bus := &voiceTestBus{}
	callService := voiceapp.NewService(voiceinfra.NewInMemoryRepository(), bus)
	recordingService := voiceapp.NewRecordingService(voiceprovidermock.NewRecordingProvider(), voiceinfra.NewInMemoryRecordingRepository(), bus)
	transcriptService := voiceapp.NewTranscriptService(voiceprovidermock.NewTranscriptProvider(), voiceinfra.NewInMemoryTranscriptRepository(), bus)
	return voicedelivery.NewCoordinator(callService, recordingService, transcriptService)
}

func newVoiceRegistryForTest() *voiceprotocol.Registry {
	registry := voiceprotocol.NewRegistry()
	_ = registry.RegisterSignaling(sip.NewVoiceProtocolAdapter())
	_ = registry.RegisterSignaling(sipws.NewAdapter())
	_ = registry.RegisterSignaling(pstnprovider.NewAdapter())
	_ = registry.RegisterMedia(voicedelivery.NewWebRTCAdapter(voiceapp.NewService(voiceinfra.NewInMemoryRepository(), &voiceTestBus{})))
	_ = registry.RegisterMedia(voicedelivery.NewRTPAdapter())
	_ = registry.RegisterMedia(voicedelivery.NewSRTPAdapter())
	return registry
}

type voiceTestBus struct{}

func (b *voiceTestBus) Publish(ctx context.Context, event eventbus.Event) error {
	return nil
}

func TestVoiceHandlerStartRecording(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewVoiceHandler(newVoiceCoordinatorForTest(), newVoiceRegistryForTest())
	r := gin.New()
	r.POST("/voice/recordings/start", handler.StartRecording)

	body, _ := json.Marshal(map[string]string{"call_id": "call-1", "provider": "mock"})
	req := httptest.NewRequest(http.MethodPost, "/voice/recordings/start", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body=%s", w.Code, w.Body.String())
	}
}

func TestVoiceHandlerAppendTranscriptAndList(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewVoiceHandler(newVoiceCoordinatorForTest(), newVoiceRegistryForTest())
	r := gin.New()
	r.POST("/voice/transcripts", handler.AppendTranscript)
	r.GET("/voice/transcripts", handler.ListTranscripts)

	body, _ := json.Marshal(map[string]interface{}{"call_id": "call-1", "content": "hello", "language": "en", "finalized": true})
	req := httptest.NewRequest(http.MethodPost, "/voice/transcripts", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("append expected 200, got %d, body=%s", w.Code, w.Body.String())
	}

	req = httptest.NewRequest(http.MethodGet, "/voice/transcripts?call_id=call-1", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("list expected 200, got %d, body=%s", w.Code, w.Body.String())
	}
}

func TestVoiceHandlerHandleProtocolCallEvent(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewVoiceHandler(newVoiceCoordinatorForTest(), newVoiceRegistryForTest())
	r := gin.New()
	r.POST("/voice/protocols/:protocol/call-events/:event", handler.HandleProtocolCallEvent)

	body, _ := json.Marshal(map[string]interface{}{
		"payload": map[string]interface{}{
			"call_id": "call-sip-1",
			"from":    "1001",
			"to":      "1002",
		},
	})
	req := httptest.NewRequest(http.MethodPost, "/voice/protocols/sip/call-events/invite", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("protocol call event expected 200, got %d, body=%s", w.Code, w.Body.String())
	}
}

func TestVoiceHandlerListProtocols(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewVoiceHandler(newVoiceCoordinatorForTest(), newVoiceRegistryForTest())
	r := gin.New()
	r.GET("/voice/protocols", handler.ListProtocols)

	req := httptest.NewRequest(http.MethodGet, "/voice/protocols", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("list protocols expected 200, got %d, body=%s", w.Code, w.Body.String())
	}
}
