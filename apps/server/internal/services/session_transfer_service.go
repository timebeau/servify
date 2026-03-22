package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"servify/apps/server/internal/models"
	agentdelivery "servify/apps/server/internal/modules/agent/delivery"
	conversationdelivery "servify/apps/server/internal/modules/conversation/delivery"
	routingcontract "servify/apps/server/internal/modules/routing/contract"
	routingdelivery "servify/apps/server/internal/modules/routing/delivery"
	ticketdelivery "servify/apps/server/internal/modules/ticket/delivery"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// SessionTransferService 会话转接服务
type SessionTransferService struct {
	db           *gorm.DB
	logger       *logrus.Logger
	aiService    AIServiceInterface
	agentService sessionTransferAgentRuntime
	notifier     sessionTransferNotifier
	routing      routingdelivery.RuntimeService
	tickets      ticketdelivery.RuntimeService
	conversation conversationdelivery.RuntimeService
	agents       agentdelivery.RuntimeService
}

type sessionTransferAgentRuntime interface {
	FindAvailableAgent(ctx context.Context, skills []string, priority string) (*AgentInfo, error)
	GetOnlineAgent(userID uint) (*AgentInfo, bool)
	ApplySessionTransfer(sessionID string, fromAgentID *uint, toAgentID uint)
}

type SessionTransferRuntime interface {
	TransferToHuman(ctx context.Context, req *TransferRequest) (*TransferResult, error)
}

type sessionTransferNotifier interface {
	SendToSession(sessionID string, message WebSocketMessage)
}

type SessionTransferAdapters struct {
	Routing      routingdelivery.RuntimeService
	Tickets      ticketdelivery.RuntimeService
	Conversation conversationdelivery.RuntimeService
	Agents       agentdelivery.RuntimeService
}

// NewSessionTransferService 创建会话转接服务
func NewSessionTransferService(
	db *gorm.DB,
	logger *logrus.Logger,
	aiService AIServiceInterface,
	agentService sessionTransferAgentRuntime,
	notifier sessionTransferNotifier,
) *SessionTransferService {
	if logger == nil {
		logger = logrus.New()
	}

	return &SessionTransferService{
		db:           db,
		logger:       logger,
		aiService:    aiService,
		agentService: agentService,
		notifier:     notifier,
	}
}

func NewSessionTransferServiceWithAdapters(
	db *gorm.DB,
	logger *logrus.Logger,
	aiService AIServiceInterface,
	agentService sessionTransferAgentRuntime,
	notifier sessionTransferNotifier,
	adapters SessionTransferAdapters,
) *SessionTransferService {
	svc := NewSessionTransferService(db, logger, aiService, agentService, notifier)
	svc.routing = adapters.Routing
	svc.tickets = adapters.Tickets
	svc.conversation = adapters.Conversation
	svc.agents = adapters.Agents
	return svc
}

type TransferRequest = routingcontract.TransferRequest
type TransferResult = routingcontract.TransferResult

// TransferToHuman 转接到人工客服
func (s *SessionTransferService) TransferToHuman(ctx context.Context, req *TransferRequest) (*TransferResult, error) {
	session, err := s.loadTransferSession(ctx, req.SessionID)
	if err != nil {
		return nil, fmt.Errorf("session not found: %w", err)
	}

	// 检查会话状态
	if session.Status == "ended" {
		return nil, fmt.Errorf("session already ended")
	}
	if session.AgentID != nil {
		return nil, fmt.Errorf("session already assigned")
	}

	// 若已在等待队列中，直接返回（避免重复入队）
	if existing, err := s.getActiveWaitingRecord(ctx, session.ID); err == nil {
		return &TransferResult{
			Success:   true,
			SessionID: session.ID,
			IsWaiting: true,
			QueuedAt:  &existing.QueuedAt,
			Summary:   "会话已在等待队列中",
		}, nil
	}

	// 查找可用的客服
	agent, err := s.agentService.FindAvailableAgent(ctx, req.TargetSkills, req.Priority)
	if err != nil {
		// 没有可用客服，加入等待队列
		return s.addToWaitingQueue(ctx, session, req)
	}

	// 执行转接
	return s.executeTransfer(ctx, session, agent.UserID, req.Reason, req.Notes)
}

