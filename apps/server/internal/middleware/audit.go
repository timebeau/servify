package middleware

import (
	auditplatform "servify/apps/server/internal/platform/audit"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AuditMiddleware(db *gorm.DB) gin.HandlerFunc {
	return auditplatform.Middleware(auditplatform.NewGormRecorder(db))
}
