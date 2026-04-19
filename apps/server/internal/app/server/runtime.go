package server

import (
	"context"
	"net/http"

	"servify/apps/server/internal/config"
	"servify/apps/server/internal/handlers"
	agentdelivery "servify/apps/server/internal/modules/agent/delivery"
	aidelivery "servify/apps/server/internal/modules/ai/delivery"
	analyticsdelivery "servify/apps/server/internal/modules/analytics/delivery"
	automationdelivery "servify/apps/server/internal/modules/automation/delivery"
	conversationapp "servify/apps/server/internal/modules/conversation/application"
	conversationdelivery "servify/apps/server/internal/modules/conversation/delivery"
	conversationinfra "servify/apps/server/internal/modules/conversation/infra"
	customerdelivery "servify/apps/server/internal/modules/customer/delivery"
	knowledgedelivery "servify/apps/server/internal/modules/knowledge/delivery"
	routingapp "servify/apps/server/internal/modules/routing/application"
	routingdelivery "servify/apps/server/internal/modules/routing/delivery"
	routinginfra "servify/apps/server/internal/modules/routing/infra"
	ticketdelivery "servify/apps/server/internal/modules/ticket/delivery"
	voiceapp "servify/apps/server/internal/modules/voice/application"
	voicedelivery "servify/apps/server/internal/modules/voice/delivery"
	voiceinfra "servify/apps/server/internal/modules/voice/infra"
	svcmetrics "servify/apps/server/internal/observability/metrics"
	"servify/apps/server/internal/platform/eventbus"
	"servify/apps/server/internal/platform/pstnprovider"
	realtimeplatform "servify/apps/server/internal/platform/realtime"
	"servify/apps/server/internal/platform/sip"
	"servify/apps/server/internal/platform/sipws"
	"servify/apps/server/internal/platform/voiceprotocol"
	"servify/apps/server/internal/services"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// Runtime owns the service graph used by the HTTP server.
type Runtime struct {
	Config *config.Config
	Logger *logrus.Logger
	DB     *gorm.DB
	Redis  *redis.Client
	Bus    eventbus.Bus

	AIService                aidelivery.RuntimeService
	AIHandlerService         aidelivery.HandlerService
	wsRuntime                websocketRunner
	RealtimeGateway          realtimeplatform.RealtimeGateway
	RTCGateway               realtimeplatform.RTCGateway
	MessageRouter            services.MessageRouterRuntime
	ConversationHandler      conversationdelivery.HandlerService
	VoiceCoordinator         *voicedelivery.Coordinator
	VoiceProtocolRegistry    *voiceprotocol.Registry
	CustomerHandlerService   customerdelivery.HandlerService
	AgentHandlerService      agentdelivery.HandlerService
	TicketHandlerService     ticketdelivery.HandlerService
	TicketReaderService      *ticketdelivery.ReaderServiceAdapter
	TransferHandlerService   routingdelivery.HandlerService
	SatisfactionService      handlers.SatisfactionService
	WorkspaceService         services.WorkspaceOverviewReader
	MacroService             handlers.MacroService
	AppIntegrationService    handlers.AppMarketService
	CustomFieldService       handlers.CustomFieldService
	StatisticsHandlerService analyticsdelivery.HandlerService
	SLAService               handlers.SLAService
	ShiftService             handlers.ShiftService
	AutomationHandlerService automationdelivery.HandlerService
	KnowledgeDocHandler      knowledgedelivery.HandlerService
	SuggestionService        handlers.SuggestionService
	GamificationService      handlers.GamificationService
	HTTPMetrics              *svcmetrics.HTTPMetrics

	// Private fields for worker access only
	statisticsService *services.StatisticsService
	slaService        *services.SLAService
}

type websocketRunner interface {
	Run()
}

// BuildRuntime wires the current modular-monolith runtime behind an explicit assembly boundary.
func BuildRuntime(cfg *config.Config, logger *logrus.Logger, db *gorm.DB, redisClient *redis.Client, bus eventbus.Bus) (*Runtime, error) {
	rt := &Runtime{
		Config: cfg,
		Logger: logger,
		DB:     db,
		Redis:  redisClient,
		Bus:    bus,
	}

	// Initialize observability metrics if monitoring is enabled.
	if cfg.Monitoring.Enabled {
		svcmetrics.DefaultRegistry.RegisterGoCollector()
		svcmetrics.DefaultRegistry.RegisterProcessCollector()
		rt.HTTPMetrics = svcmetrics.NewHTTPMetrics(svcmetrics.DefaultRegistry)
	}

	aiAssembly, err := BuildAIAssembly(cfg, logger, AIAssemblyOptions{})
	if err != nil {
		return nil, err
	}
	rt.AIService = NewScopedAIRuntimeService(cfg, logger, db, aiAssembly.RuntimeService)
	rt.AIHandlerService = NewScopedAIHandlerService(cfg, logger, db, aiAssembly.Service)

	wsHub := services.NewWebSocketHub()
	rt.wsRuntime = wsHub
	rt.RealtimeGateway = realtimeplatform.NewWebSocketAdapter(wsHub)

	conversationRepo := conversationinfra.NewGormRepository(db)
	conversationService := conversationapp.NewService(conversationRepo, bus)
	rt.ConversationHandler = conversationdelivery.NewHandlerService(conversationService)
	wsHub.SetConversationMessageWriter(conversationdelivery.NewWebSocketMessageAdapter(conversationService))

	routingRepo := routinginfra.NewGormRepository(db)
	routingService := routingapp.NewService(routingRepo, bus)

	webrtcService := services.NewWebRTCService(cfg.WebRTC.STUNServer, wsHub)
	wsHub.SetWebRTCService(webrtcService)
	rt.RTCGateway = realtimeplatform.NewWebRTCAdapter(webrtcService)

	rt.MessageRouter = services.NewMessageRouter(rt.AIService, wsHub, db)
	wsHub.SetAIService(rt.AIService)

	voiceService := voiceapp.NewService(voiceinfra.NewGormRepository(db), bus)
	recordingProvider, err := buildVoiceRecordingProvider(cfg, logger)
	if err != nil {
		return nil, err
	}
	transcriptProvider, err := buildVoiceTranscriptProvider(cfg, logger)
	if err != nil {
		return nil, err
	}
	recordingService := voiceapp.NewRecordingService(
		recordingProvider,
		voiceinfra.NewGormRecordingRepository(db),
		bus,
	)
	transcriptService := voiceapp.NewTranscriptService(
		transcriptProvider,
		voiceinfra.NewGormTranscriptRepository(db),
		bus,
	)
	rt.VoiceCoordinator = voicedelivery.NewCoordinator(voiceService, recordingService, transcriptService)
	webrtcService.SetVoiceLifecycle(rt.VoiceCoordinator)
	rt.VoiceProtocolRegistry = voiceprotocol.NewRegistry()
	_ = rt.VoiceProtocolRegistry.RegisterSignaling(sip.NewVoiceProtocolAdapter())
	_ = rt.VoiceProtocolRegistry.RegisterSignaling(sipws.NewAdapter())
	_ = rt.VoiceProtocolRegistry.RegisterSignaling(pstnprovider.NewAdapter())
	_ = rt.VoiceProtocolRegistry.RegisterMedia(voicedelivery.NewWebRTCAdapter(voiceService))
	_ = rt.VoiceProtocolRegistry.RegisterMedia(voicedelivery.NewRTPAdapter())
	_ = rt.VoiceProtocolRegistry.RegisterMedia(voicedelivery.NewSRTPAdapter())

	slaService := services.NewSLAService(db, logger)
	rt.SLAService = slaService
	rt.slaService = slaService
	automationService := services.NewAutomationService(db, logger)
	rt.AutomationHandlerService = automationdelivery.NewHandlerService(db)
	automationService.SetEventBus(bus)
	slaService.SetAutomationService(automationService)

	rt.CustomerHandlerService = customerdelivery.NewHandlerService(db)
	agentAssembly := services.BuildAgentServiceAssembly(db, logger, redisClient)
	rt.AgentHandlerService = agentAssembly.Service
	go agentAssembly.Maintenance.Start()

	transferService := routingdelivery.NewHandlerService(routingdelivery.HandlerDependencies{
		DB:           db,
		Logger:       logger,
		AI:           rt.AIService,
		Agents:       agentAssembly.Service,
		Notifier:     newRoutingTransferNotifier(rt.RealtimeGateway),
		Routing:      routingdelivery.NewSessionTransferAdapter(routingService, bus),
		Tickets:      ticketdelivery.NewRuntimeAdapter(bus),
		Conversation: conversationdelivery.NewRuntimeAdapter(db, bus),
		AgentLoad:    agentdelivery.NewTransferRuntimeAdapter(),
	})
	rt.TransferHandlerService = transferService

	statisticsService := services.NewStatisticsService(db, logger)
	statisticsService.SetEventBus(bus)
	rt.StatisticsHandlerService = statisticsService
	rt.statisticsService = statisticsService

	satisfactionService := services.NewSatisfactionService(db, logger)
	rt.SatisfactionService = satisfactionService

	rt.ShiftService = services.NewShiftService(db, logger)
	rt.WorkspaceService = services.NewWorkspaceService(db, agentAssembly.Service)
	rt.MacroService = services.NewMacroService(db)
	rt.AppIntegrationService = services.NewAppIntegrationService(db, logger)
	rt.CustomFieldService = services.NewCustomFieldService(db)
	rt.KnowledgeDocHandler = knowledgedelivery.NewHandlerServiceWithProvider(db, aiAssembly.KnowledgeProvider(cfg))
	rt.SuggestionService = services.NewSuggestionService(db)
	rt.GamificationService = services.NewGamificationService(db)
	rt.TicketHandlerService = ticketdelivery.NewHandlerServiceWithDependencies(ticketdelivery.HandlerAssemblyDependencies{
		DB:           db,
		Logger:       logger,
		Bus:          bus,
		SLA:          slaService,
		Satisfaction: satisfactionService,
	})
	rt.TicketReaderService = ticketdelivery.NewReaderServiceAdapter(db)

	wsHub.SetSessionTransferService(transferService)
	return rt, nil
}

func (rt *Runtime) Start() error {
	if rt.wsRuntime != nil {
		go rt.wsRuntime.Run()
	}
	return rt.MessageRouter.Start()
}

func (rt *Runtime) Stop(context.Context) error {
	if rt.MessageRouter == nil {
		return nil
	}
	return rt.MessageRouter.Stop()
}

// Router builds the HTTP router for this runtime.
func (rt *Runtime) Router() http.Handler {
	return BuildRouter(rt.RouterDependencies())
}

// StatisticsServiceForWorker returns the concrete statistics service for worker use.
func (rt *Runtime) StatisticsServiceForWorker() *services.StatisticsService {
	return rt.statisticsService
}

// SLAServiceForWorker returns the concrete SLA service for worker use.
func (rt *Runtime) SLAServiceForWorker() *services.SLAService {
	return rt.slaService
}

func (rt *Runtime) RouterDependencies() Dependencies {
	return Dependencies{
		Config:                   rt.Config,
		Logger:                   rt.Logger,
		DB:                       rt.DB,
		Redis:                    rt.Redis,
		AIService:                rt.AIService,
		AIHandlerService:         rt.AIHandlerService,
		RealtimeGateway:          rt.RealtimeGateway,
		RTCGateway:               rt.RTCGateway,
		MessageRouter:            rt.MessageRouter,
		ConversationHandler:      rt.ConversationHandler,
		VoiceCoordinator:         rt.VoiceCoordinator,
		VoiceProtocolRegistry:    rt.VoiceProtocolRegistry,
		CustomerHandlerService:   rt.CustomerHandlerService,
		AgentHandlerService:      rt.AgentHandlerService,
		TicketHandlerService:     rt.TicketHandlerService,
		TicketReaderService:      rt.TicketReaderService,
		TransferHandlerService:   rt.TransferHandlerService,
		SatisfactionService:      rt.SatisfactionService,
		WorkspaceService:         rt.WorkspaceService,
		MacroService:             rt.MacroService,
		AppIntegrationService:    rt.AppIntegrationService,
		CustomFieldService:       rt.CustomFieldService,
		StatisticsHandlerService: rt.StatisticsHandlerService,
		SLAService:               rt.SLAService,
		ShiftService:             rt.ShiftService,
		AutomationHandlerService: rt.AutomationHandlerService,
		KnowledgeDocHandler:      rt.KnowledgeDocHandler,
		SuggestionService:        rt.SuggestionService,
		GamificationService:      rt.GamificationService,
		HTTPMetrics:              rt.HTTPMetrics,
	}
}
