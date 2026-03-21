package server

import (
	"context"

	"servify/apps/server/internal/config"
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
	WSHub                    *services.WebSocketHub
	RealtimeGateway          realtimeplatform.RealtimeGateway
	RTCGateway               realtimeplatform.RTCGateway
	MessageRouter            *services.MessageRouter
	VoiceCoordinator         *voicedelivery.Coordinator
	VoiceProtocolRegistry    *voiceprotocol.Registry
	CustomerHandlerService   customerdelivery.HandlerService
	CustomerService          *services.CustomerService
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
	StatisticsService        *services.StatisticsService
	SLAService               *services.SLAService
	ShiftService             *services.ShiftService
	AutomationHandlerService automationdelivery.HandlerService
	AutomationService        *services.AutomationService
	KnowledgeDocHandler      knowledgedelivery.HandlerService
	KnowledgeDocService      *services.KnowledgeDocService
	SuggestionService        *services.SuggestionService
	GamificationService      *services.GamificationService
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

	rt.WSHub = services.NewWebSocketHub()
	rt.RealtimeGateway = realtimeplatform.NewWebSocketAdapter(rt.WSHub)
	rt.WSHub.SetDB(db)

	conversationRepo := conversationinfra.NewGormRepository(db)
	conversationService := conversationapp.NewService(conversationRepo, bus)
	rt.WSHub.SetConversationMessageWriter(conversationdelivery.NewWebSocketMessageAdapter(conversationService))

	routingRepo := routinginfra.NewGormRepository(db)
	routingService := routingapp.NewService(routingRepo, bus)

	webrtcService := services.NewWebRTCService(cfg.WebRTC.STUNServer, rt.WSHub)
	rt.RTCGateway = realtimeplatform.NewWebRTCAdapter(webrtcService)

	rt.MessageRouter = services.NewMessageRouter(rt.AIService, rt.WSHub, db)
	rt.WSHub.SetAIService(rt.AIService)

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

	rt.SLAService = services.NewSLAService(db, logger)
	rt.AutomationService = services.NewAutomationService(db, logger)
	rt.AutomationHandlerService = services.NewAutomationHandlerAdapter(rt.AutomationService)
	rt.AutomationService.SetEventBus(bus)
	rt.SLAService.SetAutomationService(rt.AutomationService)

	rt.CustomerService = services.NewCustomerService(db, logger)
	rt.CustomerHandlerService = customerdelivery.NewHandlerServiceAdapter(rt.CustomerService)
	rt.AgentService = services.NewAgentService(db, logger)
	rt.AgentHandlerService = rt.AgentService

	ticketService := services.NewTicketService(db, logger, rt.SLAService)
	ticketService.SetEventBus(bus)
	ticketService.SetAutomationService(rt.AutomationService)

	rt.TransferService = services.NewSessionTransferService(db, logger, rt.AIService, rt.AgentService, rt.WSHub)
	rt.TransferService.SetRoutingAdapter(routingdelivery.NewSessionTransferAdapter(routingService, bus))
	rt.TransferService.SetTicketRuntime(ticketdelivery.NewRuntimeAdapter(bus))
	rt.TransferService.SetConversationRuntime(conversationdelivery.NewRuntimeAdapter(bus))
	rt.TransferService.SetAgentRuntime(agentdelivery.NewTransferRuntimeAdapter())
	rt.TransferHandlerService = rt.TransferService

	rt.StatisticsService = services.NewStatisticsService(db, logger)
	rt.StatisticsService.SetEventBus(bus)
	rt.StatisticsHandlerService = rt.StatisticsService

	rt.SatisfactionService = services.NewSatisfactionService(db, logger)
	ticketService.SetSatisfactionService(rt.SatisfactionService)

	rt.ShiftService = services.NewShiftService(db, logger)
	rt.WorkspaceService = services.NewWorkspaceService(db, rt.AgentService)
	rt.MacroService = services.NewMacroService(db)
	rt.AppIntegrationService = services.NewAppIntegrationService(db, logger)
	rt.CustomFieldService = services.NewCustomFieldService(db)
	rt.KnowledgeDocService = services.NewKnowledgeDocService(db)
	rt.KnowledgeDocHandler = knowledgedelivery.NewHandlerServiceAdapter(rt.KnowledgeDocService)
	rt.SuggestionService = services.NewSuggestionService(db)
	rt.GamificationService = services.NewGamificationService(db)
	rt.TicketHandlerService = ticketdelivery.NewHandlerServiceAdapter(db, ticketService.ModuleCommandService(), ticketService.Orchestrator())
	rt.TicketReaderService = ticketdelivery.NewReaderServiceAdapter(db)

	rt.WSHub.SetSessionTransferService(rt.TransferService)
	return rt, nil
}

func (rt *Runtime) Start() error {
	go rt.WSHub.Run()
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
		AgentService:             rt.AgentService,
		TicketHandlerService:     rt.TicketHandlerService,
		TicketReaderService:      rt.TicketReaderService,
		TransferHandlerService:   rt.TransferHandlerService,
		TransferService:          rt.TransferService,
		SatisfactionService:      rt.SatisfactionService,
		WorkspaceService:         rt.WorkspaceService,
		MacroService:             rt.MacroService,
		AppIntegrationService:    rt.AppIntegrationService,
		CustomFieldService:       rt.CustomFieldService,
		StatisticsHandlerService: rt.StatisticsHandlerService,
		SLAService:               rt.SLAService,
		ShiftService:             rt.ShiftService,
		AutomationHandlerService: rt.AutomationHandlerService,
		AutomationService:        rt.AutomationService,
		KnowledgeDocHandler:      rt.KnowledgeDocHandler,
		SuggestionService:        rt.SuggestionService,
		GamificationService:      rt.GamificationService,
	}
}
