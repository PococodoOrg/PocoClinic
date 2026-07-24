package staticserve

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

// ResolveDir returns the first directory containing index.html.
func ResolveDir(candidates ...string) string {
	for _, candidate := range candidates {
		if candidate == "" {
			continue
		}
		indexPath := filepath.Join(candidate, "index.html")
		if info, err := os.Stat(indexPath); err == nil && !info.IsDir() {
			return candidate
		}
	}
	return ""
}

// RegisterSPA serves a Vite/React build: /assets/* as files, other GET routes → index.html.
func RegisterSPA(router *gin.Engine, dir string) {
	assetsDir := filepath.Join(dir, "assets")
	router.Static("/assets", assetsDir)

	indexPath := filepath.Join(dir, "index.html")
	router.GET("/", func(c *gin.Context) {
		c.File(indexPath)
	})

	router.NoRoute(func(c *gin.Context) {
		if c.Request.Method != http.MethodGet {
			c.Status(http.StatusNotFound)
			return
		}
		path := c.Request.URL.Path
		if strings.HasPrefix(path, "/api/") || path == "/health" {
			c.Status(http.StatusNotFound)
			return
		}
		c.File(indexPath)
	})
}
