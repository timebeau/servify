package server

import (
	"context"

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
	voiceprovidermock "servify/apps/server/internal/modules/voice/provider/mock"
	"servify/apps/server/internal/platform/eventbus"
	"servify/apps/server/internal/platform/pstnprovider"
	realtimeplatform "servify/apps/server/internal/platform/realtime"
	"servify/apps/server/internal/platform/sip"
	"servify/apps/server/internal/platform/sipws"
	"servify/apps/server/internal/platform/voiceprotocol"
	"servify/apps/server/internal/services"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// Runtime owns the service graph used by the HTTP server.
type Runtime struct {
	Config *config.Config
	Logger *logrus.Logger
	DB     *gorm.DB
	Bus    eventbus.Bus

	AIService                services.AIServiceInterface
	AIHandlerService         aidelivery.HandlerService
	wsRuntime                websocketRunner
	RealtimeGateway          realtimeplatform.RealtimeGateway
	RTCGateway               realtimeplatform.RTCGateway
	MessageRouter            services.MessageRouterRuntime
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
}

type websocketRunner interface {
	Run()
}

// BuildRuntime wires the current modular-monolith runtime behind an explicit assembly boundary.
func BuildRuntime(cfg *config.Config, logger *logrus.Logger, db *gorm.DB, bus eventbus.Bus) (*Runtime, error) {
	rt := &Runtime{
		Config: cfg,
		Logger: logger,
		DB:     db,
		Bus:    bus,
	}

	aiAssembly, err := BuildAIAssembly(cfg, logger, AIAssemblyOptions{})
	if err != nil {
		return nil, err
	}
	rt.AIService = aiAssembly.RuntimeService
	rt.AIHandlerService = aiAssembly.Service

	wsHub := services.NewWebSocketHub()
	rt.wsRuntime = wsHub
	rt.RealtimeGateway = realtimeplatform.NewWebSocketAdapter(wsHub)
	wsHub.SetDB(db)

	conversationRepo := conversationinfra.NewGormRepository(db)
	conversationService := conversationapp.NewService(conversationRepo, bus)
	wsHub.SetConversationMessageWriter(conversationdelivery.NewWebSocketMessageAdapter(conversationService))

	routingRepo := routinginfra.NewGormRepository(db)
	routingService := routingapp.NewService(routingRepo, bus)

	webrtcService := services.NewWebRTCService(cfg.WebRTC.STUNServer, wsHub)
	rt.RTCGateway = realtimeplatform.NewWebRTCAdapter(webrtcService)

	rt.MessageRouter = services.NewMessageRouter(rt.AIService, wsHub, db)
	wsHub.SetAIService(rt.AIService)

	voiceService := voiceapp.NewService(voiceinfra.NewInMemoryRepository(), bus)
	recordingService := voiceapp.NewRecordingService(
		voiceprovidermock.NewRecordingProvider(),
		voiceinfra.NewInMemoryRecordingRepository(),
		bus,
	)
	transcriptService := voiceapp.NewTranscriptService(
		voiceprovidermock.NewTranscriptProvider(),
		voiceinfra.NewInMemoryTranscriptRepository(),
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
	automationService := services.NewAutomationService(db, logger)
	rt.AutomationHandlerService = services.NewAutomationHandlerAdapter(automationService)
	automationService.SetEventBus(bus)
	slaService.SetAutomationService(automationService)

	customerService := services.NewCustomerService(db, logger)
	rt.CustomerHandlerService = customerdelivery.NewHandlerServiceAdapter(customerService)
	agentAssembly := services.BuildAgentServiceAssembly(db, logger)
	rt.AgentHandlerService = agentAssembly.Service
	go agentAssembly.Maintenance.Start()

	ticketService := services.NewTicketService(db, logger, slaService)
	ticketService.SetEventBus(bus)
	ticketService.SetAutomationService(automationService)

	transferService := services.NewSessionTransferService(db, logger, rt.AIService, agentAssembly.Service, wsHub)
	transferService.SetRoutingAdapter(routingdelivery.NewSessionTransferAdapter(routingService, bus))
	transferService.SetTicketRuntime(ticketdelivery.NewRuntimeAdapter(bus))
	transferService.SetConversationRuntime(conversationdelivery.NewRuntimeAdapter(bus))
	transferService.SetAgentRuntime(agentdelivery.NewTransferRuntimeAdapter())
	rt.TransferHandlerService = transferService

	statisticsService := services.NewStatisticsService(db, logger)
	statisticsService.SetEventBus(bus)
	rt.StatisticsHandlerService = statisticsService

	satisfactionService := services.NewSatisfactionService(db, logger)
	rt.SatisfactionService = satisfactionService
	ticketService.SetSatisfactionService(satisfactionService)

	rt.ShiftService = services.NewShiftService(db, logger)
	rt.WorkspaceService = services.NewWorkspaceService(db, agentAssembly.Service)
	rt.MacroService = services.NewMacroService(db)
	rt.AppIntegrationService = services.NewAppIntegrationService(db, logger)
	rt.CustomFieldService = services.NewCustomFieldService(db)
	knowledgeDocService := services.NewKnowledgeDocService(db)
	rt.KnowledgeDocHandler = knowledgedelivery.NewHandlerServiceAdapter(knowledgeDocService)
	rt.SuggestionService = services.NewSuggestionService(db)
	rt.GamificationService = services.NewGamificationService(db)
	rt.TicketHandlerService = ticketdelivery.NewHandlerServiceAdapter(db, ticketService.ModuleCommandService(), ticketService.Orchestrator())
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

func (rt *Runtime) RouterDependencies() Dependencies {
	return Dependencies{
		Config:                   rt.Config,
		Logger:                   rt.Logger,
		DB:                       rt.DB,
		AIService:                rt.AIService,
		AIHandlerService:         rt.AIHandlerService,
		RealtimeGateway:          rt.RealtimeGateway,
		RTCGateway:               rt.RTCGateway,
		MessageRouter:            rt.MessageRouter,
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
	}
}
