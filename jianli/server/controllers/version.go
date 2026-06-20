package controllers

import (
	"jianli-server/db"
	"jianli-server/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type VersionController struct{}

func NewVersionController() *VersionController {
	return &VersionController{}
}

func (vc *VersionController) GetVersions(c *gin.Context) {
	resumeIDStr := c.Param("id")
	resumeID, err := strconv.ParseUint(resumeIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的简历ID",
		})
		return
	}

	var resume models.Resume
	if err := db.DB.First(&resume, uint(resumeID)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "简历不存在",
		})
		return
	}

	var versions []models.Version
	if err := db.DB.Where("resume_id = ?", uint(resumeID)).Order("version_number DESC").Find(&versions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取版本列表失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    versions,
	})
}

func (vc *VersionController) GetVersion(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的版本ID",
		})
		return
	}

	var version models.Version
	if err := db.DB.First(&version, uint(id)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "版本不存在",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    version,
	})
}

func (vc *VersionController) RestoreVersion(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的版本ID",
		})
		return
	}

	var sourceVersion models.Version
	if err := db.DB.First(&sourceVersion, uint(id)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "版本不存在",
		})
		return
	}

	var resume models.Resume
	if err := db.DB.First(&resume, sourceVersion.ResumeID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "关联简历不存在",
		})
		return
	}

	tx := db.DB.Begin()

	newVersionNumber := resume.CurrentVersion + 1
	snapshotName := "恢复至版本 " + strconv.FormatUint(uint64(sourceVersion.VersionNumber), 10)

	newVersion := models.Version{
		ResumeID:      resume.ID,
		VersionNumber: newVersionNumber,
		SnapshotName:  snapshotName,
		Content:       sourceVersion.Content,
		StyleConfig:   sourceVersion.StyleConfig,
		CreatedAt:     time.Now(),
	}

	if err := tx.Create(&newVersion).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "创建新版本失败",
			"error":   err.Error(),
		})
		return
	}

	resume.Content = sourceVersion.Content
	resume.StyleConfig = sourceVersion.StyleConfig
	resume.CurrentVersion = newVersionNumber
	resume.UpdatedAt = time.Now()

	if err := tx.Save(&resume).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "更新简历失败",
			"error":   err.Error(),
		})
		return
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "恢复成功",
		"data": gin.H{
			"resume":  resume,
			"version": newVersion,
		},
	})
}
