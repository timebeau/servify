package delivery

import (
	"context"

	gamificationcontract "servify/apps/server/internal/modules/gamification/contract"
	"time"
)

type LeaderboardRequest struct {
	StartDate  time.Time
	EndDate    time.Time
	Limit      int
	Department string
}

type HandlerService interface {
	GetLeaderboard(ctx context.Context, req *LeaderboardRequest) (*gamificationcontract.LeaderboardResponse, error)
}
