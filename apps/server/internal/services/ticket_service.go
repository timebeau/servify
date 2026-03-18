package services

import (
	"context"
	"fmt"

	"servify/apps/server/internal/models"
	ticketapp "servify/apps/server/internal/modules/ticket/application"
	ticketcontract "servify/apps/server/internal/modules/ticket/contract"
	ticketinfra "servify/apps/server/internal/modules/ticket/infra"
	ticketorchestration "servify/apps/server/internal/modules/ticket/orchestration"
	"servify/apps/server/internal/platform/eventbus"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// TicketService 工单管理服务
type TicketService struct {
	db           *gorm.DB
	logger       *logrus.Logger
	slaService   *SLAService
	automation   *AutomationService
	satisfaction *SatisfactionService
	orchestrator *ticketorchestration.TicketOrchestrator
	moduleRepo   *ticketinfra.GormRepository
	moduleQuery  *ticketapp.QueryService
	moduleCmd    *ticketapp.CommandService
	moduleBus    eventbus.Bus
}

// NewTicketService 创建工单服务
func NewTicketService(db *gorm.DB, logger *logrus.Logger, slaService *SLAService) *TicketService {
	if logger == nil {
		logger = logrus.New()
	}

	service := &TicketService{
		db:         db,
		logger:     logger,
		slaService: slaService,
	}
	service.rebuildTicketModule()
	return service
}

// SetAutomationService 注入自动化服务
func (s *TicketService) SetAutomationService(automation *AutomationService) {
	s.automation = automation
}

// SetSatisfactionService 注入满意度服务
func (s *TicketService) SetSatisfactionService(satisfaction *SatisfactionService) {
	s.satisfaction = satisfaction
	s.rebuildTicketOrchestrator()
}

// SetEventBus injects the runtime event bus for the modular ticket service path.
func (s *TicketService) SetEventBus(bus eventbus.Bus) {
	s.moduleBus = bus
	s.rebuildTicketModule()
}

func (s *TicketService) rebuildTicketModule() {
	if s.db == nil {
		s.moduleRepo = nil
		s.moduleQuery = nil
		s.moduleCmd = nil
		return
	}
	s.moduleRepo = ticketinfra.NewGormRepository(s.db)
	s.moduleQuery = ticketapp.NewQueryService(s.moduleRepo)
	s.moduleCmd = ticketapp.NewCommandServiceWithBus(s.moduleRepo, s.moduleBus)
	s.rebuildTicketOrchestrator()
}

func (s *TicketService) rebuildTicketOrchestrator() {
	s.orchestrator = ticketorchestration.NewTicketOrchestrator(
		s.db,
		s.logger,
		s.slaService,
		s.satisfaction,
		s.moduleBus,
		s.moduleRepo.CustomerExists,
		s.moduleRepo.FindAutoAssignableAgent,
		s.buildCustomFieldValues,
		s.prepareTicketCustomFieldMutation,
		s.GetTicketByID,
		s.AssignTicket,
		s.AddComment,
	)
}

func (s *TicketService) ModuleCommandService() *ticketapp.CommandService {
	return s.moduleCmd
}

func (s *TicketService) Orchestrator() *ticketorchestration.TicketOrchestrator {
	return s.orchestrator
}

func (s *TicketService) ApplyAssignTicketSideEffects(ctx context.Context, originalTicket, updatedTicket *models.Ticket) {
	s.orchestrator.ApplyAssignTicketSideEffects(ctx, originalTicket, updatedTicket)
}

func (s *TicketService) ApplyCloseTicketSideEffects(ctx context.Context, ticketID uint, userID uint, reason string) {
	s.orchestrator.ApplyCloseTicketSideEffects(ctx, ticketID, userID, reason)
}

func (s *TicketService) PrepareCreateTicket(ctx context.Context, req *TicketCreateRequest) (*TicketCreatePreparation, error) {
	return s.orchestrator.PrepareCreateTicket(ctx, req)
}

func (s *TicketService) ApplyCreateTicketSideEffects(ctx context.Context, ticket *models.Ticket) (*models.Ticket, error) {
	return s.orchestrator.ApplyCreateTicketSideEffects(ctx, ticket)
}

func (s *TicketService) PrepareUpdateTicket(ctx context.Context, ticketID uint, req *TicketUpdateRequest, userID uint) (*TicketUpdatePreparation, error) {
	return s.orchestrator.PrepareUpdateTicket(ctx, ticketID, req, userID)
}

func (s *TicketService) ApplyUpdateTicketSideEffects(ctx context.Context, prepared *TicketUpdatePreparation, ticketID uint) (*models.Ticket, error) {
	return s.orchestrator.ApplyUpdateTicketSideEffects(ctx, prepared, ticketID)
}

type TicketCreateRequest = ticketcontract.CreateTicketRequest
type TicketUpdateRequest = ticketcontract.UpdateTicketRequest
type TicketListRequest = ticketcontract.ListTicketRequest
type TicketBulkUpdateRequest = ticketcontract.BulkUpdateTicketRequest
type TicketBulkUpdateFailure = ticketcontract.BulkUpdateFailure
type TicketBulkUpdateResult = ticketcontract.BulkUpdateResult

type TicketCreatePreparation = ticketorchestration.TicketCreatePreparation
type TicketUpdatePreparation = ticketorchestration.TicketUpdatePreparation

// CreateTicket 创建工单
func (s *TicketService) CreateTicket(ctx context.Context, req *TicketCreateRequest) (*models.Ticket, error) {
	prepared, err := s.PrepareCreateTicket(ctx, req)
	if err != nil {
		return nil, err
	}
	initialStatus := &models.TicketStatus{
		UserID:     0,
		FromStatus: "",
		ToStatus:   "open",
		Reason:     "工单创建",
	}
	if err := s.moduleRepo.CreateTicketModelWithCustomFieldsAndStatus(ctx, prepared.Ticket, prepared.CustomFieldValues, initialStatus); err != nil {
		return nil, fmt.Errorf("failed to create ticket: %w", err)
	}
	return s.ApplyCreateTicketSideEffects(ctx, prepared.Ticket)
}

// GetTicketByID 根据ID获取工单
func (s *TicketService) GetTicketByID(ctx context.Context, ticketID uint) (*models.Ticket, error) {
	return s.moduleRepo.LoadTicketModelByID(ctx, ticketID)
}

// UpdateTicket 更新工单
func (s *TicketService) UpdateTicket(ctx context.Context, ticketID uint, req *TicketUpdateRequest, userID uint) (*models.Ticket, error) {
	prepared, err := s.PrepareUpdateTicket(ctx, ticketID, req, userID)
	if err != nil {
		return nil, err
	}
	clearAll := false
	deleteFieldIDs := []uint(nil)
	upserts := []models.TicketCustomFieldValue(nil)
	if prepared.Mutation != nil {
		clearAll = prepared.Mutation.ClearAll
		deleteFieldIDs = prepared.Mutation.DeleteFieldIDs
		upserts = ticketapp.MapMutationToModelValues(ticketID, prepared.Mutation)
	}
	if err := s.moduleRepo.UpdateTicketModelWithStatusAndCustomFields(
		ctx,
		ticketID,
		prepared.Updates,
		prepared.StatusChange,
		clearAll,
		deleteFieldIDs,
		upserts,
	); err != nil {
		return nil, fmt.Errorf("failed to update ticket: %w", err)
	}
	return s.ApplyUpdateTicketSideEffects(ctx, prepared, ticketID)
}

// ListTickets 获取工单列表
func (s *TicketService) ListTickets(ctx context.Context, req *TicketListRequest) ([]models.Ticket, int64, error) {
	query := ticketapp.ListTicketsQuery{
		Page:               req.Page,
		PageSize:           req.PageSize,
		Status:             req.Status,
		Priority:           req.Priority,
		Category:           req.Category,
		AgentID:            req.AgentID,
		CustomerID:         req.CustomerID,
		Search:             req.Search,
		SortBy:             ticketcontract.NormalizeTicketSortBy(req.SortBy),
		SortOrder:          req.SortOrder,
		CustomFieldFilters: req.CustomFieldFilters,
	}
	return s.moduleRepo.ListTicketModels(ctx, query)
}

// AssignTicket 分配工单给客服
func (s *TicketService) AssignTicket(ctx context.Context, ticketID uint, agentID uint, assignerID uint) error {
	originalTicket, err := s.GetTicketByID(ctx, ticketID)
	if err != nil {
		return err
	}
	if originalTicket.AgentID != nil && *originalTicket.AgentID == agentID {
		return nil
	}

	if _, err := s.moduleCmd.AssignTicket(ctx, ticketID, ticketapp.AssignTicketCommand{
		AgentID: agentID,
		UserID:  assignerID,
	}); err != nil {
		return err
	}

	s.logger.Infof("Assigned ticket %d to agent %d", ticketID, agentID)
	updatedTicket, err := s.GetTicketByID(ctx, ticketID)
	if err == nil {
		s.ApplyAssignTicketSideEffects(ctx, originalTicket, updatedTicket)
	} else {
		s.logger.Warnf("Failed to fetch ticket %d after assignment for SLA evaluation: %v", ticketID, err)
	}

	return nil
}

// UnassignTicket 取消工单指派（将 agent_id 置空）
func (s *TicketService) UnassignTicket(ctx context.Context, ticketID uint, operatorID uint, reason string) error {
	originalTicket, err := s.GetTicketByID(ctx, ticketID)
	if err != nil {
		return err
	}
	if originalTicket.AgentID == nil {
		return nil
	}
	if reason == "" {
		reason = "取消指派"
	}

	if _, err := s.moduleCmd.UnassignTicket(ctx, ticketID, ticketapp.UnassignTicketCommand{
		UserID: operatorID,
		Reason: reason,
	}); err != nil {
		return err
	}

	updatedTicket, err := s.GetTicketByID(ctx, ticketID)
	if err == nil {
		s.ApplyAssignTicketSideEffects(ctx, originalTicket, updatedTicket)
	} else {
		s.logger.Warnf("Failed to fetch ticket %d after unassignment for SLA evaluation: %v", ticketID, err)
	}
	return nil
}

// BulkUpdateTickets 批量更新工单（支持：状态、标签、指派/取消指派）
func (s *TicketService) BulkUpdateTickets(ctx context.Context, req *TicketBulkUpdateRequest, userID uint) (*TicketBulkUpdateResult, error) {
	result, err := s.moduleCmd.BulkUpdateTickets(ctx, ticketapp.BulkUpdateTicketsCommand{
		TicketIDs:     req.TicketIDs,
		Status:        req.Status,
		SetTags:       req.SetTags,
		AddTags:       req.AddTags,
		RemoveTags:    req.RemoveTags,
		AgentID:       req.AgentID,
		UnassignAgent: req.UnassignAgent,
		UserID:        userID,
	})
	if err != nil {
		return nil, err
	}
	out := &TicketBulkUpdateResult{
		Updated: result.Updated,
		Failed:  make([]TicketBulkUpdateFailure, 0, len(result.Failed)),
	}
	for _, failure := range result.Failed {
		out.Failed = append(out.Failed, TicketBulkUpdateFailure{
			TicketID: failure.TicketID,
			Error:    failure.Error,
		})
	}
	return out, nil
}

// AddComment 添加工单评论
func (s *TicketService) AddComment(ctx context.Context, ticketID uint, userID uint, content string, commentType string) (*models.TicketComment, error) {
	comment, err := s.moduleCmd.AddComment(ctx, ticketID, ticketapp.AddCommentCommand{
		UserID:      userID,
		Content:     content,
		CommentType: commentType,
	})
	if err != nil {
		return nil, err
	}
	return &models.TicketComment{
		ID:        comment.ID,
		TicketID:  ticketID,
		UserID:    comment.UserID,
		Content:   comment.Content,
		Type:      comment.Type,
		CreatedAt: comment.CreatedAt,
	}, nil
}

// CloseTicket 关闭工单
func (s *TicketService) CloseTicket(ctx context.Context, ticketID uint, userID uint, reason string) error {
	if _, err := s.moduleCmd.CloseTicket(ctx, ticketID, ticketapp.CloseTicketCommand{
		UserID: userID,
		Reason: reason,
	}); err != nil {
		return err
	}
	s.ApplyCloseTicketSideEffects(ctx, ticketID, userID, reason)
	return nil
}

// GetTicketStats 获取工单统计
func (s *TicketService) GetTicketStats(ctx context.Context, agentID *uint) (*TicketStats, error) {
	stats, err := s.moduleQuery.GetTicketStats(ctx, agentID)
	if err != nil {
		return nil, err
	}
	return ticketcontract.MapTicketStats(stats), nil
}

type TicketStats = ticketcontract.TicketStats
type StatusCount = ticketcontract.StatusCount
type PriorityCount = ticketcontract.PriorityCount

type ticketCustomFieldMutation = ticketapp.CustomFieldMutation

// ListTicketCustomFields returns ticket custom field definitions.
func (s *TicketService) ListTicketCustomFields(ctx context.Context, activeOnly bool) ([]models.CustomField, error) {
	return s.moduleRepo.ListTicketCustomFields(ctx, activeOnly)
}

func (s *TicketService) buildCustomFieldValues(ctx context.Context, provided map[string]interface{}, ticketCtx map[string]interface{}, enforceRequired bool) ([]models.TicketCustomFieldValue, error) {
	fields, err := s.moduleRepo.ListTicketCustomFields(ctx, true)
	if err != nil {
		return nil, err
	}
	return ticketapp.BuildModelCustomFieldValues(fields, provided, ticketCtx, enforceRequired)
}

func (s *TicketService) prepareTicketCustomFieldMutation(ctx context.Context, ticketID uint, provided map[string]interface{}, ticketCtx map[string]interface{}) (*ticketCustomFieldMutation, error) {
	fields, err := s.moduleRepo.ListTicketCustomFields(ctx, false)
	if err != nil {
		return nil, err
	}
	return ticketapp.PrepareCustomFieldMutation(fields, ticketID, provided, ticketCtx)
}
