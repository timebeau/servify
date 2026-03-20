package server

import (
	"servify/apps/server/internal/config"
	"servify/apps/server/internal/handlers"
	"servify/apps/server/internal/middleware"
	agentdelivery "servify/apps/server/internal/modules/agent/delivery"
	aidelivery "servify/apps/server/internal/modules/ai/delivery"
	analyticsdelivery "servify/apps/server/internal/modules/analytics/delivery"
	customerdelivery "servify/apps/server/internal/modules/customer/delivery"
	routingdelivery "servify/apps/server/internal/modules/routing/delivery"
	ticketdelivery "servify/apps/server/internal/modules/ticket/delivery"
	voicedelivery "servify/apps/server/internal/modules/voice/delivery"
	realtimeplatform "servify/apps/server/internal/platform/realtime"
	"servify/apps/server/internal/platform/voiceprotocol"
	"servify/apps/server/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// Dependencies contains the runtime services required to assemble the HTTP router.
type Dependencies struct {
	Config                   *config.Config
	Logger                   *logrus.Logger
	DB                       *gorm.DB
	AIService                services.AIServiceInterface
	AIHandlerService         aidelivery.HandlerService
	RealtimeGateway          realtimeplatform.RealtimeGateway
	RTCGateway               realtimeplatform.RTCGateway
	MessageRouter            *services.MessageRouter
	VoiceCoordinator         *voicedelivery.Coordinator
	VoiceProtocolRegistry    *voiceprotocol.Registry
	CustomerHandlerService   customerdelivery.HandlerService
	AgentHandlerService      agentdelivery.HandlerService
	AgentService             *services.AgentService
	TicketHandlerService     ticketdelivery.HandlerService
	TicketReaderService      *ticketdelivery.ReaderServiceAdapter
	TransferHandlerService   routingdelivery.HandlerService
	TransferService          *services.SessionTransferService
	SatisfactionService      *services.SatisfactionService
	WorkspaceService         *services.WorkspaceService
	MacroService             *services.MacroService
	AppIntegrationService    *services.AppIntegrationService
	CustomFieldService       *services.CustomFieldService
	StatisticsHandlerService analyticsdelivery.HandlerService
	SLAService               *services.SLAService
	ShiftService             *services.ShiftService
	AutomationService        *services.AutomationService
	KnowledgeDocService      *services.KnowledgeDocService
	SuggestionService        *services.SuggestionService
	GamificationService      *services.GamificationService
}

// BuildRouter assembles the HTTP routes and middleware around already-wired services.
func BuildRouter(deps Dependencies) *gin.Engine {
	r := gin.New()
	registerBaseMiddleware(r, deps.Config)
	registerHealthRoutes(r, deps)
	registerManagementRoutes(r, deps)
	registerPublicRoutes(r, deps)
	registerRealtimeRoutes(r, deps)
	registerStatic(r)
	return r
}

func registerManagementRoutes(r *gin.Engine, deps Dependencies) {
	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware(deps.Config))

	customersAPI := api.Group("/")
	customersAPI.Use(middleware.RequireResourcePermission("customers"))
	handlers.RegisterCustomerRoutes(customersAPI, handlers.NewCustomerHandler(deps.CustomerHandlerService, deps.Logger))

	agentsAPI := api.Group("/")
	agentsAPI.Use(middleware.RequireResourcePermission("agents"))
	handlers.RegisterAgentRoutes(agentsAPI, handlers.NewAgentHandler(deps.AgentHandlerService, deps.Logger))

	ticketsAPI := api.Group("/")
	ticketsAPI.Use(middleware.RequireResourcePermission("tickets"))
	handlers.RegisterTicketRoutes(ticketsAPI, handlers.NewTicketHandler(deps.TicketHandlerService, deps.Logger))

	sessionTransferAPI := api.Group("/")
	sessionTransferAPI.Use(middleware.RequireResourcePermission("session_transfer"))
	handlers.RegisterSessionTransferRoutes(sessionTransferAPI, handlers.NewSessionTransferHandler(deps.TransferHandlerService, deps.Logger))

	satisfactionAPI := api.Group("/")
	satisfactionAPI.Use(middleware.RequireResourcePermission("satisfaction"))
	handlers.RegisterSatisfactionRoutes(satisfactionAPI, handlers.NewSatisfactionHandler(deps.SatisfactionService, deps.Logger))

	workspaceAPI := api.Group("/")
	workspaceAPI.Use(middleware.RequireResourcePermission("workspace"))
	handlers.RegisterWorkspaceRoutes(workspaceAPI, handlers.NewWorkspaceHandler(deps.WorkspaceService))

	macrosAPI := api.Group("/")
	macrosAPI.Use(middleware.RequireResourcePermission("macros"))
	handlers.RegisterMacroRoutes(macrosAPI, handlers.NewMacroHandler(deps.MacroService))

	integrationsAPI := api.Group("/")
	integrationsAPI.Use(middleware.RequireResourcePermission("integrations"))
	handlers.RegisterAppIntegrationRoutes(integrationsAPI, handlers.NewAppMarketHandler(deps.AppIntegrationService))

	customFieldsAPI := api.Group("/")
	customFieldsAPI.Use(middleware.RequireResourcePermission("custom_fields"))
	handlers.RegisterCustomFieldRoutes(customFieldsAPI, handlers.NewCustomFieldHandler(deps.CustomFieldService))

	statisticsAPI := api.Group("/")
	statisticsAPI.Use(middleware.RequireResourcePermission("statistics"))
	handlers.RegisterStatisticsRoutes(statisticsAPI, handlers.NewStatisticsHandler(deps.StatisticsHandlerService, deps.Logger))

	slaAPI := api.Group("/")
	slaAPI.Use(middleware.RequireResourcePermission("sla"))
	handlers.RegisterSLARoutes(slaAPI, handlers.NewSLAHandler(deps.SLAService, deps.TicketReaderService))

	shiftAPI := api.Group("/")
	shiftAPI.Use(middleware.RequireResourcePermission("shift"))
	handlers.RegisterShiftRoutes(shiftAPI, handlers.NewShiftHandler(deps.ShiftService))

	automationAPI := api.Group("/")
	automationAPI.Use(middleware.RequireResourcePermission("automation"))
	handlers.RegisterAutomationRoutes(automationAPI, handlers.NewAutomationHandler(deps.AutomationService))

	knowledgeAPI := api.Group("/")
	knowledgeAPI.Use(middleware.RequireResourcePermission("knowledge"))
	handlers.RegisterKnowledgeDocRoutes(knowledgeAPI, handlers.NewKnowledgeDocHandler(deps.KnowledgeDocService))

	assistAPI := api.Group("/")
	assistAPI.Use(middleware.RequireResourcePermission("assist"))
	handlers.RegisterSuggestionRoutes(assistAPI, handlers.NewSuggestionHandler(deps.SuggestionService))

	gamificationAPI := api.Group("/")
	gamificationAPI.Use(middleware.RequireResourcePermission("gamification"))
	handlers.RegisterGamificationRoutes(gamificationAPI, handlers.NewGamificationHandler(deps.GamificationService))

	voiceAPI := api.Group("/")
	voiceAPI.Use(middleware.RequireResourcePermission("assist"))
	handlers.RegisterVoiceRoutes(voiceAPI, handlers.NewVoiceHandler(deps.VoiceCoordinator, deps.VoiceProtocolRegistry))
}

