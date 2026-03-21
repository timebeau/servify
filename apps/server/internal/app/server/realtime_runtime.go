package server

import (
	"context"

	"servify/apps/server/internal/config"
	aidelivery "servify/apps/server/internal/modules/ai/delivery"
	conversationapp "servify/apps/server/internal/modules/conversation/application"
	conversationdelivery "servify/apps/server/internal/modules/conversation/delivery"
	conversationinfra "servify/apps/server/internal/modules/conversation/infra"
	realtimeplatform "servify/apps/server/internal/platform/realtime"
	"servify/apps/server/internal/services"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// RealtimeRuntime assembles realtime and AI-facing primitives used by lightweight runtimes.
type RealtimeRuntime struct {
	Config           *config.Config
	Logger           *logrus.Logger
	DB               *gorm.DB
	AIService        services.AIServiceInterface
	AIHandlerService aidelivery.HandlerService
	WSHub            *services.WebSocketHub
	WebRTCService    *services.WebRTCService
	RealtimeGateway  realtimeplatform.RealtimeGateway
	RTCGateway       realtimeplatform.RTCGateway
	MessageRouter    services.MessageRouterRuntime
}

func BuildRealtimeRuntime(cfg *config.Config, logger *logrus.Logger, db *gorm.DB, ai services.AIServiceInterface, handlerAI aidelivery.HandlerService) *RealtimeRuntime {
	wsHub := services.NewWebSocketHub()
	if db != nil {
		wsHub.SetDB(db)
		conversationRepo := conversationinfra.NewGormRepository(db)
		conversationService := conversationapp.NewService(conversationRepo, nil)
		wsHub.SetConversationMessageWriter(conversationdelivery.NewWebSocketMessageAdapter(conversationService))
	}
	wsHub.SetAIService(ai)

	webrtcService := services.NewWebRTCService(cfg.WebRTC.STUNServer, wsHub)
	return &RealtimeRuntime{
		Config:           cfg,
		Logger:           logger,
		DB:               db,
		AIService:        ai,
		AIHandlerService: handlerAI,
		WSHub:            wsHub,
		WebRTCService:    webrtcService,
		RealtimeGateway:  realtimeplatform.NewWebSocketAdapter(wsHub),
		RTCGateway:       realtimeplatform.NewWebRTCAdapter(webrtcService),
		MessageRouter:    services.NewMessageRouter(ai, wsHub, db),
	}
}

func (rt *RealtimeRuntime) Start() error {
	go rt.WSHub.Run()
	return rt.MessageRouter.Start()
}

func (rt *RealtimeRuntime) Stop(context.Context) error {
	if rt.MessageRouter == nil {
		return nil
	}
	return rt.MessageRouter.Stop()
}
