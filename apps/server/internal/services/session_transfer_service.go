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
	agentService *AgentService
	wsHub        *WebSocketHub
	routing      routingdelivery.RuntimeService
	tickets      ticketdelivery.RuntimeService
	conversation conversationdelivery.RuntimeService
	agents       agentdelivery.RuntimeService
}

// NewSessionTransferService 创建会话转接服务
func NewSessionTransferService(
	db *gorm.DB,
	logger *logrus.Logger,
	aiService AIServiceInterface,
	agentService *AgentService,
	wsHub *WebSocketHub,
) *SessionTransferService {
	if logger == nil {
		logger = logrus.New()
	}

	return &SessionTransferService{
		db:           db,
		logger:       logger,
		aiService:    aiService,
		agentService: agentService,
		wsHub:        wsHub,
	}
}

type TransferRequest = routingcontract.TransferRequest
type TransferResult = routingcontract.TransferResult

// TransferToHuman 转接到人工客服
func (s *SessionTransferService) TransferToHuman(ctx context.Context, req *TransferRequest) (*TransferResult, error) {
	// 获取会话信息
	var session models.Session
	if err := s.db.Preload("User").Preload("Messages").First(&session, "id = ?", req.SessionID).Error; err != nil {
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
		return s.addToWaitingQueue(ctx, &session, req)
	}

	// 执行转接
	return s.executeTransfer(ctx, &session, agent.UserID, req.Reason, req.Notes)
}