// TransferToAgent 转接到指定客服
func (s *SessionTransferService) TransferToAgent(ctx context.Context, sessionID string, targetAgentID uint, reason string) (*TransferResult, error) {
	// 获取会话信息
	session, err := s.loadTransferSession(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("session not found: %w", err)
	}

	// 检查目标客服是否可用
	agentInfo, ok := s.agentService.GetOnlineAgent(targetAgentID)
	if !ok {
		return nil, fmt.Errorf("target agent is not online")
	}

	if agentInfo.CurrentLoad >= agentInfo.MaxConcurrent {
		return nil, fmt.Errorf("target agent is at maximum capacity")
	}

	// 执行转接
	return s.executeTransfer(ctx, session, targetAgentID, reason, "")
}

// executeTransfer 执行转接
func (s *SessionTransferService) executeTransfer(ctx context.Context, session *conversationdelivery.TransferSession, targetAgentID uint, reason, notes string) (*TransferResult, error) {
	if session.Status == "ended" {
		return nil, fmt.Errorf("session already ended")
	}
	if session.AgentID != nil && *session.AgentID == targetAgentID {
		return &TransferResult{
			Success:    true,
			SessionID:  session.ID,
			NewAgentID: targetAgentID,
			Summary:    "会话已指派给目标客服",
		}, nil
	}

	fromAgentID := session.AgentID
	transferAt := time.Now()

	// 生成会话摘要（在事务外，避免长事务）
	summary, err := s.generateSessionSummary(session)
	if err != nil {
		s.logger.Warnf("Failed to generate session summary: %v", err)
		summary = "无法生成会话摘要"
	}

	transferMessageContent := s.buildTransferMessage(reason, notes)

	transferRecord := &models.TransferRecord{
		SessionID:      session.ID,
		FromAgentID:    fromAgentID,
		ToAgentID:      &targetAgentID,
		Reason:         reason,
		Notes:          notes,
		SessionSummary: summary,
		TransferredAt:  transferAt,
	}

	// 原子化：会话指派 + 工时负载 + 记录/消息
	if err := s.db.Transaction(func(tx *gorm.DB) error {
		if err := s.syncTransferSession(ctx, tx, session, targetAgentID); err != nil {
			return fmt.Errorf("update session: %w", err)
		}

		if err := s.syncTransferTicket(ctx, tx, session, targetAgentID, transferAt); err != nil {
			return fmt.Errorf("update ticket: %w", err)
		}

		if err := s.syncTransferAgentLoad(ctx, tx, fromAgentID, targetAgentID); err != nil {
			return fmt.Errorf("sync agent load: %w", err)
		}

		if err := s.appendTransferSystemMessage(ctx, tx, session.ID, targetAgentID, transferMessageContent, transferAt); err != nil {
			return fmt.Errorf("create transfer message: %w", err)
		}

		createdRecord, err := s.createTransferRecord(ctx, tx, routingdelivery.AssignAgentCommand{
			SessionID:      session.ID,
			AgentID:        targetAgentID,
			FromAgentID:    fromAgentID,
			Reason:         reason,
			Notes:          notes,
			SessionSummary: summary,
			AssignedAt:     transferAt,
		}, transferRecord)
		if err != nil {
			return fmt.Errorf("create transfer record: %w", err)
		}
		transferRecord = createdRecord

		if err := s.markWaitingTransferred(ctx, tx, session.ID, targetAgentID, transferAt); err != nil {
			return fmt.Errorf("sync waiting record: %w", err)
		}

		return nil
	}); err != nil {
		return nil, err
	}

	// 更新内存状态（在线客服负载/会话映射），避免与 DB 状态长期漂移
	if s.agentService != nil {
		s.agentService.ApplySessionTransfer(session.ID, fromAgentID, targetAgentID)
	}

	// 发送实时通知
	s.notifyTransfer(session.ID, targetAgentID, transferMessageContent)

	s.logger.Infof("Successfully transferred session %s to agent %d", session.ID, targetAgentID)

	return &TransferResult{
		Success:       true,
		SessionID:     session.ID,
		NewAgentID:    targetAgentID,
		TransferredAt: transferRecord.TransferredAt,
		Summary:       summary,
	}, nil
}

