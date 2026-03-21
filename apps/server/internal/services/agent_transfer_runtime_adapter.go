package services

import (
	"context"

	"gorm.io/gorm"
)

type AgentTransferRuntimeAdapter struct{}

func NewAgentTransferRuntimeAdapter() *AgentTransferRuntimeAdapter {
	return &AgentTransferRuntimeAdapter{}
}

func (a *AgentTransferRuntimeAdapter) SyncTransferLoad(ctx context.Context, tx *gorm.DB, fromAgentID *uint, toAgentID uint) error {
	if fromAgentID != nil && *fromAgentID != toAgentID {
		if err := tx.WithContext(ctx).Exec(
			`UPDATE agents SET current_load = CASE WHEN current_load > 0 THEN current_load - 1 ELSE 0 END WHERE user_id = ?`,
			*fromAgentID,
		).Error; err != nil {
			return err
		}
	}
	return tx.WithContext(ctx).
		Table("agents").
		Where("user_id = ?", toAgentID).
		UpdateColumn("current_load", gorm.Expr("current_load + 1")).Error
}

var _ AgentTransferRuntime = (*AgentTransferRuntimeAdapter)(nil)
