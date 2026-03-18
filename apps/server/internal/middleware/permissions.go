package middleware

import (
	platformauth "servify/apps/server/internal/platform/auth"

	"github.com/gin-gonic/gin"
)

func HasPermission(granted []string, required string) bool {
	return platformauth.HasPermission(granted, required)
}

func RequirePermissionsAny(required ...string) gin.HandlerFunc {
	return platformauth.RequirePermissionsAny(required...)
}

func RequirePermissionsAll(required ...string) gin.HandlerFunc {
	return platformauth.RequirePermissionsAll(required...)
}

func RequireResourcePermission(resource string) gin.HandlerFunc {
	return platformauth.RequireResourcePermission(resource)
}
