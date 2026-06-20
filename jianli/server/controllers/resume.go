package controllers

import (
	"jianli-server/db"
	"jianli-server/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type ResumeController struct{}

func NewResumeController() *ResumeController {
	return &ResumeController{}
}

func (rc *ResumeController) GetResumes(c *gin.Context) {
	userIdentity := c.Query("user_identity")
	if userIdentity == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "缺少 user_identity 参数",
		})
		return
	}

	var resumes []models.Resume
	if err := db.DB.Where("user_identity = ?", userIdentity).Order("updated_at DESC").Find(&resumes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取简历列表失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    resumes,
	})
}

func (rc *ResumeController) GetResume(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的简历ID",
		})
		return
	}

	var resume models.Resume
	if err := db.DB.First(&resume, uint(id)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "简历不存在",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    resume,
	})
}

func (rc *ResumeController) CreateResume(c *gin.Context) {
	var req struct {
		TemplateID   uint   `json:"template_id" binding:"required"`
		Title        string `json:"title" binding:"required"`
		UserIdentity string `json:"user_identity" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
			"error":   err.Error(),
		})
		return
	}

	var template models.Template
	if err := db.DB.First(&template, req.TemplateID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "模板不存在",
		})
		return
	}

	tx := db.DB.Begin()

	resume := models.Resume{
		UserIdentity:   req.UserIdentity,
		TemplateID:     req.TemplateID,
		Title:          req.Title,
		CurrentVersion: 1,
		Content:        template.Content,
		StyleConfig:    template.StyleConfig,
	}

	if err := tx.Create(&resume).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "创建简历失败",
			"error":   err.Error(),
		})
		return
	}

	version := models.Version{
		ResumeID:      resume.ID,
		VersionNumber: 1,
		SnapshotName:  "初始版本",
		Content:       resume.Content,
		StyleConfig:   resume.StyleConfig,
		CreatedAt:     time.Now(),
	}

	if err := tx.Create(&version).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "创建初始版本失败",
			"error":   err.Error(),
		})
		return
	}

	tx.Commit()

	c.JSON(http.StatusCreated, gin.H{
		"code":    201,
		"message": "创建成功",
		"data":    resume,
	})
}

func (rc *ResumeController) UpdateResume(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的简历ID",
		})
		return
	}

	var req struct {
		Content      models.JSON `json:"content" binding:"required"`
		StyleConfig  models.JSON `json:"style_config" binding:"required"`
		SnapshotName string      `json:"snapshot_name"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
			"error":   err.Error(),
		})
		return
	}

	var resume models.Resume
	if err := db.DB.First(&resume, uint(id)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "简历不存在",
		})
		return
	}

	tx := db.DB.Begin()

	newVersion := resume.CurrentVersion + 1
	snapshotName := req.SnapshotName
	if snapshotName == "" {
		snapshotName = "版本 " + strconv.FormatUint(uint64(newVersion), 10)
	}

	version := models.Version{
		ResumeID:      resume.ID,
		VersionNumber: newVersion,
		SnapshotName:  snapshotName,
		Content:       req.Content,
		StyleConfig:   req.StyleConfig,
		CreatedAt:     time.Now(),
	}

	if err := tx.Create(&version).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "保存版本失败",
			"error":   err.Error(),
		})
		return
	}

	resume.Content = req.Content
	resume.StyleConfig = req.StyleConfig
	resume.CurrentVersion = newVersion
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
		"message": "更新成功",
		"data":    resume,
	})
}

func (rc *ResumeController) DeleteResume(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的简历ID",
		})
		return
	}

	var resume models.Resume
	if err := db.DB.First(&resume, uint(id)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "简历不存在",
		})
		return
	}

	tx := db.DB.Begin()

	if err := tx.Where("resume_id = ?", uint(id)).Delete(&models.Version{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "删除版本记录失败",
			"error":   err.Error(),
		})
		return
	}

	if err := tx.Delete(&resume).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "删除简历失败",
			"error":   err.Error(),
		})
		return
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "删除成功",
	})
}
