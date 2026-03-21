package delivery

import (
	"context"
	"time"

	"servify/apps/server/internal/models"
	routingapp "servify/apps/server/internal/modules/routing/application"
)

type SessionTransferAdapter struct {
	service *routingapp.Service
}

func NewSessionTransferAdapter(service *routingapp.Service) *SessionTransferAdapter {
	return &SessionTransferAdapter{service: service}
}

func (a *SessionTransferAdapter) AddToWaitingQueue(
	ctx context.Context,
	sessionID string,
	reason string,
	targetSkills []string,
	priority string,
	notes string,
) (*models.WaitingRecord, error) {
	entry, err := a.service.AddToWaitingQueue(ctx, routingapp.AddToWaitingQueueCommand{
		SessionID:    sessionID,
		Reason:       reason,
		TargetSkills: targetSkills,
		Priority:     priority,
		Notes:        notes,
	})
	if err != nil {
		return nil, err
	}
	return mapWaitingRecord(entry), nil
}

func (a *SessionTransferAdapter) GetTransferHistory(ctx context.Context, sessionID string) ([]models.TransferRecord, error) {
	items, err := a.service.GetTransferHistory(ctx, sessionID)
	if err != nil {
		return nil, err
	}
	out := make([]models.TransferRecord, 0, len(items))
	for _, item := range items {
		out = append(out, mapTransferRecord(item))
	}
	return out, nil
}

func (a *SessionTransferAdapter) ListWaitingRecords(ctx context.Context, status string, limit int) ([]models.WaitingRecord, error) {
	items, err := a.service.ListWaitingEntries(ctx, status, limit)
	if err != nil {
		return nil, err
	}
	out := make([]models.WaitingRecord, 0, len(items))
	for _, item := range items {
		out = append(out, *mapWaitingRecord(&item))
	}
	return out, nil
}

func (a *SessionTransferAdapter) GetWaitingRecord(ctx context.Context, sessionID string) (*models.WaitingRecord, error) {
	entry, err := a.service.GetWaitingEntry(ctx, sessionID)
	if err != nil {
		return nil, err
	}
	return mapWaitingRecord(entry), nil
}

func (a *SessionTransferAdapter) CancelWaiting(ctx context.Context, sessionID string, reason string) (*models.WaitingRecord, error) {
	entry, err := a.service.CancelWaiting(ctx, routingapp.CancelWaitingCommand{
		SessionID: sessionID,
		Reason:    reason,
	})
	if err != nil {
		return nil, err
	}
	return mapWaitingRecord(entry), nil
}

func (a *SessionTransferAdapter) MarkWaitingTransferred(ctx context.Context, sessionID string, agentID uint, assignedAt time.Time) (*models.WaitingRecord, error) {
	entry, err := a.service.MarkWaitingTransferred(ctx, routingapp.MarkWaitingTransferredCommand{
		SessionID:  sessionID,
		AssignedTo: agentID,
		AssignedAt: assignedAt,
	})
	if err != nil {
		return nil, err
	}
	return mapWaitingRecord(entry), nil
}

func mapWaitingRecord(item *routingapp.QueueEntryDTO) *models.WaitingRecord {
	if item == nil {
		return nil
	}
	return &models.WaitingRecord{
		SessionID:    item.SessionID,
		Reason:       item.Reason,
		TargetSkills: marshalSkills(item.TargetSkills),
		Priority:     item.Priority,
		Notes:        item.Notes,
		Status:       item.Status,
		QueuedAt:     item.QueuedAt,
		AssignedAt:   item.AssignedAt,
		AssignedTo:   item.AssignedTo,
		CreatedAt:    item.QueuedAt,
	}
}

func marshalSkills(skills []string) string {
	if len(skills) == 0 {
		return ""
	}
	// legacy waiting records previously stored comma-separated values
	return joinSkills(skills)
}

func joinSkills(skills []string) string {
	if len(skills) == 0 {
		return ""
	}
	out := ""
	for i, skill := range skills {
		if i > 0 {
			out += ","
		}
		out += skill
	}
	return out
}

func cloneTimePtr(in *time.Time) *time.Time {
	if in == nil {
		return nil
	}
	t := *in
	return &t
}

func mapTransferRecord(item routingapp.TransferRecordDTO) models.TransferRecord {
	return models.TransferRecord{
		SessionID:      item.SessionID,
		FromAgentID:    item.FromAgentID,
		ToAgentID:      item.ToAgentID,
		Reason:         item.Reason,
		Notes:          item.Notes,
		SessionSummary: item.SessionSummary,
		TransferredAt:  item.TransferredAt,
		CreatedAt:      item.TransferredAt,
	}
}
