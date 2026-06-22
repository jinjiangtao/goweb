package api

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"

	"qrcode/internal/store"
)

type Service struct {
	Store   *store.Store
	DataDir string
}

func NewRouter(svc *Service) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(corsMiddleware())

	api := r.Group("/api")
	{
		api.POST("/qrcode/preview", svc.preview)
		api.POST("/qrcode/generate", svc.generate)
		api.POST("/qrcode/batch", svc.batch)

		api.POST("/parse/text", svc.parseText)
		api.POST("/parse/image", svc.parseImage)

		api.GET("/templates", svc.listTemplates)
		api.POST("/templates", svc.createTemplate)
		api.GET("/templates/:id/preview", svc.templatePreview)
		api.PUT("/templates/:id", svc.updateTemplate)
		api.DELETE("/templates/:id", svc.deleteTemplate)

		api.GET("/records", svc.listRecords)
		api.GET("/records/:id", svc.getRecord)
		api.GET("/records/:id/image", svc.recordImage)
		api.PATCH("/records/:id/status", svc.updateStatus)
		api.PUT("/records/:id/content", svc.updateContent)
		api.DELETE("/records/:id", svc.deleteRecord)
		api.POST("/records/export", svc.exportRecords)

		api.GET("/stats", svc.stats)
	}
	return r
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

func writeJSON(c *gin.Context, status int, data interface{}) {
	c.JSON(status, data)
}

func (svc *Service) absPath(p string) string {
	return filepath.Join(svc.DataDir, p)
}

func (svc *Service) savePNG(b []byte) (string, error) {
	dir := filepath.Join(svc.DataDir, "qrcodes")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	name := fmt.Sprintf("%s.png", randID())
	full := filepath.Join(dir, name)
	if err := os.WriteFile(full, b, 0o644); err != nil {
		return "", err
	}
	return filepath.ToSlash(filepath.Join("qrcodes", name)), nil
}

func (svc *Service) overwritePNG(path string, b []byte) error {
	full := svc.absPath(path)
	return os.WriteFile(full, b, 0o644)
}

func randID() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

func (svc *Service) imagePNG(c *gin.Context, b []byte) {
	c.Header("Content-Type", "image/png")
	c.Header("Cache-Control", "no-store")
	c.Data(http.StatusOK, "image/png", b)
}

func sanitizeFilename(s string, max int) string {
	s = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '-' || r == '_' {
			return r
		}
		return '_'
	}, s)
	if len(s) > max {
		s = s[:max]
	}
	if s == "" {
		s = "qr"
	}
	return s
}
