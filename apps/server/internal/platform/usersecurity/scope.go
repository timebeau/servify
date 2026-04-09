package usersecurity

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"servify/apps/server/internal/models"
	platformauth "servify/apps/server/internal/platform/auth"

	"gorm.io/gorm"
)

func scopedUserIDs(ctx context.Context, db *gorm.DB, userIDs []uint) (map[uint]struct{}, error) {
	if db == nil {
		return nil, gorm.ErrInvalidDB
	}
	if len(userIDs) == 0 {
		return map[uint]struct{}{}, nil
	}

	tenantID := strings.TrimSpace(platformauth.TenantIDFromContext(ctx))
	workspaceID := strings.TrimSpace(platformauth.WorkspaceIDFromContext(ctx))

	allowed := make(map[uint]struct{}, len(userIDs))
	if tenantID == "" && workspaceID == "" {
		for _, userID := range userIDs {
			if userID != 0 {
				allowed[userID] = struct{}{}
			}
		}
		return allowed, nil
	}

	collect := func(model any) error {
		var scoped []uint
		tx := db.WithContext(ctx).Model(model).Distinct().Where("user_id IN ?", userIDs)
		if tenantID != "" {
			tx = tx.Where("tenant_id = ?", tenantID)
		}
		if workspaceID != "" {
			tx = tx.Where("workspace_id = ?", workspaceID)
		}
		if err := tx.Pluck("user_id", &scoped).Error; err != nil {
			return err
		}
		for _, userID := range scoped {
			allowed[userID] = struct{}{}
		}
		return nil
	}

	if err := collect(&models.Agent{}); err != nil {
		return nil, err
	}
	if err := collect(&models.Customer{}); err != nil {
		return nil, err
	}
	return allowed, nil
}

func ensureScopedUserAccess(ctx context.Context, db *gorm.DB, userID uint) error {
	if userID == 0 {
		return fmt.Errorf("user_id required")
	}
	allowed, err := scopedUserIDs(ctx, db, []uint{userID})
	if err != nil {
		return err
	}
	if _, ok := allowed[userID]; !ok {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func hasRequestScope(ctx context.Context) bool {
	tenantID := strings.TrimSpace(platformauth.TenantIDFromContext(ctx))
	workspaceID := strings.TrimSpace(platformauth.WorkspaceIDFromContext(ctx))
	return tenantID != "" || workspaceID != ""
}

func scopedUserIDQuery(ctx context.Context, db *gorm.DB, model any) *gorm.DB {
	tenantID := strings.TrimSpace(platformauth.TenantIDFromContext(ctx))
	workspaceID := strings.TrimSpace(platformauth.WorkspaceIDFromContext(ctx))

	tx := db.WithContext(ctx).Model(model).Distinct().Select("user_id")
	if tenantID != "" {
		tx = tx.Where("tenant_id = ?", tenantID)
	}
	if workspaceID != "" {
		tx = tx.Where("workspace_id = ?", workspaceID)
	}
	return tx
}

func orderedUniqueUserIDs(userIDs []uint) ([]uint, error) {
	if len(userIDs) == 0 {
		return []uint{}, nil
	}

	ordered := make([]uint, 0, len(userIDs))
	seen := make(map[uint]struct{}, len(userIDs))
	for _, userID := range userIDs {
		if userID == 0 {
			return nil, fmt.Errorf("user_id required")
		}
		if _, ok := seen[userID]; ok {
			continue
		}
		seen[userID] = struct{}{}
		ordered = append(ordered, userID)
	}
	return ordered, nil
}

func missingScopedUserIDs(expected []uint, allowed map[uint]struct{}, found map[uint]struct{}) []uint {
	missing := make([]uint, 0)
	for _, userID := range expected {
		if _, ok := allowed[userID]; !ok {
			missing = append(missing, userID)
			continue
		}
		if _, ok := found[userID]; !ok {
			missing = append(missing, userID)
		}
	}
	sort.Slice(missing, func(i, j int) bool { return missing[i] < missing[j] })
	return missing
}