// addToWaitingQueue 添加到等待队列
func (s *SessionTransferService) addToWaitingQueue(ctx context.Context, session *conversationdelivery.TransferSession, req *TransferRequest) (*TransferResult, error) {
	// 若已在等待队列中，直接返回（避免重复入队）
	if existing, err := s.getActiveWaitingRecord(ctx, session.ID); err == nil {
		return &TransferResult{
			Success:   true,
			SessionID: session.ID,
			IsWaiting: true,
			QueuedAt:  &existing.QueuedAt,
			Summary:   "会话已在等待队列中",
		}, nil
	}

	var waitingRecord *models.WaitingRecord
	waitingMessage := &models.Message{
		SessionID: session.ID,
		UserID:    session.CustomerID,
		Content:   "您的会话已加入人工客服等待队列，我们会尽快为您安排客服。请耐心等待。",
		Type:      "system",
		Sender:    "system",
	}
	if err := s.db.Transaction(func(tx *gorm.DB) error {
		if err := s.ensureWaitingSessionState(ctx, tx, session); err != nil {
			return fmt.Errorf("failed to ensure session active: %w", err)
		}

		createdRecord, err := s.createWaitingRecord(ctx, tx, session.ID, req)
		if err != nil {
			return fmt.Errorf("failed to create waiting record: %w", err)
		}
		waitingRecord = createdRecord

		if err := s.appendWaitingSystemMessage(ctx, tx, waitingMessage); err != nil {
			return fmt.Errorf("create waiting message: %w", err)
		}
		return nil
	}); err != nil {
		return nil, err
	}

	// 发送实时通知
	s.notifyWaiting(session.ID, waitingMessage.Content)

	s.logger.Infof("Added session %s to waiting queue", session.ID)

	return &TransferResult{
		Success:   true,
		SessionID: session.ID,
		IsWaiting: true,
		QueuedAt:  &waitingRecord.QueuedAt,
		Summary:   "会话已加入等待队列",
	}, nil
}

// ProcessWaitingQueue 处理等待队列
func (s *SessionTransferService) ProcessWaitingQueue(ctx context.Context) error {
	waitingRecords, err := s.loadWaitingQueue(ctx, 10)
	if err != nil {
		return err
	}

	for _, record := range waitingRecords {
		// 查找可用客服
		skills := []string{}
		if record.TargetSkills != "" {
			skills = strings.Split(record.TargetSkills, ",")
		}

		agent, err := s.agentService.FindAvailableAgent(ctx, skills, record.Priority)
		if err != nil {
			continue // 没有可用客服，继续下一个
		}

		// 获取会话信息
		session, err := s.loadTransferSession(ctx, record.SessionID)
		if err != nil {
			continue
		}

		// 执行转接
		result, err := s.executeTransfer(ctx, session, agent.UserID, record.Reason, record.Notes)
		if err != nil {
			s.logger.Errorf("Failed to transfer waiting session %s: %v", record.SessionID, err)
			continue
		}

		s.logger.Infof("Successfully transferred waiting session %s to agent %d",
			result.SessionID, result.NewAgentID)
	}

	return nil
}

func (s *SessionTransferService) loadWaitingQueue(ctx context.Context, limit int) ([]models.WaitingRecord, error) {
	return s.listWaitingRecords(ctx, "waiting", limit)
}

