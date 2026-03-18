package server

import (
	"os"

	"github.com/gin-gonic/gin"
)

var defaultStaticRoots = []string{
	"./apps/demo-web",
	"../demo-web",
	"/app/apps/demo-web",
}

func registerStatic(r staticRegistrar) {
	r.Static("/", detectStaticRoot(defaultStaticRoots))
}

func detectStaticRoot(candidates []string) string {
	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	return "./apps/demo-web"
}

type staticRegistrar interface {
	Static(string, string) gin.IRoutes
}
