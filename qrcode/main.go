package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"

	"qrcode/internal/api"
	"qrcode/internal/config"
	"qrcode/internal/store"
)

//go:embed all:web/dist
var distFS embed.FS

func main() {
	cfg := config.Default()

	if err := os.MkdirAll(cfg.DataDir, 0o755); err != nil {
		log.Fatalf("创建数据目录失败: %v", err)
	}
	if dir := filepath.Dir(cfg.DBPath); dir != "" {
		_ = os.MkdirAll(dir, 0o755)
	}

	st, err := store.Open(cfg.DBPath)
	if err != nil {
		log.Fatalf("打开数据库失败: %v", err)
	}

	svc := &api.Service{Store: st, DataDir: cfg.DataDir}
	r := api.NewRouter(svc)

	dist, err := fs.Sub(distFS, "web/dist")
	if err != nil {
		log.Fatalf("加载前端资源失败: %v", err)
	}
	fileServer := http.FileServer(http.FS(dist))

	r.NoRoute(func(c *gin.Context) {
		p := c.Request.URL.Path
		if strings.HasPrefix(p, "/api/") {
			c.JSON(http.StatusNotFound, gin.H{"code": 1, "message": "not found"})
			return
		}
		clean := strings.TrimPrefix(p, "/")
		if clean != "" {
			if f, err := dist.Open(clean); err == nil {
				_ = f.Close()
				fileServer.ServeHTTP(c.Writer, c.Request)
				return
			}
		}
		c.Request.URL.Path = "/"
		fileServer.ServeHTTP(c.Writer, c.Request)
	})

	addr := ":" + strings.TrimPrefix(cfg.Port, ":")
	fmt.Printf("QRForge 启动中...\n  数据库: %s\n  数据目录: %s\n  监听: http://localhost%s\n", cfg.DBPath, cfg.DataDir, addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("服务启动失败: %v", err)
	}
}
