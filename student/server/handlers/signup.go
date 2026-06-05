package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"student-signup-server/models"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
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
	status := c.Query("status")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	result, err := models.GetSignups(page, pageSize, keyword, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取列表失败"})
		return
	}

	c.JSON(http.StatusOK, result)
}

type AdminSignupRequest struct {
	Name   string `json:"name" binding:"required"`
	Phone  string `json:"phone" binding:"required"`
	Age    int    `json:"age" binding:"required"`
	Hukou  string `json:"hukou" binding:"required"`
	School string `json:"school" binding:"required"`
	Status string `json:"status"`
}

func AdminCreateSignup(c *gin.Context) {
	var req AdminSignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	if req.Status == "" {
		req.Status = "pending"
	}
	if req.Status != "pending" && req.Status != "approved" && req.Status != "rejected" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "状态值无效"})
		return
	}

	signup := models.Signup{
		Name:   req.Name,
		Phone:  req.Phone,
		Age:    req.Age,
		Hukou:  req.Hukou,
		School: req.School,
		Status: req.Status,
	}

	err := models.CreateSignup(&signup)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "创建成功", "data": signup})
}

type AdminUpdateSignupRequest struct {
	Name   string `json:"name" binding:"required"`
	Phone  string `json:"phone" binding:"required"`
	Age    int    `json:"age" binding:"required"`
	Hukou  string `json:"hukou" binding:"required"`
	School string `json:"school" binding:"required"`
}

func AdminUpdateSignup(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID格式错误"})
		return
	}

	var req AdminUpdateSignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	signup, err := models.GetSignupByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "记录不存在"})
		return
	}

	signup.Name = req.Name
	signup.Phone = req.Phone
	signup.Age = req.Age
	signup.Hukou = req.Hukou
	signup.School = req.School

	err = models.UpdateSignup(signup)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "更新成功", "data": signup})
}

func ExportSignups(c *gin.Context) {
	keyword := c.Query("keyword")
	status := c.Query("status")

	signups, err := models.GetAllSignupsForExport(keyword, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取数据失败"})
		return
	}

	f := excelize.NewFile()
	index, err := f.NewSheet("报名记录")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建Excel失败"})
		return
	}
	f.SetCellValue("报名记录", "A1", "姓名")
	f.SetCellValue("报名记录", "B1", "手机号")
	f.SetCellValue("报名记录", "C1", "年龄")
	f.SetCellValue("报名记录", "D1", "户口地址")
	f.SetCellValue("报名记录", "E1", "学校")
	f.SetCellValue("报名记录", "F1", "状态")
	f.SetCellValue("报名记录", "G1", "提交时间")

	for i, signup := range signups {
		row := i + 2
		f.SetCellValue("报名记录", fmt.Sprintf("A%d", row), signup.Name)
		f.SetCellValue("报名记录", fmt.Sprintf("B%d", row), signup.Phone)
		f.SetCellValue("报名记录", fmt.Sprintf("C%d", row), signup.Age)
		f.SetCellValue("报名记录", fmt.Sprintf("D%d", row), signup.Hukou)
		f.SetCellValue("报名记录", fmt.Sprintf("E%d", row), signup.School)
		f.SetCellValue("报名记录", fmt.Sprintf("F%d", row), getStatusText(signup.Status))
		f.SetCellValue("报名记录", fmt.Sprintf("G%d", row), signup.CreatedAt.String())
	}

	f.SetActiveSheet(index)

	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename=signups.xlsx")

	if err := f.Write(c.Writer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "导出失败"})
		return
	}
}

func getStatusText(status string) string {
	switch status {
	case "pending":
		return "报名中"
	case "approved":
		return "报名成功"
	case "rejected":
		return "报名失败"
	default:
		return status
	}
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
