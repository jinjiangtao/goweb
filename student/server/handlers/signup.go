package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"student-signup-server/models"
)

type SignupRequest struct {
	Name   string `json:"name" binding:"required"`
	Phone  string `json:"phone" binding:"required"`
	Age    int    `json:"age" binding:"required"`
	Hukou  string `json:"hukou" binding:"required"`
	School string `json:"school" binding:"required"`
}

func Signup(c *gin.Context) {
	var req SignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	signup := models.Signup{
		Name:   req.Name,
		Phone:  req.Phone,
		Age:    req.Age,
		Hukou:  req.Hukou,
		School: req.School,
		Status: "pending",
	}

	err := models.CreateSignup(&signup)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "报名失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "报名成功"})
}

func GetSignups(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	keyword := c.Query("keyword")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	result, err := models.GetSignups(page, pageSize, keyword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取列表失败"})
		return
	}

	c.JSON(http.StatusOK, result)
}

type UpdateStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

func UpdateSignupStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID格式错误"})
		return
	}

	var req UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	if req.Status != "pending" && req.Status != "approved" && req.Status != "rejected" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "状态值无效"})
		return
	}

	err = models.UpdateSignupStatus(uint(id), req.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}

func GetStats(c *gin.Context) {
	stats, err := models.GetStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取统计失败"})
		return
	}

	c.JSON(http.StatusOK, stats)
}
