package handlers

import (
	"net/http"
	"shuiti/database"
	"shuiti/models"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type QuestionListResponse struct {
	Total    int64             `json:"total"`
	List     []models.Question `json:"list"`
	Page     int               `json:"page"`
	PageSize int               `json:"page_size"`
}

func GetQuestions(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	qType := c.Query("type")
	category := c.Query("category")
	keyword := c.Query("keyword")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	query := database.DB.Model(&models.Question{})

	if qType != "" {
		query = query.Where("type = ?", qType)
	}
	if category != "" {
		query = query.Where("category LIKE ?", "%"+category+"%")
	}
	if keyword != "" {
		query = query.Where("content LIKE ?", "%"+keyword+"%")
	}

	var total int64
	query.Count(&total)

	var questions []models.Question
	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&questions).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, QuestionListResponse{
		Total:    total,
		List:     questions,
		Page:     page,
		PageSize: pageSize,
	})
}

func GetQuestion(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var question models.Question
	err = database.DB.First(&question, uint(id)).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Question not found"})
		return
	}

	c.JSON(http.StatusOK, question)
}

func CreateQuestion(c *gin.Context) {
	var question models.Question
	if err := c.ShouldBindJSON(&question); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if question.Type == "" {
		question.Type = models.SingleChoice
	}

	err := database.DB.Create(&question).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, question)
}

func BatchCreateQuestions(c *gin.Context) {
	var questions []models.Question
	if err := c.ShouldBindJSON(&questions); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for i := range questions {
		if questions[i].Type == "" {
			questions[i].Type = models.SingleChoice
		}
	}

	err := database.DB.Create(&questions).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"count":   len(questions),
		"data":    questions,
	})
}

func UpdateQuestion(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var question models.Question
	err = database.DB.First(&question, uint(id)).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Question not found"})
		return
	}

	var updateData models.Question
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = database.DB.Model(&question).Updates(updateData).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	database.DB.First(&question, uint(id))
	c.JSON(http.StatusOK, question)
}

func DeleteQuestion(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	result := database.DB.Delete(&models.Question{}, uint(id))
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Question not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

func BatchDeleteQuestions(c *gin.Context) {
	var req struct {
		IDs []uint `json:"ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := database.DB.Delete(&models.Question{}, req.IDs)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"deleted": result.RowsAffected,
	})
}

func GetCategories(c *gin.Context) {
	var categories []string
	database.DB.Model(&models.Question{}).Distinct("category").Where("category != ''").Pluck("category", &categories)

	var typeStats []struct {
		Type  string `json:"type"`
		Count int64  `json:"count"`
	}
	database.DB.Model(&models.Question{}).Select("type, count(*) as count").Group("type").Scan(&typeStats)

	var total int64
	database.DB.Model(&models.Question{}).Count(&total)

	c.JSON(http.StatusOK, gin.H{
		"categories": categories,
		"type_stats": typeStats,
		"total":      total,
	})
}

func GetQuestionsNoPaging(c *gin.Context) {
	qType := c.Query("type")
	category := c.Query("category")
	ids := c.Query("ids")

	query := database.DB.Model(&models.Question{})

	if ids != "" {
		idList := strings.Split(ids, ",")
		var uintIds []uint
		for _, idStr := range idList {
			id, err := strconv.ParseUint(strings.TrimSpace(idStr), 10, 32)
			if err == nil {
				uintIds = append(uintIds, uint(id))
			}
		}
		if len(uintIds) > 0 {
			query = query.Where("id IN ?", uintIds)
		}
	} else {
		if qType != "" {
			query = query.Where("type = ?", qType)
		}
		if category != "" {
			query = query.Where("category LIKE ?", "%"+category+"%")
		}
	}

	var questions []models.Question
	err := query.Order("id DESC").Find(&questions).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, questions)
}
