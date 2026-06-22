package controllers

import (
	"net/http"
	"time"

	"xiangce/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateShareRequest struct {
	AlbumID  uint   `json:"album_id" binding:"required"`
	Password string `json:"password"`
	MaxViews int    `json:"max_views"`
	Days     int    `json:"days"`
}

type VerifyShareRequest struct {
	Password string `json:"password"`
}

func CreateShareLink(c *gin.Context) {
	userID := getUserIDFromContext(c)

	var req CreateShareRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误: " + err.Error()})
		return
	}

	var album models.Album
	if result := models.DB.Where("id = ? AND user_id = ?", req.AlbumID, userID).First(&album); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "相册不存在"})
		return
	}

	token := uuid.New().String()

	var expiresAt *time.Time
	if req.Days > 0 {
		t := time.Now().Add(time.Duration(req.Days) * 24 * time.Hour)
		expiresAt = &t
	}

	share := models.ShareLink{
		UserID:    userID,
		AlbumID:   req.AlbumID,
		Token:     token,
		Password:  req.Password,
		MaxViews:  req.MaxViews,
		ExpiresAt: expiresAt,
		IsActive:  true,
	}

	if result := models.DB.Create(&share); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建分享链接失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"share": share,
		"link":  "/share/" + token,
	})
}

func GetShareLinks(c *gin.Context) {
	userID := getUserIDFromContext(c)
	albumID, _ := parseUintParam(c, "albumId")

	var shares []models.ShareLink
	query := models.DB.Where("user_id = ?", userID)
	if albumID > 0 {
		query = query.Where("album_id = ?", albumID)
	}
	query.Order("created_at DESC").Find(&shares)

	for i := range shares {
		var album models.Album
		models.DB.Select("id, name, cover").First(&album, shares[i].AlbumID)
		shares[i].Album = album
	}

	c.JSON(http.StatusOK, gin.H{"shares": shares})
}

func DeleteShareLink(c *gin.Context) {
	userID := getUserIDFromContext(c)
	shareID, err := parseUintParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的分享ID"})
		return
	}

	var share models.ShareLink
	if result := models.DB.Where("id = ? AND user_id = ?", shareID, userID).First(&share); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "分享链接不存在"})
		return
	}

	models.DB.Delete(&share)
	c.JSON(http.StatusOK, gin.H{"message": "分享链接已删除"})
}

func GetShareByToken(c *gin.Context) {
	token := c.Param("token")

	var share models.ShareLink
	if result := models.DB.Where("token = ? AND is_active = ?", token, true).First(&share); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "分享链接无效或已失效"})
		return
	}

	if share.ExpiresAt != nil && time.Now().After(*share.ExpiresAt) {
		c.JSON(http.StatusForbidden, gin.H{"error": "分享链接已过期"})
		return
	}

	if share.MaxViews > 0 && share.ViewCount >= share.MaxViews {
		c.JSON(http.StatusForbidden, gin.H{"error": "分享链接访问次数已达上限"})
		return
	}

	needPassword := share.Password != ""

	var album models.Album
	models.DB.First(&album, share.AlbumID)

	var medias []models.Media
	models.DB.Where("album_id = ?", share.AlbumID).
		Order("sort_order ASC, created_at DESC").
		Find(&medias)

	c.JSON(http.StatusOK, gin.H{
		"share": gin.H{
			"id":            share.ID,
			"token":         share.Token,
			"need_password": needPassword,
			"expires_at":    share.ExpiresAt,
		},
		"album": album,
		"media": medias,
	})
}

func VerifySharePassword(c *gin.Context) {
	token := c.Param("token")

	var req VerifyShareRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	var share models.ShareLink
	if result := models.DB.Where("token = ? AND is_active = ?", token, true).First(&share); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "分享链接无效或已失效"})
		return
	}

	if share.Password == "" {
		c.JSON(http.StatusOK, gin.H{"valid": true})
		return
	}

	if share.Password != req.Password {
		c.JSON(http.StatusForbidden, gin.H{"error": "密码错误"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"valid": true})
}

func RecordVisit(c *gin.Context) {
	token := c.Param("token")

	var share models.ShareLink
	if result := models.DB.Where("token = ? AND is_active = ?", token, true).First(&share); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "分享链接无效"})
		return
	}

	share.ViewCount++
	models.DB.Save(&share)

	visit := models.VisitRecord{
		ShareID:   share.ID,
		IPAddress: c.ClientIP(),
		UserAgent: c.GetHeader("User-Agent"),
		VisitedAt: time.Now(),
	}
	models.DB.Create(&visit)

	c.JSON(http.StatusOK, gin.H{"message": "访问已记录"})
}

func GetVisitRecords(c *gin.Context) {
	userID := getUserIDFromContext(c)
	shareID, err := parseUintParam(c, "shareId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的分享ID"})
		return
	}

	var share models.ShareLink
	if result := models.DB.Where("id = ? AND user_id = ?", shareID, userID).First(&share); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "分享链接不存在"})
		return
	}

	var records []models.VisitRecord
	models.DB.Where("share_id = ?", shareID).
		Order("visited_at DESC").
		Limit(100).
		Find(&records)

	c.JSON(http.StatusOK, gin.H{
		"records": records,
		"total":   len(records),
	})
}
