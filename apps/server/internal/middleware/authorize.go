package middleware

import (
	platformauth "servify/apps/server/internal/platform/auth"

	"github.com/gin-gonic/gin"
)

func RequireRolesAny(required ...string) gin.HandlerFunc {
	return platformauth.RequireRolesAny(required...)
}

func RequirePrincipalKinds(required ...string) gin.HandlerFunc {
	return platformauth.RequirePrincipalKinds(required...)
}
