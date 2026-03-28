package server

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

var defaultStaticRoots = []string{
	"./apps/admin",
	"../admin",
	"/app/apps/admin",
}

func registerStatic(r staticRegistrar) {
	root := detectStaticRoot(defaultStaticRoots)
	// Use NoRoute to serve static files as a fallback
	r.NoRoute(func(c *gin.Context) {
		// Don't handle API routes
		if strings.HasPrefix(c.Request.URL.Path, "/api") || strings.HasPrefix(c.Request.URL.Path, "/public") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
			return
		}

		// Try to serve the requested file
		reqPath := c.Request.URL.Path
		if reqPath == "/" {
			reqPath = "/index.html"
		}
		fullPath := filepath.Join(root, filepath.Clean(reqPath))

		// Check if file exists
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			// Fall back to index.html for SPA routing
			c.File(filepath.Join(root, "index.html"))
			return
		}

		c.File(fullPath)
	})
}

func detectStaticRoot(candidates []string) string {
	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	return "./apps/admin"
}

type staticRegistrar interface {
	NoRoute(handlers ...gin.HandlerFunc)
}
