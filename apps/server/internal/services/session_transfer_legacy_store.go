package services

import (
	"context"
	"strings"
	"time"

	"servify/apps/server/internal/models"
	conversationdelivery "servify/apps/server/internal/modules/conversation/delivery"
	routingcontract "servify/apps/server/internal/modules/routing/contract"
	ticketapp "servify/apps/server/internal/modules/ticket/application"

	"gorm.io/gorm"
)

// sessionTransferLegacyStore keeps the remaining DB-backed compatibility logic
// used when runtime adapters are not available yet.
type sessionTransferLegacyStore struct {
	db *gorm.DB
}

func newSessionTransferLegacyStore(db *gorm.DB) *sessionTransferLegacyStore {
	return &sessionTransferLegacyStore{db: db}
}

func (s *sessionTransferLegacyStore) loadWaitingRecord(ctx context.Context, sessionID string) (*models.WaitingRecord, error) {
	var record models.WaitingRecord
	if err := s.db.WithContext(ctx).
		Where("session_id = ?", sessionID).
		Order("queued_at DESC").
		First(&record).Error; err != nil {
		return nil, err
	}
	return &record, nil
}

func (s *sessionTransferLegacyStore) listWaitingRecords(ctx context.Context, status string, limit int) ([]models.WaitingRecord, error) {
	var waitingRecords []models.WaitingRecord
	if err := s.db.WithContext(ctx).
		Where("status = ?", status).
		Order("priority DESC, queued_at ASC").
		Limit(limit).
		Find(&waitingRecords).Error; err != nil {
		return nil, err
	}
	return waitingRecords, nil
}

func (s *sessionTransferLegacyStore) listTransferHistory(ctx context.Context, sessionID string) ([]models.TransferRecord, error) {
	var records []models.TransferRecord
	if err := s.db.WithContext(ctx).
		Where("session_id = ?", sessionID).
		Order("transferred_at DESC").
		Find(&records).Error; err != nil {
		return nil, err
	}
	return records, nil
}

