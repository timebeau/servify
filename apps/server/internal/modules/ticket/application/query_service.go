package application

import (
	"context"
	"fmt"
)

type QueryService struct {
	repo QueryRepository
}

func NewQueryService(repo QueryRepository) *QueryService {
	return &QueryService{repo: repo}
}

func (s *QueryService) GetTicketByID(ctx context.Context, ticketID uint) (*TicketDetailsDTO, error) {
	if ticketID == 0 {
		return nil, fmt.Errorf("ticket id required")
	}
	details, err := s.repo.GetTicketByID(ctx, ticketID)
	if err != nil {
		return nil, err
	}
	return MapTicketDetails(details), nil
}

func (s *QueryService) ListTickets(ctx context.Context, query ListTicketsQuery) (*ListTicketsResultDTO, error) {
	if query.Page <= 0 {
		query.Page = 1
	}
	if query.PageSize <= 0 {
		query.PageSize = 20
	}
	items, total, err := s.repo.ListTickets(ctx, query)
	if err != nil {
		return nil, err
	}
	return &ListTicketsResultDTO{
		Items: MapTickets(items),
		Total: total,
	}, nil
}

func (s *QueryService) GetTicketStats(ctx context.Context, agentID *uint) (*TicketStatsDTO, error) {
	return s.repo.GetTicketStats(ctx, agentID)
}
