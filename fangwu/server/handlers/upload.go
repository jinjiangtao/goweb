package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"fangwu-server/utils"

	"github.com/gin-gonic/gin"
)

func UploadImage(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		utils.Fail(c, http.StatusBadRequest, "请选择要上传的图片")
		return
	}

	ext := filepath.Ext(file.Filename)
	if ext == "" {
		ext = ".jpg"
	}

	allowedExts := map[string]bool{
		".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true,
	}
	if !allowedExts[ext] {
		utils.Fail(c, http.StatusBadRequest, "仅支持 jpg/png/gif/webp 格式的图片")
		return
	}

	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	uploadDir := "uploads"
	os.MkdirAll(uploadDir, 0755)
	savePath := filepath.Join(uploadDir, filename)

	if err := c.SaveUploadedFile(file, savePath); err != nil {
		utils.Fail(c, http.StatusInternalServerError, "图片保存失败")
		return
	}

	url := "/uploads/" + filename
	utils.Success(c, gin.H{"url": url})
}
