package handlers

import (
	"gongshi/database"
	"gongshi/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CreateWorkRecordRequest struct {
	UserID    uint    `json:"user_id"`
	ProjectID uint    `json:"project_id" binding:"required"`
	WorkDate  string  `json:"work_date" binding:"required"`
	Hours     float64 `json:"hours" binding:"required"`
	Content   string  `json:"content" binding:"required"`
}

func validateWorkRecord(req CreateWorkRecordRequest, excludeID *uint) (string, bool) {
	if req.Hours <= 0 || req.Hours > 24 {
		return "工时必须在0-24小时之间", false
	}

	if req.WorkDate == "" {
		return "工作日期不能为空", false
	}

	_, err := time.Parse("2006-01-02", req.WorkDate)
	if err != nil {
		return "日期格式不正确，应为YYYY-MM-DD", false
	}

	var totalHours float64
	query := database.DB.Model(&models.WorkRecord{}).
		Where("user_id = ? AND work_date = ?", req.UserID, req.WorkDate)
	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}
	query.Select("COALESCE(SUM(hours), 0)").Scan(&totalHours)

	if totalHours+req.Hours > 24 {
		return "该日累计工时不能超过24小时，当前已有" + strconv.FormatFloat(totalHours, 'f', 1, 64) + "小时", false
	}

	var dup models.WorkRecord
	dupQuery := database.DB.Where("user_id = ? AND work_date = ? AND project_id = ?",
		req.UserID, req.WorkDate, req.ProjectID)
	if excludeID != nil {
		dupQuery = dupQuery.Where("id != ?", *excludeID)
	}
	if dupQuery.First(&dup).Error == nil {
		return "该日期下该项目已存在工时记录，请修改原记录", false
	}

	return "", true
}

func GetWorkRecords(c *gin.Context) {
	userID, _ := strconv.Atoi(c.Query("user_id"))
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	status := c.Query("status")
	projectID, _ := strconv.Atoi(c.Query("project_id"))

	role := c.GetHeader("X-User-Role")
	currentUserID, _ := strconv.Atoi(c.GetHeader("X-User-ID"))
	if currentUserID == 0 {
		currentUserID = 2
	}

	var records []models.WorkRecord
	query := database.DB.Preload("User").Preload("Project").Preload("Reviewer")

	if role != "admin" || userID > 0 {
		if userID > 0 {
			query = query.Where("work_records.user_id = ?", userID)
		} else {
			query = query.Where("work_records.user_id = ?", currentUserID)
		}
	}

	if startDate != "" {
		query = query.Where("work_date >= ?", startDate)
	}
	if endDate != "" {
		query = query.Where("work_date <= ?", endDate)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if projectID > 0 {
		query = query.Where("project_id = ?", projectID)
	}

	query.Order("work_date DESC, created_at DESC").Find(&records)
	c.JSON(http.StatusOK, records)
}

func CreateWorkRecord(c *gin.Context) {
	var req CreateWorkRecordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误: " + err.Error()})
		return
	}

	currentUserID, _ := strconv.Atoi(c.GetHeader("X-User-ID"))
	if currentUserID == 0 {
		currentUserID = 2
	}
	req.UserID = uint(currentUserID)

	if msg, ok := validateWorkRecord(req, nil); !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": msg, "type": "validation"})
		return
	}

	var project models.Project
	if err := database.DB.First(&project, req.ProjectID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "项目不存在"})
		return
	}

	record := models.WorkRecord{
		UserID:    req.UserID,
		ProjectID: req.ProjectID,
		WorkDate:  req.WorkDate,
		Hours:     req.Hours,
		Content:   req.Content,
		Status:    "pending",
	}

	if err := database.DB.Create(&record).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建失败: " + err.Error()})
		return
	}

	database.DB.Preload("User").Preload("Project").First(&record, record.ID)
	c.JSON(http.StatusOK, record)
}

func BatchCreateWorkRecords(c *gin.Context) {
	var records []CreateWorkRecordRequest
	if err := c.ShouldBindJSON(&records); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	currentUserID, _ := strconv.Atoi(c.GetHeader("X-User-ID"))
	if currentUserID == 0 {
		currentUserID = 2
	}

	created := []models.WorkRecord{}
	errors := []gin.H{}

	for i, req := range records {
		req.UserID = uint(currentUserID)
		if msg, ok := validateWorkRecord(req, nil); !ok {
			errors = append(errors, gin.H{"index": i, "error": msg})
			continue
		}

		record := models.WorkRecord{
			UserID:    req.UserID,
			ProjectID: req.ProjectID,
			WorkDate:  req.WorkDate,
			Hours:     req.Hours,
			Content:   req.Content,
			Status:    "pending",
		}
		if err := database.DB.Create(&record).Error; err != nil {
			errors = append(errors, gin.H{"index": i, "error": err.Error()})
			continue
		}
		database.DB.Preload("User").Preload("Project").First(&record, record.ID)
		created = append(created, record)
	}

	c.JSON(http.StatusOK, gin.H{
		"created": created,
		"errors":  errors,
	})
}