func registerPublicRoutes(r *gin.Engine, deps Dependencies) {
	public := r.Group("/public")
	handlers.RegisterCSATSurveyRoutes(public, handlers.NewCSATSurveyHandler(deps.SatisfactionService))
	handlers.RegisterPublicKnowledgeBaseRoutes(public, handlers.NewKnowledgeDocHandler(deps.KnowledgeDocService))
	public.GET("/portal/config", handlers.NewPortalConfigHandler(deps.Config).Get)
}

func registerRealtimeRoutes(r *gin.Engine, deps Dependencies) {
	v1 := r.Group("/api/v1")

	wsHandler := handlers.NewWebSocketHandler(deps.RealtimeGateway)
	v1.GET("/ws", wsHandler.HandleWebSocket)
	v1.GET("/ws/stats", wsHandler.GetStats)

	webrtcHandler := handlers.NewWebRTCHandler(deps.RTCGateway)
	v1.GET("/webrtc/stats", webrtcHandler.GetStats)
	v1.GET("/webrtc/connections", webrtcHandler.GetConnections)

	messageHandler := handlers.NewMessageHandler(deps.MessageRouter)
	v1.GET("/messages/platforms", messageHandler.GetPlatformStats)

	aiHandler := handlers.NewAIHandler(deps.AIHandlerService)
	aiAPI := v1.Group("/ai")
	aiAPI.POST("/query", aiHandler.ProcessQuery)
	aiAPI.GET("/status", aiHandler.GetStatus)
	aiAPI.GET("/metrics", aiHandler.GetMetrics)
	if deps.Config != nil && deps.Config.WeKnora.Enabled {
		aiAPI.POST("/knowledge/upload", aiHandler.UploadDocument)
		aiAPI.POST("/knowledge/sync", aiHandler.SyncKnowledgeBase)
		aiAPI.PUT("/weknora/enable", aiHandler.EnableWeKnora)
		aiAPI.PUT("/weknora/disable", aiHandler.DisableWeKnora)
		aiAPI.POST("/circuit-breaker/reset", aiHandler.ResetCircuitBreaker)
	}

	ingest := handlers.NewMetricsIngestHandler(handlers.NewMetricsAggregator())
	v1.POST("/metrics/ingest", ingest.Ingest)
}
