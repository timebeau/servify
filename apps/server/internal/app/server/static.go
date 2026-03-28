package server

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

var defaultStaticRoots = []string{
	"./apps/admin/dist",
	"./apps/admin",
	"../admin/dist",
	"../admin",
	"/app/apps/admin/dist",
	"/app/apps/admin",
}

// demoStaticDirs contains paths that serve the demo site and SDK assets.
var demoStaticDirs = []string{
	"./apps/demo",
	"./apps/demo-sdk",
}

func registerStatic(r staticRegistrar) {
	root := detectStaticRoot(defaultStaticRoots)

	// Serve demo-sdk assets directly (no auth)
	r.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path

		// Don't handle API routes
		if strings.HasPrefix(path, "/api") || strings.HasPrefix(path, "/public") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
			return
		}

		// Serve demo-sdk assets from ./apps/demo-sdk/
		if rest, ok := strings.CutPrefix(path, "/demo-sdk/"); ok {
			filePath := filepath.Join(".", "apps", "demo-sdk", filepath.Clean(rest))
			if _, err := os.Stat(filePath); err == nil {
				c.File(filePath)
				return
			}
		}

		// Try to serve the requested file from admin dist
		reqPath := path
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
	return "./apps/admin/dist"
}

type staticRegistrar interface {
	NoRoute(handlers ...gin.HandlerFunc)
}