func normalizeWaitingRecordQuery(status string, limit int) (string, int) {
	if status == "" {
		status = "waiting"
	}
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	return status, limit
}

func (s *SessionTransferService) listWaitingRecords(ctx context.Context, status string, limit int) ([]models.WaitingRecord, error) {
	status, limit = normalizeWaitingRecordQuery(status, limit)
	if s.routing != nil {
		records, err := s.routing.ListWaitingRecords(ctx, status, limit)
		if err != nil {
			return nil, fmt.Errorf("failed to list waiting records: %w", err)
		}
		return records, nil
	}

	var waitingRecords []models.WaitingRecord
	if err := s.db.Where("status = ?", status).
		Order("priority DESC, queued_at ASC").
		Limit(limit).
		Find(&waitingRecords).Error; err != nil {
		return nil, fmt.Errorf("failed to list waiting records: %w", err)
	}
	return waitingRecords, nil
}

func (s *SessionTransferService) listTransferHistory(ctx context.Context, sessionID string) ([]models.TransferRecord, error) {
	if s.routing != nil {
		records, err := s.routing.GetTransferHistory(ctx, sessionID)
		if err != nil {
			return nil, fmt.Errorf("failed to get transfer history: %w", err)
		}
		return records, nil
	}

	var records []models.TransferRecord
	if err := s.db.Where("session_id = ?", sessionID).
		Order("transferred_at DESC").
		Find(&records).Error; err != nil {
		return nil, fmt.Errorf("failed to get transfer history: %w", err)
	}

	return records, nil
}

func (s *SessionTransferService) loadWaitingRecord(ctx context.Context, sessionID string) (*models.WaitingRecord, error) {
	if s.routing != nil {
		record, err := s.routing.GetWaitingRecord(ctx, sessionID)
		if err != nil {
			return nil, err
		}
		return record, nil
	}

	var record models.WaitingRecord
	if err := s.db.Where("session_id = ?", sessionID).Order("queued_at DESC").First(&record).Error; err != nil {
		return nil, err
	}
	return &record, nil
}

func (s *SessionTransferService) syncTransferSession(ctx context.Context, tx *gorm.DB, session *conversationdelivery.TransferSession, targetAgentID uint) error {
	if s.conversation != nil {
		return s.conversation.SyncTransferAssignment(ctx, tx, session.ID, session.CustomerID, targetAgentID)
	}
	return tx.Model(&models.Session{}).
		Where("id = ?", session.ID).
		Updates(map[string]interface{}{
			"agent_id": targetAgentID,
			"status":   "active",
			"ended_at": nil,
		}).Error
}

func (s *SessionTransferService) syncTransferTicket(ctx context.Context, tx *gorm.DB, session *conversationdelivery.TransferSession, targetAgentID uint, transferAt time.Time) error {
	if session.TicketID == nil || *session.TicketID == 0 {
		return nil
	}
	if s.tickets != nil {
		err := s.tickets.SyncTransferAssignment(ctx, tx, *session.TicketID, targetAgentID, targetAgentID)
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		return err
	}

	ticket, err := s.loadTransferTicket(tx, *session.TicketID)
	if err == gorm.ErrRecordNotFound {
		return nil
	}
	if err != nil {
		return err
	}

	updates, fromStatus, toStatus := buildTransferTicketUpdate(targetAgentID, ticket.Status)
	if err := tx.Model(&models.Ticket{}).Where("id = ?", ticket.ID).Updates(updates).Error; err != nil {
		return err
	}

	s.appendTransferTicketStatus(tx, ticket.ID, targetAgentID, fromStatus, toStatus, transferAt)
	return nil
}

func (s *SessionTransferService) loadTransferTicket(tx *gorm.DB, ticketID uint) (*models.Ticket, error) {
	var ticket models.Ticket
	if err := tx.Select("id", "agent_id", "status").First(&ticket, "id = ?", ticketID).Error; err != nil {
		return nil, err
	}
	return &ticket, nil
}