func UpdateWorkRecord(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	var record models.WorkRecord
	if err := database.DB.First(&record, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "记录不存在"})
		return
	}

	currentUserID, _ := strconv.Atoi(c.GetHeader("X-User-ID"))
	if currentUserID == 0 {
		currentUserID = 2
	}

	if record.Status == "approved" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "已审核通过的记录无法修改"})
		return
	}

	if record.UserID != uint(currentUserID) && c.GetHeader("X-User-Role") != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权修改"})
		return
	}

	var req CreateWorkRecordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	req.UserID = record.UserID

	recordID := uint(id)
	if msg, ok := validateWorkRecord(req, &recordID); !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": msg, "type": "validation"})
		return
	}

	record.ProjectID = req.ProjectID
	record.WorkDate = req.WorkDate
	record.Hours = req.Hours
	record.Content = req.Content
	record.Status = "pending"
	record.RejectReason = ""
	record.ReviewerID = nil
	record.ReviewedAt = nil

	database.DB.Save(&record)
	database.DB.Preload("User").Preload("Project").First(&record, record.ID)
	c.JSON(http.StatusOK, record)
}

func DeleteWorkRecord(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	var record models.WorkRecord
	if err := database.DB.First(&record, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "记录不存在"})
		return
	}

	currentUserID, _ := strconv.Atoi(c.GetHeader("X-User-ID"))
	if currentUserID == 0 {
		currentUserID = 2
	}

	if record.Status == "approved" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "已审核通过的记录无法删除"})
		return
	}

	if record.UserID != uint(currentUserID) && c.GetHeader("X-User-Role") != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权删除"})
		return
	}

	database.DB.Delete(&record)
	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

type ReviewRequest struct {
	RejectReason string `json:"reject_reason"`
}

func ApproveWorkRecord(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	var record models.WorkRecord
	if err := database.DB.First(&record, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "记录不存在"})
		return
	}

	if record.Status != "pending" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "该记录当前状态不可审核"})
		return
	}

	reviewerID, _ := strconv.Atoi(c.GetHeader("X-User-ID"))
	if reviewerID == 0 {
		reviewerID = 1
	}
	reviewerIDUint := uint(reviewerID)
	now := models.LocalTime(time.Now())

	record.Status = "approved"
	record.ReviewerID = &reviewerIDUint
	record.ReviewedAt = &now
	record.RejectReason = ""

	database.DB.Save(&record)
	database.DB.Preload("User").Preload("Project").Preload("Reviewer").First(&record, record.ID)
	c.JSON(http.StatusOK, record)
}

func RejectWorkRecord(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	var record models.WorkRecord
	if err := database.DB.First(&record, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "记录不存在"})
		return
	}

	if record.Status != "pending" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "该记录当前状态不可审核"})
		return
	}

	var req ReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	if req.RejectReason == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "驳回原因不能为空"})
		return
	}

	reviewerID, _ := strconv.Atoi(c.GetHeader("X-User-ID"))
	if reviewerID == 0 {
		reviewerID = 1
	}
	reviewerIDUint := uint(reviewerID)
	now := models.LocalTime(time.Now())

	record.Status = "rejected"
	record.ReviewerID = &reviewerIDUint
	record.ReviewedAt = &now
	record.RejectReason = req.RejectReason

	database.DB.Save(&record)
	database.DB.Preload("User").Preload("Project").Preload("Reviewer").First(&record, record.ID)
	c.JSON(http.StatusOK, record)
}

