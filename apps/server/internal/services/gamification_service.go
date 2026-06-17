package services

import (
	"context"

	gamificationcontract "servify/apps/server/internal/modules/gamification/contract"
	gamificationdelivery "servify/apps/server/internal/modules/gamification/delivery"

	"gorm.io/gorm"
)

type GamificationService struct {
	service gamificationdelivery.HandlerService
}

func NewGamificationService(db *gorm.DB) *GamificationService {
	return &GamificationService{service: gamificationdelivery.NewHandlerService(db)}
}

type Badge = gamificationcontract.Badge
type LeaderboardEntry = gamificationcontract.LeaderboardEntry
type LeaderboardResponse = gamificationcontract.LeaderboardResponse
type LeaderboardRequest = gamificationdelivery.LeaderboardRequest

func (s *GamificationService) GetLeaderboard(ctx context.Context, req *LeaderboardRequest) (*LeaderboardResponse, error) {
	return s.service.GetLeaderboard(ctx, req)
}

var (
	computeGamificationScore = func(resolved int64, csatAvg float64, csatCount int64, avgResponseTimeSec int64) float64 {
		csatWeight := 20.0
		if csatCount < 3 {
			csatWeight = 10.0
		}
		respPenalty := 0.0
		if avgResponseTimeSec > 0 {
			respPenalty = float64(avgResponseTimeSec) * 0.03
		}
		return float64(resolved)*10.0 + csatAvg*csatWeight - respPenalty
	}
	applyBadges = func(entries []LeaderboardEntry) {
		if len(entries) == 0 {
			return
		}
		var (
			bestResolvedID uint
			bestResolved   int64
			bestCSATID     uint
			bestCSAT       float64
			bestCSATN      int64
			bestSpeedID    uint
			bestSpeed      int64
		)
		for _, e := range entries {
			if e.ResolvedTickets > bestResolved {
				bestResolved = e.ResolvedTickets
				bestResolvedID = e.AgentID
			}
			if e.CSATCount >= 3 {
				if e.CSATAvg > bestCSAT || (e.CSATAvg == bestCSAT && e.CSATCount > bestCSATN) {
					bestCSAT = e.CSATAvg
					bestCSATN = e.CSATCount
					bestCSATID = e.AgentID
				}
			}
			if e.AvgResponseTime > 0 {
				if bestSpeed == 0 || e.AvgResponseTime < bestSpeed {
					bestSpeed = e.AvgResponseTime
					bestSpeedID = e.AgentID
				}
			}
		}
		for i := range entries {
			var badges []Badge
			if entries[i].AgentID == bestResolvedID && bestResolved > 0 {
				badges = append(badges, Badge{ID: "top_resolver", Name: "解决王", Description: "本周期解决工单数最高"})
			}
			if entries[i].AgentID == bestCSATID && bestCSATID != 0 {
				badges = append(badges, Badge{ID: "customer_hero", Name: "满意之星", Description: "本周期 CSAT 均分最高（至少 3 条评价）"})
			}
			if entries[i].AgentID == bestSpeedID && bestSpeedID != 0 {
				badges = append(badges, Badge{ID: "speedster", Name: "极速响应", Description: "本周期平均响应时长最低（基于 Agent.avg_response_time）"})
			}
			entries[i].Badges = badges
		}
	}
)
