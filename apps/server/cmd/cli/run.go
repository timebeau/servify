//go:build !weknora
// +build !weknora

package cli

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	appbootstrap "servify/apps/server/internal/app/bootstrap"
	appserver "servify/apps/server/internal/app/server"
	"servify/apps/server/internal/config"
	"servify/apps/server/internal/handlers"
	"servify/apps/server/internal/middleware"
	"servify/apps/server/internal/platform/llm/openai"
	"servify/apps/server/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the servify application",
	Long:  `Run the servify application`,
	Run:   run,
}

func init() {
	rootCmd.AddCommand(runCmd)
}

func run(cmd *cobra.Command, args []string) {
	cfg, err := appbootstrap.LoadConfig("")
	if err != nil {
		logrus.Fatalf("Failed to load config: %v", err)
	}
	app := appbootstrap.BuildApp(cfg)
	if app.Logger, err = appbootstrap.InitLogging(cfg); err != nil {
		logrus.Fatalf("Failed to initialize logger: %v", err)
	}
	if err := appbootstrap.SetupObservability(context.Background(), cfg, app); err != nil {
		logrus.Warnf("init tracing: %v", err)
	}

	db, err := appbootstrap.OpenDatabase(cfg, appbootstrap.DatabaseOptions{})
	if err != nil {
		logrus.Warnf("DB connect failed, message persistence disabled: %v", err)
	}
	app.DB = db

	openAIProvider := openai.NewProvider(cfg.AI.OpenAI.APIKey, cfg.AI.OpenAI.BaseURL)
	aiService := services.NewOrchestratedAIService(openAIProvider, nil)
	runtime := appserver.BuildRealtimeRuntime(cfg, app.Logger, db, aiService)
	if err := runtime.Start(); err != nil {
		logrus.Fatalf("Failed to start message router: %v", err)
	}

	// 设置 Gin 模式
	if cfg.Server.Host != "localhost" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 创建路由
	router := setupRouter(cfg, runtime)

	server := appbootstrap.NewHTTPServer(cfg, router, appbootstrap.HTTPServerOptions{})
	appbootstrap.StartHTTPServer(server, logrus.StandardLogger(), fmt.Sprintf("Starting server on %s", server.Addr))
	appbootstrap.WaitForShutdownSignal()

	logrus.Info("Shutting down server...")

	ctx, cancel := appbootstrap.ShutdownContext(30 * time.Second)
	defer cancel()

	// 停止消息路由
	if err := runtime.Stop(ctx); err != nil {
		logrus.Errorf("Failed to stop message router: %v", err)
	}

	// 关闭服务器
	if err := server.Shutdown(ctx); err != nil {
		logrus.Errorf("Server forced to shutdown: %v", err)
	}
	if err := app.RunShutdownHooks(); err != nil {
		logrus.Errorf("Failed to run shutdown hooks: %v", err)
	}

	logrus.Info("Server exited")
}

func setupRouter(cfg *config.Config, runtime *appserver.RealtimeRuntime) *gin.Engine {
	router := gin.New()

	// 中间件
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(corsMiddlewareWithConfig(cfg))
	router.Use(middleware.RateLimitMiddlewareFromConfig(cfg))
	if cfg.Monitoring.Tracing.Enabled {
		router.Use(otelgin.Middleware(cfg.Monitoring.Tracing.ServiceName))
	}

	// 健康检查
	healthHandler := handlers.NewHealthHandler()
	router.GET("/health", healthHandler.Health)
	router.GET("/ready", healthHandler.Ready)

	// API 路由组
	api := router.Group("/api/v1")
	{
		// WebSocket 连接
		wsHandler := handlers.NewWebSocketHandler(runtime.RealtimeGateway)
		api.GET("/ws", wsHandler.HandleWebSocket)
		api.GET("/ws/stats", wsHandler.GetStats)

		// WebRTC 相关
		webrtcHandler := handlers.NewWebRTCHandler(runtime.RTCGateway)
		api.GET("/webrtc/stats", webrtcHandler.GetStats)
		api.GET("/webrtc/connections", webrtcHandler.GetConnections)

		// 消息路由
		messageHandler := handlers.NewMessageHandler(runtime.MessageRouter)
		api.GET("/messages/platforms", messageHandler.GetPlatformStats)

		// 轻量指标上报（可选）
		ingest := handlers.NewMetricsIngestHandler(handlers.NewMetricsAggregator())
		api.POST("/metrics/ingest", ingest.Ingest)
	}

	// 静态文件服务（尝试多路径）
	staticRoots := []string{
		"./apps/demo-web",
		"../demo-web",
		"/app/apps/demo-web",
	}
	sr := "./apps/demo-web"
	for _, p := range staticRoots {
		if fi, err := os.Stat(p); err == nil && fi.IsDir() {
			sr = p
			break
		}
	}
	router.Static("/", sr)

	return router
}

func corsMiddlewareWithConfig(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		origins := "*"
		methods := "GET, POST, PUT, DELETE, OPTIONS"
		headers := "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With"
		if cfg != nil && cfg.Security.CORS.Enabled {
			if len(cfg.Security.CORS.AllowedOrigins) > 0 {
				origins = strings.Join(cfg.Security.CORS.AllowedOrigins, ", ")
			}
			if len(cfg.Security.CORS.AllowedMethods) > 0 {
				methods = strings.Join(cfg.Security.CORS.AllowedMethods, ", ")
			}
			if len(cfg.Security.CORS.AllowedHeaders) > 0 {
				headers = strings.Join(cfg.Security.CORS.AllowedHeaders, ", ")
			}
		}
		c.Header("Access-Control-Allow-Origin", origins)
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", headers)
		c.Header("Access-Control-Allow-Methods", methods)
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
