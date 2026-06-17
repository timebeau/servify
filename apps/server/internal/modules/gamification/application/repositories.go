package application

import "context"

type Repository interface {
	ListAgentProfiles(ctx context.Context, department string) ([]AgentProfile, error)
	ListResolvedCounts(ctx context.Context, startDate, endDate string) ([]AgentResolvedCount, error)
	ListCSATStats(ctx context.Context, startDate, endDate string) ([]AgentCSAT, error)
}
