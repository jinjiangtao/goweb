package controllers

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"xiangce/models"
	"xiangce/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type MediaListRequest struct {
	AlbumID   uint   `form:"album_id" binding:"required"`
	SortBy    string `form:"sort_by"`
	SortOrder string `form:"sort_order"`
	Page      int    `form:"page,default=1"`
	PageSize  int    `form:"page_size,default=50"`
}

type MoveMediaRequest struct {
	MediaIDs      []uint `json:"media_ids" binding:"required"`
	TargetAlbumID uint   `json:"target_album_id" binding:"required"`
}

type SortMediaRequest struct {
	MediaIDs []uint `json:"media_ids" binding:"required"`
}

func UploadMedia(c *gin.Context) {
	userID := getUserIDFromContext(c)
	albumIDStr := c.PostForm("album_id")
	albumID, err := strconv.ParseUint(albumIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的相册ID"})
		return
	}

	var album models.Album
	if result := models.DB.Where("id = ? AND user_id = ?", albumID, userID).First(&album); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "相册不存在"})
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "获取文件失败"})
		return
	}

	files := form.File["files"]
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请选择要上传的文件"})
		return
	}

	uploadDir := "./uploads"
	thumbDir := "./uploads/thumbnails"

	var uploadedMedia []models.Media
	var maxSort int
	models.DB.Model(&models.Media{}).Where("album_id = ?", albumID).Select("COALESCE(MAX(sort_order), -1)").Scan(&maxSort)

	for i, file := range files {
		if file.Size > 100*1024*1024 {
			continue
		}

		imgInfo, err := utils.SaveUploadedFile(file, uploadDir, thumbDir)
		if err != nil {
			continue
		}

		if imgInfo.FileType == "image" {
			utils.CompressImage(imgInfo.FilePath, 85)
		}

		media := models.Media{
			UserID:       userID,
			AlbumID:      uint(albumID),
			FileName:     imgInfo.FileName,
			OriginalName: file.Filename,
			FileSize:     imgInfo.FileSize,
			FileType:     imgInfo.FileType,
			MimeType:     imgInfo.MimeType,
			FilePath:     imgInfo.FilePath,
			ThumbPath:    imgInfo.ThumbPath,
			Width:        imgInfo.Width,
			Height:       imgInfo.Height,
			SortOrder:    maxSort + i + 1,
		}

		if result := models.DB.Create(&media); result.Error == nil {
			uploadedMedia = append(uploadedMedia, media)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("成功上传 %d 个文件", len(uploadedMedia)),
		"media":   uploadedMedia,
	})
}

func GetMediaList(c *gin.Context) {
	userID := getUserIDFromContext(c)
	albumID, err := parseUintParam(c, "albumId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的相册ID"})
		return
	}

	sortBy := c.DefaultQuery("sort_by", "sort_order")
	sortOrder := c.DefaultQuery("sort_order", "asc")

	validSortFields := map[string]bool{
		"sort_order": true,
		"created_at": true,
		"file_name":  true,
		"file_size":  true,
	}
	if !validSortFields[sortBy] {
		sortBy = "sort_order"
	}
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "asc"
	}

	orderClause := sortBy + " " + sortOrder
	if sortBy == "sort_order" {
		orderClause = "sort_order ASC, created_at DESC"
	}

	var medias []models.Media
	models.DB.Where("album_id = ? AND user_id = ?", albumID, userID).
		Order(orderClause).
		Find(&medias)

	c.JSON(http.StatusOK, gin.H{"media": medias})
}

func GetMedia(c *gin.Context) {
	userID := getUserIDFromContext(c)
	mediaID, err := parseUintParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的媒体ID"})
		return
	}

	var media models.Media
	if result := models.DB.Where("id = ? AND user_id = ?", mediaID, userID).First(&media); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文件不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"media": media})
}