func BatchApprove(c *gin.Context) {
	var body struct {
		IDs []uint `json:"ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	reviewerID, _ := strconv.Atoi(c.GetHeader("X-User-ID"))
	if reviewerID == 0 {
		reviewerID = 1
	}
	reviewerIDUint := uint(reviewerID)
	now := models.LocalTime(time.Now())

	result := database.DB.Model(&models.WorkRecord{}).
		Where("id IN ? AND status = ?", body.IDs, "pending").
		Updates(map[string]interface{}{
			"status":        "approved",
			"reviewer_id":   reviewerIDUint,
			"reviewed_at":   now,
			"reject_reason": "",
		})

	c.JSON(http.StatusOK, gin.H{
		"message": "批量审核成功",
		"count":   result.RowsAffected,
	})
}

func GetDailyHoursSummary(c *gin.Context) {
	userID, _ := strconv.Atoi(c.Query("user_id"))
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	currentUserID, _ := strconv.Atoi(c.GetHeader("X-User-ID"))
	if currentUserID == 0 {
		currentUserID = 2
	}
	if userID == 0 {
		userID = currentUserID
	}

	type DailyResult struct {
		WorkDate string  `json:"work_date"`
		Total    float64 `json:"total"`
		Approved float64 `json:"approved"`
		Pending  float64 `json:"pending"`
		Rejected float64 `json:"rejected"`
	}

	var results []DailyResult
	query := database.DB.Model(&models.WorkRecord{}).
		Select("work_date, "+
			"COALESCE(SUM(hours),0) as total, "+
			"COALESCE(SUM(CASE WHEN status='approved' THEN hours ELSE 0 END),0) as approved, "+
			"COALESCE(SUM(CASE WHEN status='pending' THEN hours ELSE 0 END),0) as pending, "+
			"COALESCE(SUM(CASE WHEN status='rejected' THEN hours ELSE 0 END),0) as rejected").
		Where("user_id = ?", userID).
		Group("work_date")

	if startDate != "" {
		query = query.Where("work_date >= ?", startDate)
	}
	if endDate != "" {
		query = query.Where("work_date <= ?", endDate)
	}

	query.Order("work_date").Scan(&results)
	if results == nil {
		results = []DailyResult{}
	}
	c.JSON(http.StatusOK, results)
}

func GetProjectHoursSummary(c *gin.Context) {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	userID, _ := strconv.Atoi(c.Query("user_id"))

	type ProjectResult struct {
		ProjectID   uint    `json:"project_id"`
		ProjectName string  `json:"project_name"`
		ProjectCode string  `json:"project_code"`
		Total       float64 `json:"total"`
		Approved    float64 `json:"approved"`
		Pending     float64 `json:"pending"`
		Rejected    float64 `json:"rejected"`
		RecordCount int64   `json:"record_count"`
	}

	var results []ProjectResult
	query := database.DB.Model(&models.WorkRecord{}).
		Select("work_records.project_id, projects.name as project_name, projects.code as project_code, " +
			"COALESCE(SUM(work_records.hours),0) as total, " +
			"COALESCE(SUM(CASE WHEN work_records.status='approved' THEN work_records.hours ELSE 0 END),0) as approved, " +
			"COALESCE(SUM(CASE WHEN work_records.status='pending' THEN work_records.hours ELSE 0 END),0) as pending, " +
			"COALESCE(SUM(CASE WHEN work_records.status='rejected' THEN work_records.hours ELSE 0 END),0) as rejected, " +
			"COUNT(work_records.id) as record_count").
		Joins("LEFT JOIN projects ON projects.id = work_records.project_id").
		Group("work_records.project_id, projects.name, projects.code")

	if startDate != "" {
		query = query.Where("work_records.work_date >= ?", startDate)
	}
	if endDate != "" {
		query = query.Where("work_records.work_date <= ?", endDate)
	}
	if userID > 0 {
		query = query.Where("work_records.user_id = ?", userID)
	}

	query.Order("total DESC").Scan(&results)
	if results == nil {
		results = []ProjectResult{}
	}
	c.JSON(http.StatusOK, results)
}

func GetOverallStats(c *gin.Context) {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	userID, _ := strconv.Atoi(c.Query("user_id"))

	query := database.DB.Model(&models.WorkRecord{})
	if startDate != "" {
		query = query.Where("work_date >= ?", startDate)
	}
	if endDate != "" {
		query = query.Where("work_date <= ?", endDate)
	}
	if userID > 0 {
		query = query.Where("user_id = ?", userID)
	}

	type StatusCount struct {
		Status string  `json:"status"`
		Hours  float64 `json:"hours"`
		Count  int64   `json:"count"`
	}
	var byStatus []StatusCount
	database.DB.Table("(?) as t", query.Select("status, hours")).
		Select("status, COALESCE(SUM(hours),0) as hours, COUNT(*) as count").
		Group("status").Scan(&byStatus)

	type UserStat struct {
		UserID   uint    `json:"user_id"`
		UserName string  `json:"user_name"`
		Total    float64 `json:"total"`
		Approved float64 `json:"approved"`
	}
	var byUser []UserStat
	uQuery := database.DB.Model(&models.WorkRecord{}).
		Select("user_id, users.name as user_name, " +
			"COALESCE(SUM(hours),0) as total, " +
			"COALESCE(SUM(CASE WHEN status='approved' THEN hours ELSE 0 END),0) as approved").
		Joins("LEFT JOIN users ON users.id = work_records.user_id").
		Group("user_id, users.name")
	if startDate != "" {
		uQuery = uQuery.Where("work_date >= ?", startDate)
	}
	if endDate != "" {
		uQuery = uQuery.Where("work_date <= ?", endDate)
	}
	uQuery.Order("total DESC").Limit(10).Scan(&byUser)

	var totalRecords int64
	var totalHours float64
	query.Session(&gorm.Session{}).Count(&totalRecords)
	database.DB.Table("(?) as t", query.Select("hours")).
		Select("COALESCE(SUM(hours),0)").Scan(&totalHours)

	if byStatus == nil {
		byStatus = []StatusCount{}
	}
	if byUser == nil {
		byUser = []UserStat{}
	}

	c.JSON(http.StatusOK, gin.H{
		"total_records": totalRecords,
		"total_hours":   totalHours,
		"by_status":     byStatus,
		"top_users":     byUser,
	})
}