// TransferToAgent 转接到指定客服
func (s *SessionTransferService) TransferToAgent(ctx context.Context, sessionID string, targetAgentID uint, reason string) (*TransferResult, error) {
	// 获取会话信息
	var session models.Session
	if err := s.db.Preload("User").First(&session, "id = ?", sessionID).Error; err != nil {
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
	return s.executeTransfer(ctx, &session, targetAgentID, reason, "")
}

// executeTransfer 执行转接
func (s *SessionTransferService) executeTransfer(ctx context.Context, session *models.Session, targetAgentID uint, reason, notes string) (*TransferResult, error) {
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
		// 更新会话：active/ended；是否分配由 agent_id 判断
		if s.conversation != nil {
			if err := s.conversation.SyncTransferAssignment(ctx, tx, session.ID, session.UserID, targetAgentID); err != nil {
				return fmt.Errorf("update session: %w", err)
			}
		} else if err := tx.Model(&models.Session{}).
			Where("id = ?", session.ID).
			Updates(map[string]interface{}{
				"agent_id": targetAgentID,
				"status":   "active",
				"ended_at": nil,
			}).Error; err != nil {
			return fmt.Errorf("update session: %w", err)
		}

		// 若会话关联了工单，则同步工单的指派（与会话转接保持一致）
		if session.TicketID != nil && *session.TicketID != 0 {
			if s.tickets != nil {
				if err := s.tickets.SyncTransferAssignment(ctx, tx, *session.TicketID, targetAgentID, targetAgentID); err != nil && err != gorm.ErrRecordNotFound {
					return fmt.Errorf("update ticket: %w", err)
				}
			} else {
				var ticket models.Ticket
				if err := tx.Select("id", "agent_id", "status").First(&ticket, "id = ?", *session.TicketID).Error; err != nil {
					if err != gorm.ErrRecordNotFound {
						return fmt.Errorf("load ticket: %w", err)
					}
				} else {
					updates := map[string]interface{}{
						"agent_id": targetAgentID,
					}
					fromStatus := ticket.Status
					toStatus := fromStatus
					if fromStatus == "open" || fromStatus == "" {
						toStatus = "assigned"
						updates["status"] = toStatus
					}
					if err := tx.Model(&models.Ticket{}).Where("id = ?", ticket.ID).Updates(updates).Error; err != nil {
						return fmt.Errorf("update ticket: %w", err)
					}

					// Best-effort 记录状态变更（不影响主流程）
					_ = tx.Create(&models.TicketStatus{
						TicketID:   ticket.ID,
						UserID:     targetAgentID,
						FromStatus: fromStatus,
						ToStatus:   toStatus,
						Reason:     fmt.Sprintf("会话转接同步指派至客服 %d", targetAgentID),
						CreatedAt:  transferAt,
					}).Error
				}
			}
		}

		// 负载：转移需要先减后加（最佳努力不低于 0）
		if s.agents != nil {
			if err := s.agents.SyncTransferLoad(ctx, tx, fromAgentID, targetAgentID); err != nil {
				return fmt.Errorf("sync agent load: %w", err)
			}
		} else {
			if fromAgentID != nil && *fromAgentID != targetAgentID {
				if err := tx.Exec(`UPDATE agents SET current_load = CASE WHEN current_load > 0 THEN current_load - 1 ELSE 0 END WHERE user_id = ?`, *fromAgentID).Error; err != nil {
					return fmt.Errorf("decrement from agent load: %w", err)
				}
			}
			if err := tx.Exec(`UPDATE agents SET current_load = current_load + 1 WHERE user_id = ?`, targetAgentID).Error; err != nil {
				return fmt.Errorf("increment target agent load: %w", err)
			}
		}

		if s.conversation != nil {
			if err := s.conversation.AppendSystemMessage(ctx, tx, session.ID, transferMessageContent, transferAt); err != nil {
				return fmt.Errorf("create transfer message: %w", err)
			}
		} else {
			transferMessage := &models.Message{
				SessionID: session.ID,
				UserID:    targetAgentID,
				Content:   transferMessageContent,
				Type:      "system",
				Sender:    "system",
				CreatedAt: transferAt,
			}
			if err := tx.Create(transferMessage).Error; err != nil {
				return fmt.Errorf("create transfer message: %w", err)
			}
		}

		if s.routing != nil {
			createdRecord, err := s.routing.AssignAgent(ctx, tx, routingdelivery.AssignAgentCommand{
				SessionID:      session.ID,
				AgentID:        targetAgentID,
				FromAgentID:    fromAgentID,
				Reason:         reason,
				Notes:          notes,
				SessionSummary: summary,
				AssignedAt:     transferAt,
			})
			if err != nil {
				return fmt.Errorf("create transfer record: %w", err)
			}
			transferRecord = createdRecord
		} else if err := tx.Create(transferRecord).Error; err != nil {
			return fmt.Errorf("create transfer record: %w", err)
		}

		if s.routing != nil {
			if _, err := s.routing.MarkWaitingTransferred(ctx, tx, session.ID, targetAgentID, transferAt); err != nil && err != gorm.ErrRecordNotFound {
				return fmt.Errorf("sync waiting record: %w", err)
			}
		} else if err := tx.Model(&models.WaitingRecord{}).
			Where("session_id = ? AND status = ?", session.ID, "waiting").
			Updates(map[string]interface{}{
				"status":      "transferred",
				"assigned_at": transferAt,
				"assigned_to": targetAgentID,
			}).Error; err != nil {
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
func (s *SessionTransferService) addToWaitingQueue(ctx context.Context, session *models.Session, req *TransferRequest) (*TransferResult, error) {
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
		UserID:    session.UserID,
		Content:   "您的会话已加入人工客服等待队列，我们会尽快为您安排客服。请耐心等待。",
		Type:      "system",
		Sender:    "system",
	}
	if err := s.db.Transaction(func(tx *gorm.DB) error {
		// 会话保持 active，等待队列由 WaitingRecord 表达
		if err := tx.Model(&models.Session{}).
			Where("id = ?", session.ID).
			Updates(map[string]interface{}{
				"status":   "active",
				"agent_id": nil,
				"ended_at": nil,
			}).Error; err != nil {
			return fmt.Errorf("failed to ensure session active: %w", err)
		}

		if s.routing != nil {
			createdRecord, err := s.routing.AddToWaitingQueue(ctx, tx, session.ID, req.Reason, req.TargetSkills, req.Priority, req.Notes)
			if err != nil {
				return fmt.Errorf("failed to create waiting record: %w", err)
			}
			waitingRecord = createdRecord
		} else {
			waitingRecord = &models.WaitingRecord{
				SessionID:    session.ID,
				Reason:       req.Reason,
				TargetSkills: strings.Join(req.TargetSkills, ","),
				Priority:     req.Priority,
				Notes:        req.Notes,
				Status:       "waiting",
				QueuedAt:     time.Now(),
			}

			if err := tx.Create(waitingRecord).Error; err != nil {
				return fmt.Errorf("failed to create waiting record: %w", err)
			}
		}

		if s.conversation != nil {
			if err := s.conversation.AppendSystemMessage(ctx, tx, session.ID, waitingMessage.Content, waitingMessage.CreatedAt); err != nil {
				return fmt.Errorf("create waiting message: %w", err)
			}
		} else if err := tx.Create(waitingMessage).Error; err != nil {
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
		var session models.Session
		if err := s.db.First(&session, "id = ?", record.SessionID).Error; err != nil {
			continue
		}

		// 执行转接
		result, err := s.executeTransfer(ctx, &session, agent.UserID, record.Reason, record.Notes)
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
	if s.routing != nil {
		records, err := s.routing.ListWaitingRecords(ctx, "waiting", limit)
		if err != nil {
			return nil, fmt.Errorf("failed to get waiting records: %w", err)
		}
		return records, nil
	}

	var waitingRecords []models.WaitingRecord
	if err := s.db.Where("status = ?", "waiting").
		Order("priority DESC, queued_at ASC").
		Limit(limit).
		Find(&waitingRecords).Error; err != nil {
		return nil, fmt.Errorf("failed to get waiting records: %w", err)
	}
	return waitingRecords, nil
}

func (s *SessionTransferService) getActiveWaitingRecord(ctx context.Context, sessionID string) (*models.WaitingRecord, error) {
	if s.routing != nil {
		record, err := s.routing.GetWaitingRecord(ctx, sessionID)
		if err != nil {
			return nil, err
		}
		if record.Status != "waiting" {
			return nil, gorm.ErrRecordNotFound
		}
		return record, nil
	}

	var existing models.WaitingRecord
	if err := s.db.Where("session_id = ? AND status = ?", sessionID, "waiting").First(&existing).Error; err != nil {
		return nil, err
	}
	return &existing, nil
}

// generateSessionSummary 生成会话摘要
func (s *SessionTransferService) generateSessionSummary(session *models.Session) (string, error) {
	// 获取会话消息
	var messages []models.Message
	if err := s.db.Where("session_id = ?", session.ID).
		Order("created_at ASC").
		Find(&messages).Error; err != nil {
		return "", err
	}

	// 如果消息太少，返回简单摘要
	if len(messages) < 3 {
		userLabel := session.User.Username
		if userLabel == "" {
			userLabel = fmt.Sprintf("ID=%d", session.UserID)
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
	if s.wsHub != nil {
		s.wsHub.SendToSession(sessionID, WebSocketMessage{
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
	if s.wsHub != nil {
		s.wsHub.SendToSession(sessionID, WebSocketMessage{
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
	if s.routing != nil {
		return s.routing.GetTransferHistory(ctx, sessionID)
	}
	var records []models.TransferRecord
	err := s.db.Where("session_id = ?", sessionID).
		Order("transferred_at DESC").
		Find(&records).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get transfer history: %w", err)
	}

	return records, nil
}

// ListWaitingRecords 列出等待队列记录（默认 status=waiting）
func (s *SessionTransferService) ListWaitingRecords(ctx context.Context, status string, limit int) ([]models.WaitingRecord, error) {
	if s.routing != nil {
		return s.routing.ListWaitingRecords(ctx, status, limit)
	}
	if status == "" {
		status = "waiting"
	}
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	var records []models.WaitingRecord
	if err := s.db.Where("status = ?", status).
		Order("priority DESC, queued_at ASC").
		Limit(limit).
		Find(&records).Error; err != nil {
		return nil, fmt.Errorf("failed to list waiting records: %w", err)
	}
	return records, nil
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
		if s.routing != nil {
			if _, err := s.routing.CancelWaiting(ctx, tx, sessionID, reason); err != nil {
				if err == gorm.ErrRecordNotFound {
					return nil
				}
				return fmt.Errorf("update waiting record: %w", err)
			}
		} else {
			var wr models.WaitingRecord
			if err := tx.Where("session_id = ? AND status = ?", sessionID, "waiting").First(&wr).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					return nil
				}
				return fmt.Errorf("load waiting record: %w", err)
			}

			if err := tx.Model(&models.WaitingRecord{}).
				Where("id = ?", wr.ID).
				Updates(map[string]interface{}{
					"status": "cancelled",
				}).Error; err != nil {
				return fmt.Errorf("update waiting record: %w", err)
			}
		}

		// 记录系统消息（可选，不影响主流程）
		msg := &models.Message{
			SessionID: sessionID,
			UserID:    operatorID,
			Content:   fmt.Sprintf("已取消人工客服等待队列（原因：%s）", reason),
			Type:      "system",
			Sender:    "system",
			CreatedAt: now,
		}
		_ = tx.Create(msg).Error
		return nil
	})
}

func (s *SessionTransferService) SetRoutingAdapter(adapter routingdelivery.RuntimeService) {
	s.routing = adapter
}

func (s *SessionTransferService) SetTicketRuntime(adapter ticketdelivery.RuntimeService) {
	s.tickets = adapter
}

func (s *SessionTransferService) SetConversationRuntime(adapter conversationdelivery.RuntimeService) {
	s.conversation = adapter
}

func (s *SessionTransferService) SetAgentRuntime(adapter agentdelivery.RuntimeService) {
	s.agents = adapter
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