func DeleteMedia(c *gin.Context) {
	userID := getUserIDFromContext(c)
	mediaID, err := parseUintParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的媒体ID"})
		return
	}

	var media models.Media
	if result := models.DB.Where("id = ? AND user_id = ?", mediaID, userID).First(&media); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文件不存在"})
		return
	}

	os.Remove(media.FilePath)
	os.Remove(media.ThumbPath)

	models.DB.Delete(&media)
	c.JSON(http.StatusOK, gin.H{"message": "文件已删除"})
}

func DownloadMedia(c *gin.Context) {
	userID := getUserIDFromContext(c)
	mediaID, err := parseUintParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的媒体ID"})
		return
	}

	var media models.Media
	if result := models.DB.Where("id = ? AND user_id = ?", mediaID, userID).First(&media); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文件不存在"})
		return
	}

	c.FileAttachment(media.FilePath, media.OriginalName)
}

func BatchDeleteMedia(c *gin.Context) {
	userID := getUserIDFromContext(c)

	var req struct {
		MediaIDs []uint `json:"media_ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	for _, id := range req.MediaIDs {
		var media models.Media
		if result := models.DB.Where("id = ? AND user_id = ?", id, userID).First(&media); result.Error == nil {
			os.Remove(media.FilePath)
			os.Remove(media.ThumbPath)
			models.DB.Delete(&media)
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "批量删除成功"})
}

func BatchDownloadMedia(c *gin.Context) {
	userID := getUserIDFromContext(c)

	var req struct {
		MediaIDs []uint `json:"media_ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	zipFileName := fmt.Sprintf("photos_%s.zip", uuid.New().String()[:8])
	zipFilePath := filepath.Join("./uploads", zipFileName)

	zipFile, err := os.Create(zipFilePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建压缩包失败"})
		return
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	for _, id := range req.MediaIDs {
		var media models.Media
		if result := models.DB.Where("id = ? AND user_id = ?", id, userID).First(&media); result.Error != nil {
			continue
		}

		file, err := os.Open(media.FilePath)
		if err != nil {
			continue
		}

		w, err := zipWriter.Create(media.OriginalName)
		if err != nil {
			file.Close()
			continue
		}

		io.Copy(w, file)
		file.Close()
	}

	c.FileAttachment(zipFilePath, fmt.Sprintf("相册下载_%s.zip", time.Now().Format("20060102")))
	go func() {
		time.Sleep(10 * time.Minute)
		os.Remove(zipFilePath)
	}()
}

func MoveMedia(c *gin.Context) {
	userID := getUserIDFromContext(c)

	var req MoveMediaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	var targetAlbum models.Album
	if result := models.DB.Where("id = ? AND user_id = ?", req.TargetAlbumID, userID).First(&targetAlbum); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "目标相册不存在"})
		return
	}

	models.DB.Model(&models.Media{}).
		Where("id IN ? AND user_id = ?", req.MediaIDs, userID).
		Update("album_id", req.TargetAlbumID)

	c.JSON(http.StatusOK, gin.H{"message": "移动成功"})
}

func UpdateMediaSort(c *gin.Context) {
	userID := getUserIDFromContext(c)
	albumID, err := parseUintParam(c, "albumId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的相册ID"})
		return
	}

	var req SortMediaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	for i, id := range req.MediaIDs {
		models.DB.Model(&models.Media{}).
			Where("id = ? AND album_id = ? AND user_id = ?", id, albumID, userID).
			Update("sort_order", i)
	}

	c.JSON(http.StatusOK, gin.H{"message": "排序已更新"})
}

func ServeFile(c *gin.Context) {
	filename := c.Param("filename")
	filePath := filepath.Join("./uploads", filename)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "文件不存在"})
		return
	}

	c.File(filePath)
}

func ServeThumbnail(c *gin.Context) {
	filename := c.Param("filename")
	filePath := filepath.Join("./uploads/thumbnails", filename)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		originalPath := filepath.Join("./uploads", filename)
		if _, err := os.Stat(originalPath); os.IsNotExist(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "文件不存在"})
			return
		}
		c.File(originalPath)
		return
	}

	c.File(filePath)
}