func (s *sessionTransferLegacyStore) loadTransferSession(ctx context.Context, sessionID string) (*conversationdelivery.TransferSession, error) {
	var model models.Session
	if err := s.db.WithContext(ctx).
		Preload("User").
		First(&model, "id = ?", sessionID).Error; err != nil {
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

func (s *sessionTransferLegacyStore) loadSessionMessages(ctx context.Context, sessionID string) ([]models.Message, error) {
	var messages []models.Message
	if err := s.db.WithContext(ctx).
		Where("session_id = ?", sessionID).
		Order("created_at ASC").
		Find(&messages).Error; err != nil {
		return nil, err
	}
	return messages, nil
}

func (s *sessionTransferLegacyStore) syncTransferSession(ctx context.Context, tx *gorm.DB, session *conversationdelivery.TransferSession, targetAgentID uint) error {
	return tx.WithContext(ctx).
		Model(&models.Session{}).
		Where("id = ?", session.ID).
		Updates(map[string]interface{}{
			"agent_id": targetAgentID,
			"status":   "active",
			"ended_at": nil,
		}).Error
}

func (s *sessionTransferLegacyStore) loadTransferTicket(tx *gorm.DB, ticketID uint) (*models.Ticket, error) {
	var ticket models.Ticket
	if err := tx.Select("id", "agent_id", "status").First(&ticket, "id = ?", ticketID).Error; err != nil {
		return nil, err
	}
	return &ticket, nil
}

func (s *sessionTransferLegacyStore) syncTransferTicket(tx *gorm.DB, ticketID uint, targetAgentID uint, transferAt time.Time) error {
	ticket, err := s.loadTransferTicket(tx, ticketID)
	if err != nil {
		return err
	}

	updates, fromStatus, toStatus := ticketapp.BuildTransferAssignmentUpdate(targetAgentID, ticket.Status)
	if err := tx.Model(&models.Ticket{}).Where("id = ?", ticket.ID).Updates(updates).Error; err != nil {
		return err
	}

	s.appendTransferTicketStatus(tx, ticket.ID, targetAgentID, fromStatus, toStatus, transferAt)
	return nil
}

func (s *sessionTransferLegacyStore) appendTransferTicketStatus(tx *gorm.DB, ticketID uint, targetAgentID uint, fromStatus string, toStatus string, transferAt time.Time) {
	_ = tx.Create(&models.TicketStatus{
		TicketID:   ticketID,
		UserID:     targetAgentID,
		FromStatus: fromStatus,
		ToStatus:   toStatus,
		Reason:     "会话转接同步指派",
		CreatedAt:  transferAt,
	}).Error
}

func (s *sessionTransferLegacyStore) syncTransferAgentLoad(tx *gorm.DB, fromAgentID *uint, targetAgentID uint) error {
	if fromAgentID != nil && *fromAgentID != targetAgentID {
		if err := tx.Exec(`UPDATE agents SET current_load = CASE WHEN current_load > 0 THEN current_load - 1 ELSE 0 END WHERE user_id = ?`, *fromAgentID).Error; err != nil {
			return err
		}
	}
	return tx.Exec(`UPDATE agents SET current_load = current_load + 1 WHERE user_id = ?`, targetAgentID).Error
}

func (s *sessionTransferLegacyStore) appendTransferSystemMessage(tx *gorm.DB, sessionID string, targetAgentID uint, content string, createdAt time.Time) error {
	return tx.Create(&models.Message{
		SessionID: sessionID,
		UserID:    targetAgentID,
		Content:   content,
		Type:      "system",
		Sender:    "system",
		CreatedAt: createdAt,
	}).Error
}

func (s *sessionTransferLegacyStore) createTransferRecord(tx *gorm.DB, fallback *models.TransferRecord) (*models.TransferRecord, error) {
	if err := tx.Create(fallback).Error; err != nil {
		return nil, err
	}
	return fallback, nil
}

func (s *sessionTransferLegacyStore) updateTransferredWaitingRecord(tx *gorm.DB, sessionID string, targetAgentID uint, transferAt time.Time) error {
	return tx.Model(&models.WaitingRecord{}).
		Where("session_id = ? AND status = ?", sessionID, "waiting").
		Updates(map[string]interface{}{
			"status":      "transferred",
			"assigned_at": transferAt,
			"assigned_to": targetAgentID,
		}).Error
}

func (s *sessionTransferLegacyStore) cancelWaitingRecord(tx *gorm.DB, sessionID string) error {
	return tx.Model(&models.WaitingRecord{}).
		Where("session_id = ? AND status = ?", sessionID, "waiting").
		Updates(map[string]interface{}{"status": "cancelled"}).Error
}

func (s *sessionTransferLegacyStore) ensureWaitingSessionState(ctx context.Context, tx *gorm.DB, session *conversationdelivery.TransferSession) error {
	return tx.WithContext(ctx).
		Model(&models.Session{}).
		Where("id = ?", session.ID).
		Updates(map[string]interface{}{
			"status":   "active",
			"agent_id": nil,
			"ended_at": nil,
		}).Error
}

func (s *sessionTransferLegacyStore) createWaitingRecord(tx *gorm.DB, sessionID string, req *routingcontract.TransferRequest) (*models.WaitingRecord, error) {
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

func (s *sessionTransferLegacyStore) appendWaitingSystemMessage(tx *gorm.DB, message *models.Message) error {
	return tx.Create(message).Error
}

func (s *sessionTransferLegacyStore) appendCancellationSystemMessage(tx *gorm.DB, sessionID string, operatorID uint, content string, createdAt time.Time) error {
	return tx.Create(&models.Message{
		SessionID: sessionID,
		UserID:    operatorID,
		Content:   content,
		Type:      "system",
		Sender:    "system",
		CreatedAt: createdAt,
	}).Error
}
