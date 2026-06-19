package application_test

import (
	"testing"

	gamificationapp "servify/apps/server/internal/modules/gamification/application"
)

func TestComputeGamificationScore(t *testing.T) {
	tests := []struct {
		name               string
		resolved           int64
		csatAvg            float64
		csatCount          int64
		avgResponseTimeSec int64
		wantPositive       bool
	}{
		{
			name:               "high performer",
			resolved:           50,
			csatAvg:            4.8,
			csatCount:          20,
			avgResponseTimeSec: 300,
			wantPositive:       true,
		},
		{
			name:               "low csat samples",
			resolved:           10,
			csatAvg:            5.0,
			csatCount:          2,
			avgResponseTimeSec: 300,
			wantPositive:       true,
		},
		{
			name:               "slow response time penalty",
			resolved:           10,
			csatAvg:            4.0,
			csatCount:          5,
			avgResponseTimeSec: 5000,
			wantPositive:       true,
		},
		{
			name:               "minimal activity",
			resolved:           1,
			csatAvg:            3.0,
			csatCount:          1,
			avgResponseTimeSec: 100,
			wantPositive:       true,
		},
		{
			name:               "no response time",
			resolved:           10,
			csatAvg:            4.0,
			csatCount:          5,
			avgResponseTimeSec: 0,
			wantPositive:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score := gamificationapp.ComputeGamificationScore(tt.resolved, tt.csatAvg, tt.csatCount, tt.avgResponseTimeSec)
			if (score > 0) != tt.wantPositive {
				t.Errorf("ComputeGamificationScore() = %v, wantPositive %v", score, tt.wantPositive)
			}
		})
	}
}

func TestApplyBadges(t *testing.T) {
	tests := []struct {
		name              string
		entries           []gamificationapp.LeaderboardEntry
		wantBadgesInFirst bool
	}{
		{
			name: "single top resolver gets badge",
			entries: []gamificationapp.LeaderboardEntry{
				{
					AgentID:         1,
					AgentName:       "张三",
					ResolvedTickets: 50,
					CSATAvg:         4.5,
					CSATCount:       10,
					AvgResponseTime: 300,
				},
				{
					AgentID:         2,
					AgentName:       "李四",
					ResolvedTickets: 30,
					CSATAvg:         4.0,
					CSATCount:       8,
					AvgResponseTime: 400,
				},
			},
			wantBadgesInFirst: true,
		},
		{
			name: "top csat gets badge",
			entries: []gamificationapp.LeaderboardEntry{
				{
					AgentID:         1,
					AgentName:       "张三",
					ResolvedTickets: 30,
					CSATAvg:         5.0,
					CSATCount:       10,
					AvgResponseTime: 300,
				},
				{
					AgentID:         2,
					AgentName:       "李四",
					ResolvedTickets: 50,
					CSATAvg:         4.0,
					CSATCount:       8,
					AvgResponseTime: 400,
				},
			},
			wantBadgesInFirst: true,
		},
		{
			name: "fastest response gets badge",
			entries: []gamificationapp.LeaderboardEntry{
				{
					AgentID:         1,
					AgentName:       "张三",
					ResolvedTickets: 30,
					CSATAvg:         4.0,
					CSATCount:       8,
					AvgResponseTime: 100,
				},
				{
					AgentID:         2,
					AgentName:       "李四",
					ResolvedTickets: 30,
					CSATAvg:         4.0,
					CSATCount:       8,
					AvgResponseTime: 500,
				},
			},
			wantBadgesInFirst: true,
		},
		{
			name:              "empty entries",
			entries:           []gamificationapp.LeaderboardEntry{},
			wantBadgesInFirst: false,
		},
		{
			name: "all-rounder gets multiple badges",
			entries: []gamificationapp.LeaderboardEntry{
				{
					AgentID:         1,
					AgentName:       "全能张三",
					ResolvedTickets: 100,
					CSATAvg:         5.0,
					CSATCount:       20,
					AvgResponseTime: 50,
				},
				{
					AgentID:         2,
					AgentName:       "李四",
					ResolvedTickets: 10,
					CSATAvg:         3.0,
					CSATCount:       2,
					AvgResponseTime: 1000,
				},
			},
			wantBadgesInFirst: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gamificationapp.ApplyBadges(tt.entries)
			if len(tt.entries) > 0 {
				hasBadges := len(tt.entries[0].Badges) > 0
				if hasBadges != tt.wantBadgesInFirst {
					t.Errorf("ApplyBadges() first entry has badges = %v, want %v", hasBadges, tt.wantBadgesInFirst)
				}
			}
		})
	}
}