func buildTransferTicketUpdate(targetAgentID uint, currentStatus string) (map[string]interface{}, string, string) {
	updates := map[string]interface{}{"agent_id": targetAgentID}
	fromStatus := currentStatus
	toStatus := fromStatus
	if fromStatus == "open" || fromStatus == "" {
		toStatus = "assigned"
		updates["status"] = toStatus
	}
	return updates, fromStatus, toStatus
}

func (s *SessionTransferService) appendTransferTicketStatus(tx *gorm.DB, ticketID uint, targetAgentID uint, fromStatus string, toStatus string, transferAt time.Time) {
	_ = tx.Create(&models.TicketStatus{
		TicketID:   ticketID,
		UserID:     targetAgentID,
		FromStatus: fromStatus,
		ToStatus:   toStatus,
		Reason:     fmt.Sprintf("会话转接同步指派至客服 %d", targetAgentID),
		CreatedAt:  transferAt,
	}).Error
}

func (s *SessionTransferService) syncTransferAgentLoad(ctx context.Context, tx *gorm.DB, fromAgentID *uint, targetAgentID uint) error {
	if s.agents != nil {
		return s.agents.SyncTransferLoad(ctx, tx, fromAgentID, targetAgentID)
	}
	if fromAgentID != nil && *fromAgentID != targetAgentID {
		if err := tx.Exec(`UPDATE agents SET current_load = CASE WHEN current_load > 0 THEN current_load - 1 ELSE 0 END WHERE user_id = ?`, *fromAgentID).Error; err != nil {
			return err
		}
	}
	return tx.Exec(`UPDATE agents SET current_load = current_load + 1 WHERE user_id = ?`, targetAgentID).Error
}

func (s *SessionTransferService) appendTransferSystemMessage(ctx context.Context, tx *gorm.DB, sessionID string, targetAgentID uint, content string, createdAt time.Time) error {
	if s.conversation != nil {
		return s.conversation.AppendSystemMessage(ctx, tx, sessionID, content, createdAt)
	}
	return tx.Create(&models.Message{
		SessionID: sessionID,
		UserID:    targetAgentID,
		Content:   content,
		Type:      "system",
		Sender:    "system",
		CreatedAt: createdAt,
	}).Error
}

func (s *SessionTransferService) createTransferRecord(ctx context.Context, tx *gorm.DB, cmd routingdelivery.AssignAgentCommand, fallback *models.TransferRecord) (*models.TransferRecord, error) {
	if s.routing != nil {
		return s.routing.AssignAgent(ctx, tx, cmd)
	}
	if err := tx.Create(fallback).Error; err != nil {
		return nil, err
	}
	return fallback, nil
}

func (s *SessionTransferService) markWaitingTransferred(ctx context.Context, tx *gorm.DB, sessionID string, targetAgentID uint, transferAt time.Time) error {
	if s.routing != nil {
		_, err := s.routing.MarkWaitingTransferred(ctx, tx, sessionID, targetAgentID, transferAt)
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		return err
	}
	return tx.Model(&models.WaitingRecord{}).
		Where("session_id = ? AND status = ?", sessionID, "waiting").
		Updates(map[string]interface{}{
			"status":      "transferred",
			"assigned_at": transferAt,
			"assigned_to": targetAgentID,
		}).Error
}

