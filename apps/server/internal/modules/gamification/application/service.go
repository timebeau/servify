package application

import (
	"context"
	"errors"
	"sort"
	"strings"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetLeaderboard(ctx context.Context, req *LeaderboardRequest) (*LeaderboardResponse, error) {
	if req == nil {
		return nil, errors.New("request required")
	}
	start := req.StartDate
	end := req.EndDate
	if start.IsZero() || end.IsZero() {
		return nil, errors.New("start_date and end_date required")
	}
	if end.Before(start) {
		return nil, errors.New("end_date must be after start_date")
	}

	limit := req.Limit
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	profiles, err := s.repo.ListAgentProfiles(ctx, strings.TrimSpace(req.Department))
	if err != nil {
		return nil, err
	}
	resolved, err := s.repo.ListResolvedCounts(ctx, start.Format("2006-01-02 15:04:05"), end.Format("2006-01-02 15:04:05"))
	if err != nil {
		return nil, err
	}
	csats, err := s.repo.ListCSATStats(ctx, start.Format("2006-01-02 15:04:05"), end.Format("2006-01-02 15:04:05"))
	if err != nil {
		return nil, err
	}

	resolvedByAgent := make(map[uint]int64, len(resolved))
	for _, r := range resolved {
		resolvedByAgent[r.AgentID] = r.Count
	}
	csatByAgent := make(map[uint]AgentCSAT, len(csats))
	for _, c := range csats {
		csatByAgent[c.AgentID] = c
	}

	entries := make([]LeaderboardEntry, 0, len(profiles))
	for _, p := range profiles {
		rc := resolvedByAgent[p.UserID]
		csat := csatByAgent[p.UserID]
		if rc == 0 && csat.Count == 0 {
			continue
		}
		score := ComputeGamificationScore(rc, csat.Avg, csat.Count, p.AvgResponseTime)
		entries = append(entries, LeaderboardEntry{
			AgentID:         p.UserID,
			AgentName:       p.Name,
			Username:        p.Username,
			Department:      p.Department,
			ResolvedTickets: rc,
			CSATAvg:         csat.Avg,
			CSATCount:       csat.Count,
			AvgResponseTime: p.AvgResponseTime,
			Score:           score,
		})
	}

	ApplyBadges(entries)
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Score == entries[j].Score {
			if entries[i].ResolvedTickets == entries[j].ResolvedTickets {
				return entries[i].AgentID < entries[j].AgentID
			}
			return entries[i].ResolvedTickets > entries[j].ResolvedTickets
		}
		return entries[i].Score > entries[j].Score
	})
	if len(entries) > limit {
		entries = entries[:limit]
	}
	for i := range entries {
		entries[i].Rank = i + 1
	}

	return &LeaderboardResponse{
		StartDate: start.Format("2006-01-02"),
		EndDate:   end.Format("2006-01-02"),
		Limit:     limit,
		Entries:   entries,
	}, nil
}

func ComputeGamificationScore(resolved int64, csatAvg float64, csatCount int64, avgResponseTimeSec int64) float64 {
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

func ApplyBadges(entries []LeaderboardEntry) {
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
