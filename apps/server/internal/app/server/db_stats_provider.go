package server

import (
	"database/sql"

	"gorm.io/gorm"
)

type gormDBStatsProvider struct {
	db *gorm.DB
}

func newGormDBStatsProvider(db *gorm.DB) *gormDBStatsProvider {
	if db == nil {
		return nil
	}
	return &gormDBStatsProvider{db: db}
}

func (p *gormDBStatsProvider) Stats() (sql.DBStats, bool) {
	if p == nil || p.db == nil {
		return sql.DBStats{}, false
	}
	sqlDB, err := p.db.DB()
	if err != nil || sqlDB == nil {
		return sql.DBStats{}, false
	}
	return sqlDB.Stats(), true
}
