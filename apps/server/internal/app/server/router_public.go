package server

import (
	"servify/apps/server/internal/handlers"
	"servify/apps/server/internal/platform/configscope"

	"github.com/gin-gonic/gin"
)

func registerPublicRoutes(r *gin.Engine, deps Dependencies) {
	public := r.Group("/public")
	handlers.RegisterCSATSurveyRoutes(public, handlers.NewCSATSurveyHandler(deps.SatisfactionService))
	handlers.RegisterPublicKnowledgeBaseRoutes(public, handlers.NewKnowledgeDocHandler(deps.KnowledgeDocHandler))
	portalResolver := configscope.NewResolver(
		deps.Config,
		configscope.WithTenantPortalProvider(configscope.NewGormTenantConfigProvider(deps.DB)),
		configscope.WithWorkspacePortalProvider(configscope.NewGormWorkspaceConfigProvider(deps.DB)),
	)
	public.GET("/portal/config", handlers.NewPortalConfigHandlerWithResolver(deps.Config, portalResolver).Get)
}
