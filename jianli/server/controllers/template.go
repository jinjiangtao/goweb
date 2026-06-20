package controllers

import (
	"jianli-server/db"
	"jianli-server/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TemplateController struct{}

func NewTemplateController() *TemplateController {
	return &TemplateController{}
}

func (tc *TemplateController) GetTemplates(c *gin.Context) {
	var templates []models.Template
	if err := db.DB.Find(&templates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取模板列表失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    templates,
	})
}

func (tc *TemplateController) GetTemplate(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的模板ID",
		})
		return
	}

	var template models.Template
	if err := db.DB.First(&template, uint(id)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "模板不存在",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    template,
	})
}
