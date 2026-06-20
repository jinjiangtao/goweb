package handlers

import (
	"gongshi/database"
	"gongshi/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetProjects(c *gin.Context) {
	var projects []models.Project
	status := c.Query("status")
	query := database.DB
	if status != "" {
		query = query.Where("status = ?", status)
	}
	query.Find(&projects)
	c.JSON(http.StatusOK, projects)
}

func CreateProject(c *gin.Context) {
	var project models.Project
	if err := c.ShouldBindJSON(&project); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误: " + err.Error()})
		return
	}

	if project.Name == "" || project.Code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "项目名称和编码不能为空"})
		return
	}

	var existing models.Project
	if database.DB.Where("code = ?", project.Code).First(&existing).Error == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "项目编码已存在"})
		return
	}

	if project.Status == "" {
		project.Status = "active"
	}

	if err := database.DB.Create(&project).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建失败: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, project)
}

func UpdateProject(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	var project models.Project
	if err := database.DB.First(&project, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "项目不存在"})
		return
	}

	var input models.Project
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	database.DB.Model(&project).Updates(input)
	c.JSON(http.StatusOK, project)
}

func DeleteProject(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	var count int64
	database.DB.Model(&models.WorkRecord{}).Where("project_id = ?", id).Count(&count)
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "该项目下存在工时记录，无法删除"})
		return
	}

	result := database.DB.Delete(&models.Project{}, id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}
