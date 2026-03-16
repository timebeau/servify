package application

import (
	"context"
	"fmt"
	"strings"
	"time"

	"servify/apps/server/internal/modules/ticket/domain"
)

type CommandService struct {
	repo         CommandRepository
	statusPolicy StatusTransitionPolicy
}

func NewCommandService(repo CommandRepository) *CommandService {
	return &CommandService{
		repo:         repo,
		statusPolicy: NewStatusTransitionPolicy(),
	}
}

func (s *CommandService) CreateTicket(ctx context.Context, cmd CreateTicketCommand) (*TicketDTO, error) {
	if strings.TrimSpace(cmd.Title) == "" {
		return nil, fmt.Errorf("title required")
	}
	if cmd.CustomerID == 0 {
		return nil, fmt.Errorf("customer id required")
	}

	exists, err := s.repo.CustomerExists(ctx, cmd.CustomerID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, fmt.Errorf("customer not found")
	}

	now := time.Now()
	category := strings.TrimSpace(cmd.Category)
	if category == "" {
		category = "general"
	}
	priority := strings.TrimSpace(cmd.Priority)
	if priority == "" {
		priority = "normal"
	}
	source := strings.TrimSpace(cmd.Source)
	if source == "" {
		source = "web"
	}

	ticket := &domain.Ticket{
		Title:       strings.TrimSpace(cmd.Title),
		Description: strings.TrimSpace(cmd.Description),
		CustomerID:  cmd.CustomerID,
		Category:    category,
		Priority:    priority,
		Status:      "open",
		Source:      source,
		Tags:        strings.TrimSpace(cmd.Tags),
		SessionID:   cmd.SessionID,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.repo.CreateTicket(ctx, ticket); err != nil {
		return nil, err
	}
	dto := MapTicket(*ticket)
	return &dto, nil
}

func (s *CommandService) UpdateTicket(ctx context.Context, ticketID uint, cmd UpdateTicketCommand) (*TicketDTO, error) {
	if ticketID == 0 {
		return nil, fmt.Errorf("ticket id required")
	}
	ticket, err := s.repo.GetTicket(ctx, ticketID)
	if err != nil {
		return nil, err
	}

	if cmd.Title != nil {
		ticket.Title = strings.TrimSpace(*cmd.Title)
	}
	if cmd.Description != nil {
		ticket.Description = strings.TrimSpace(*cmd.Description)
	}
	if cmd.AgentID != nil {
		ticket.AgentID = cmd.AgentID
	}
	if cmd.Category != nil {
		ticket.Category = strings.TrimSpace(*cmd.Category)
	}
	if cmd.Priority != nil {
		ticket.Priority = strings.TrimSpace(*cmd.Priority)
	}
	if cmd.Status != nil {
		status := strings.TrimSpace(*cmd.Status)
		if status != ticket.Status {
			if err := s.statusPolicy.Validate(ticket.Status, status); err != nil {
				return nil, err
			}
			if err := s.recordStatusChange(ctx, ticket.ID, cmd.ActorID, ticket.Status, status, "update_ticket"); err != nil {
				return nil, err
			}
			ticket.Status = status
			now := time.Now()
			switch status {
			case "resolved":
				ticket.ResolvedAt = &now
			case "closed":
				ticket.ClosedAt = &now
			}
		}
	}
	if cmd.Tags != nil {
		ticket.Tags = strings.TrimSpace(*cmd.Tags)
	}
	if cmd.DueDate != nil {
		ticket.DueDate = cmd.DueDate
	}

	if strings.TrimSpace(ticket.Title) == "" {
		return nil, fmt.Errorf("title required")
	}
	ticket.UpdatedAt = time.Now()
	if err := s.repo.UpdateTicket(ctx, ticket); err != nil {
		return nil, err
	}
	dto := MapTicket(*ticket)
	return &dto, nil
}

func (s *CommandService) AssignTicket(ctx context.Context, ticketID uint, cmd AssignTicketCommand) (*TicketDTO, error) {
	if ticketID == 0 {
		return nil, fmt.Errorf("ticket id required")
	}
	if cmd.AgentID == 0 {
		return nil, fmt.Errorf("agent id required")
	}
	ticket, err := s.repo.GetTicket(ctx, ticketID)
	if err != nil {
		return nil, err
	}
	assignable, err := s.repo.AgentAssignable(ctx, cmd.AgentID)
	if err != nil {
		return nil, err
	}
	if !assignable {
		return nil, fmt.Errorf("agent not available")
	}

	ticket.AgentID = &cmd.AgentID
	if ticket.Status == "" || ticket.Status == "open" {
		if err := s.statusPolicy.Validate(ticket.Status, "assigned"); err != nil {
			return nil, err
		}
		if err := s.recordStatusChange(ctx, ticket.ID, cmd.UserID, ticket.Status, "assigned", "assign_ticket"); err != nil {
			return nil, err
		}
		ticket.Status = "assigned"
	}
	ticket.UpdatedAt = time.Now()
	if err := s.repo.UpdateTicket(ctx, ticket); err != nil {
		return nil, err
	}
	dto := MapTicket(*ticket)
	return &dto, nil
}

func (s *CommandService) AddComment(ctx context.Context, ticketID uint, cmd AddCommentCommand) (*CommentDTO, error) {
	if ticketID == 0 {
		return nil, fmt.Errorf("ticket id required")
	}
	if cmd.UserID == 0 {
		return nil, fmt.Errorf("user id required")
	}
	content := strings.TrimSpace(cmd.Content)
	if content == "" {
		return nil, fmt.Errorf("content required")
	}
	if _, err := s.repo.GetTicket(ctx, ticketID); err != nil {
		return nil, err
	}

	commentType := strings.TrimSpace(cmd.CommentType)
	if commentType == "" {
		commentType = "comment"
	}
	comment := &domain.Comment{
		UserID:    cmd.UserID,
		Content:   content,
		Type:      commentType,
		CreatedAt: time.Now(),
	}
	if err := s.repo.AddComment(ctx, ticketID, comment); err != nil {
		return nil, err
	}
	dto := MapComment(*comment)
	return &dto, nil
}

func (s *CommandService) CloseTicket(ctx context.Context, ticketID uint, cmd CloseTicketCommand) (*TicketDTO, error) {
	if ticketID == 0 {
		return nil, fmt.Errorf("ticket id required")
	}
	ticket, err := s.repo.GetTicket(ctx, ticketID)
	if err != nil {
		return nil, err
	}
	if err := s.statusPolicy.Validate(ticket.Status, "closed"); err != nil {
		return nil, err
	}
	now := time.Now()
	if err := s.recordStatusChange(ctx, ticket.ID, cmd.UserID, ticket.Status, "closed", cmd.Reason); err != nil {
		return nil, err
	}
	ticket.Status = "closed"
	ticket.ClosedAt = &now
	ticket.UpdatedAt = now
	if err := s.repo.UpdateTicket(ctx, ticket); err != nil {
		return nil, err
	}
	dto := MapTicket(*ticket)
	return &dto, nil
}

func (s *CommandService) recordStatusChange(ctx context.Context, ticketID uint, userID uint, fromStatus, toStatus, reason string) error {
	if fromStatus == toStatus {
		return nil
	}
	change := &domain.StatusChange{
		UserID:     userID,
		FromStatus: fromStatus,
		ToStatus:   toStatus,
		Reason:     reason,
		CreatedAt:  time.Now(),
	}
	return s.repo.RecordStatusChange(ctx, ticketID, change)
}