func (s *SessionTransferService) cancelWaitingRecord(ctx context.Context, tx *gorm.DB, sessionID string, reason string) error {
	if s.routing != nil {
		_, err := s.routing.CancelWaiting(ctx, tx, sessionID, reason)
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		if err != nil {
			return fmt.Errorf("update waiting record: %w", err)
		}
		return nil
	}

	record, err := s.getActiveWaitingRecord(ctx, sessionID)
	if err == gorm.ErrRecordNotFound {
		return nil
	}
	if err != nil {
		return fmt.Errorf("load waiting record: %w", err)
	}

	if err := tx.Model(&models.WaitingRecord{}).
		Where("session_id = ? AND status = ?", record.SessionID, "waiting").
		Updates(map[string]interface{}{
			"status": "cancelled",
		}).Error; err != nil {
		return fmt.Errorf("update waiting record: %w", err)
	}
	return nil
}

func (s *SessionTransferService) ensureWaitingSessionState(ctx context.Context, tx *gorm.DB, session *conversationdelivery.TransferSession) error {
	if s.conversation != nil {
		return s.conversation.SyncWaitingAssignment(ctx, tx, session.ID, session.CustomerID)
	}
	return tx.Model(&models.Session{}).
		Where("id = ?", session.ID).
		Updates(map[string]interface{}{
			"status":   "active",
			"agent_id": nil,
			"ended_at": nil,
		}).Error
}

func (s *SessionTransferService) createWaitingRecord(ctx context.Context, tx *gorm.DB, sessionID string, req *TransferRequest) (*models.WaitingRecord, error) {
	if s.routing != nil {
		return s.routing.AddToWaitingQueue(ctx, tx, sessionID, req.Reason, req.TargetSkills, req.Priority, req.Notes)
	}
	record := &models.WaitingRecord{
		SessionID:    sessionID,
		Reason:       req.Reason,
		TargetSkills: strings.Join(req.TargetSkills, ","),
		Priority:     req.Priority,
		Notes:        req.Notes,
		Status:       "waiting",
		QueuedAt:     time.Now(),
	}
	if err := tx.Create(record).Error; err != nil {
		return nil, err
	}
	return record, nil
}

func (s *SessionTransferService) appendWaitingSystemMessage(ctx context.Context, tx *gorm.DB, message *models.Message) error {
	if s.conversation != nil {
		return s.conversation.AppendSystemMessage(ctx, tx, message.SessionID, message.Content, message.CreatedAt)
	}
	return tx.Create(message).Error
}

func (s *SessionTransferService) appendCancellationSystemMessage(ctx context.Context, tx *gorm.DB, sessionID string, operatorID uint, reason string, createdAt time.Time) error {
	content := fmt.Sprintf("已取消人工客服等待队列（原因：%s）", reason)
	if s.conversation != nil {
		return s.conversation.AppendSystemMessage(ctx, tx, sessionID, content, createdAt)
	}
	return tx.Create(&models.Message{
		SessionID: sessionID,
		UserID:    operatorID,
		Content:   content,
		Type:      "system",
		Sender:    "system",
		CreatedAt: createdAt,
	}).Error
}

func (s *SessionTransferService) getActiveWaitingRecord(ctx context.Context, sessionID string) (*models.WaitingRecord, error) {
	record, err := s.loadWaitingRecord(ctx, sessionID)
	if err != nil {
		return nil, err
	}
	if record.Status != "waiting" {
		return nil, gorm.ErrRecordNotFound
	}
	return record, nil
}

// generateSessionSummary 生成会话摘要
func (s *SessionTransferService) generateSessionSummary(session *conversationdelivery.TransferSession) (string, error) {
	// 获取会话消息
	var messages []models.Message
	if err := s.db.Where("session_id = ?", session.ID).
		Order("created_at ASC").
		Find(&messages).Error; err != nil {
		return "", err
	}

	// 如果消息太少，返回简单摘要
	if len(messages) < 3 {
		userLabel := session.UserUsername
		if userLabel == "" {
			userLabel = session.UserName
		}
		if userLabel == "" {
			userLabel = fmt.Sprintf("ID=%d", session.CustomerID)
		}
		return fmt.Sprintf("用户%s的简短会话，共%d条消息", userLabel, len(messages)), nil
	}

	// 使用 AI 服务生成摘要
	return s.aiService.GetSessionSummary(messages)
}

