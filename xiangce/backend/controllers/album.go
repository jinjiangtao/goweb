package controllers

import (
	"net/http"

	"xiangce/models"

	"github.com/gin-gonic/gin"
)

type AlbumRequest struct {
	Name        string `json:"name" binding:"required,max=100"`
	Description string `json:"description"`
	IsPrivate   bool   `json:"is_private"`
	Password    string `json:"password"`
}

type AlbumSortRequest struct {
	AlbumIDs []uint `json:"album_ids" binding:"required"`
}

func GetAlbums(c *gin.Context) {
	userID := getUserIDFromContext(c)

	var albums []models.Album
	models.DB.Where("user_id = ?", userID).Order("sort_order ASC, created_at DESC").Find(&albums)

	for i := range albums {
		var count int64
		models.DB.Model(&models.Media{}).Where("album_id = ?", albums[i].ID).Count(&count)
		albums[i].PhotoCount = int(count)

		var firstMedia models.Media
		if result := models.DB.Where("album_id = ?", albums[i].ID).Order("sort_order ASC, created_at DESC").First(&firstMedia); result.Error == nil {
			albums[i].Cover = firstMedia.ThumbPath
		}
	}

	c.JSON(http.StatusOK, gin.H{"albums": albums})
}

func GetAlbum(c *gin.Context) {
	userID := getUserIDFromContext(c)
	albumID, err := parseUintParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的相册ID"})
		return
	}

	var album models.Album
	if result := models.DB.Where("id = ? AND user_id = ?", albumID, userID).First(&album); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "相册不存在"})
		return
	}

	var count int64
	models.DB.Model(&models.Media{}).Where("album_id = ?", album.ID).Count(&count)
	album.PhotoCount = int(count)

	var firstMedia models.Media
	if result := models.DB.Where("album_id = ?", album.ID).Order("sort_order ASC, created_at DESC").First(&firstMedia); result.Error == nil {
		album.Cover = firstMedia.ThumbPath
	}

	c.JSON(http.StatusOK, gin.H{"album": album})
}

func CreateAlbum(c *gin.Context) {
	userID := getUserIDFromContext(c)

	var req AlbumRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误: " + err.Error()})
		return
	}

	album := models.Album{
		UserID:      userID,
		Name:        req.Name,
		Description: req.Description,
		IsPrivate:   req.IsPrivate,
	}

	if req.Password != "" {
		album.Password = req.Password
	}

	if result := models.DB.Create(&album); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建相册失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"album": album})
}

func UpdateAlbum(c *gin.Context) {
	userID := getUserIDFromContext(c)
	albumID, err := parseUintParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的相册ID"})
		return
	}

	var album models.Album
	if result := models.DB.Where("id = ? AND user_id = ?", albumID, userID).First(&album); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "相册不存在"})
		return
	}

	var req AlbumRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	album.Name = req.Name
	album.Description = req.Description
	album.IsPrivate = req.IsPrivate
	if req.Password != "" {
		album.Password = req.Password
	}

	models.DB.Save(&album)
	c.JSON(http.StatusOK, gin.H{"album": album})
}

func DeleteAlbum(c *gin.Context) {
	userID := getUserIDFromContext(c)
	albumID, err := parseUintParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的相册ID"})
		return
	}

	var album models.Album
	if result := models.DB.Where("id = ? AND user_id = ?", albumID, userID).First(&album); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "相册不存在"})
		return
	}

	models.DB.Where("album_id = ?", albumID).Delete(&models.Media{})
	models.DB.Delete(&album)

	c.JSON(http.StatusOK, gin.H{"message": "相册已删除"})
}

func UpdateAlbumSort(c *gin.Context) {
	userID := getUserIDFromContext(c)

	var req AlbumSortRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	for i, id := range req.AlbumIDs {
		models.DB.Model(&models.Album{}).
			Where("id = ? AND user_id = ?", id, userID).
			Update("sort_order", i)
	}

	c.JSON(http.StatusOK, gin.H{"message": "排序已更新"})
}
