package delivery

import (
	"context"

	gamificationapp "servify/apps/server/internal/modules/gamification/application"
	gamificationcontract "servify/apps/server/internal/modules/gamification/contract"
	gamificationinfra "servify/apps/server/internal/modules/gamification/infra"

	"gorm.io/gorm"
)

type HandlerServiceAdapter struct {
	service *gamificationapp.Service
}

func NewHandlerService(db *gorm.DB) *HandlerServiceAdapter {
	return NewHandlerServiceAdapter(gamificationapp.NewService(gamificationinfra.NewGormRepository(db)))
}

func NewHandlerServiceAdapter(service *gamificationapp.Service) *HandlerServiceAdapter {
	return &HandlerServiceAdapter{service: service}
}

func (a *HandlerServiceAdapter) GetLeaderboard(ctx context.Context, req *LeaderboardRequest) (*gamificationcontract.LeaderboardResponse, error) {
	appReq := &gamificationapp.LeaderboardRequest{}
	if req != nil {
		appReq.StartDate = req.StartDate
		appReq.EndDate = req.EndDate
		appReq.Limit = req.Limit
		appReq.Department = req.Department
	}
	resp, err := a.service.GetLeaderboard(ctx, appReq)
	if err != nil || resp == nil {
		return nil, err
	}
	out := &gamificationcontract.LeaderboardResponse{
		StartDate: resp.StartDate,
		EndDate:   resp.EndDate,
		Limit:     resp.Limit,
		Entries:   make([]gamificationcontract.LeaderboardEntry, 0, len(resp.Entries)),
	}
	for _, entry := range resp.Entries {
		item := gamificationcontract.LeaderboardEntry{
			Rank:            entry.Rank,
			AgentID:         entry.AgentID,
			AgentName:       entry.AgentName,
			Username:        entry.Username,
			Department:      entry.Department,
			ResolvedTickets: entry.ResolvedTickets,
			CSATAvg:         entry.CSATAvg,
			CSATCount:       entry.CSATCount,
			AvgResponseTime: entry.AvgResponseTime,
			Score:           entry.Score,
			Badges:          make([]gamificationcontract.Badge, 0, len(entry.Badges)),
		}
		for _, badge := range entry.Badges {
			item.Badges = append(item.Badges, gamificationcontract.Badge{
				ID:          badge.ID,
				Name:        badge.Name,
				Description: badge.Description,
			})
		}
		out.Entries = append(out.Entries, item)
	}
	return out, nil
}

var _ HandlerService = (*HandlerServiceAdapter)(nil)