// buildTransferMessage 构建转接消息
func (s *SessionTransferService) buildTransferMessage(reason, notes string) string {
	message := "您的会话已转接至人工客服"

	if reason != "" {
		message += fmt.Sprintf("。转接原因：%s", reason)
	}

	if notes != "" {
		message += fmt.Sprintf("。备注：%s", notes)
	}

	message += "。客服将很快为您提供帮助。"

	return message
}

// notifyTransfer 发送转接通知
func (s *SessionTransferService) notifyTransfer(sessionID string, agentID uint, message string) {
	// 发送给用户
	if s.notifier != nil {
		s.notifier.SendToSession(sessionID, WebSocketMessage{
			Type: "transfer_notification",
			Data: map[string]interface{}{
				"message":   message,
				"agent_id":  agentID,
				"timestamp": time.Now(),
			},
		})
	}
}

// notifyWaiting 发送等待通知
func (s *SessionTransferService) notifyWaiting(sessionID string, message string) {
	if s.notifier != nil {
		s.notifier.SendToSession(sessionID, WebSocketMessage{
			Type: "waiting_notification",
			Data: map[string]interface{}{
				"message":   message,
				"timestamp": time.Now(),
			},
		})
	}
}

// GetTransferHistory 获取转接历史
func (s *SessionTransferService) GetTransferHistory(ctx context.Context, sessionID string) ([]models.TransferRecord, error) {
	return s.listTransferHistory(ctx, sessionID)
}

// ListWaitingRecords 列出等待队列记录（默认 status=waiting）
func (s *SessionTransferService) ListWaitingRecords(ctx context.Context, status string, limit int) ([]models.WaitingRecord, error) {
	return s.listWaitingRecords(ctx, status, limit)
}

// CancelWaitingRecord 取消等待队列中的会话（幂等）
func (s *SessionTransferService) CancelWaitingRecord(ctx context.Context, sessionID string, operatorID uint, reason string) error {
	if sessionID == "" {
		return fmt.Errorf("session_id is required")
	}
	if reason == "" {
		reason = "cancelled"
	}

	now := time.Now()
	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := s.cancelWaitingRecord(ctx, tx, sessionID, reason); err != nil {
			return err
		}

		// 记录系统消息（可选，不影响主流程）
		_ = s.appendCancellationSystemMessage(ctx, tx, sessionID, operatorID, reason, now)
		return nil
	})
}

func (s *SessionTransferService) loadTransferSession(ctx context.Context, sessionID string) (*conversationdelivery.TransferSession, error) {
	if s.conversation != nil {
		return s.conversation.LoadTransferSession(ctx, sessionID)
	}

	var model models.Session
	if err := s.db.WithContext(ctx).Preload("User").First(&model, "id = ?", sessionID).Error; err != nil {
		return nil, err
	}
	return &conversationdelivery.TransferSession{
		ID:           model.ID,
		CustomerID:   model.UserID,
		AgentID:      model.AgentID,
		TicketID:     model.TicketID,
		Status:       model.Status,
		Platform:     model.Platform,
		UserName:     model.User.Name,
		UserUsername: model.User.Username,
	}, nil
}

// AutoTransferCheck 自动转接检查
func (s *SessionTransferService) AutoTransferCheck(ctx context.Context, sessionID string, messages []models.Message) bool {
	// 检查是否需要自动转接
	lastMessages := messages
	if len(messages) > 5 {
		lastMessages = messages[len(messages)-5:]
	}

	// 构建查询字符串
	var queryBuilder strings.Builder
	for _, msg := range lastMessages {
		if msg.Sender == "user" {
			queryBuilder.WriteString(msg.Content)
			queryBuilder.WriteString(" ")
		}
	}

	query := queryBuilder.String()

	// 使用 AI 服务判断是否需要转人工
	return s.aiService.ShouldTransferToHuman(query, messages)
}
