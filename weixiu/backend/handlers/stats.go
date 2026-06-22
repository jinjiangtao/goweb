package handlers

import (
	"net/http"
	"time"
	"weixiu/config"
	"weixiu/models"

	"github.com/gin-gonic/gin"
)

func GetStatistics(c *gin.Context) {
	var totalOrders int64
	var pendingOrders int64
	var processingOrders int64
	var completedOrders int64
	var rejectedOrders int64
	var overdueOrders int64

	config.DB.Model(&models.WorkOrder{}).Count(&totalOrders)
	config.DB.Model(&models.WorkOrder{}).Where("status = ?", models.StatusPending).Count(&pendingOrders)
	config.DB.Model(&models.WorkOrder{}).Where("status = ? OR status = ?", models.StatusAssigned, models.StatusProcessing).Count(&processingOrders)
	config.DB.Model(&models.WorkOrder{}).Where("status = ?", models.StatusCompleted).Count(&completedOrders)
	config.DB.Model(&models.WorkOrder{}).Where("status = ?", models.StatusRejected).Count(&rejectedOrders)
	config.DB.Model(&models.WorkOrder{}).Where("is_overdue = ?", true).Count(&overdueOrders)

	type StatusCount struct {
		Status string `json:"status"`
		Count  int64  `json:"count"`
	}
	var statusStats []StatusCount
	config.DB.Model(&models.WorkOrder{}).
		Select("status, COUNT(*) as count").
		Group("status").
		Scan(&statusStats)

	type TypeCount struct {
		DeviceType string `json:"device_type"`
		Count      int64  `json:"count"`
	}
	var typeStats []TypeCount
	config.DB.Model(&models.WorkOrder{}).
		Select("device_type, COUNT(*) as count").
		Group("device_type").
		Scan(&typeStats)

	type AreaCount struct {
		Area  string `json:"area"`
		Count int64  `json:"count"`
	}
	var areaStats []AreaCount
	config.DB.Model(&models.WorkOrder{}).
		Select("area, COUNT(*) as count").
		Group("area").
		Scan(&areaStats)

	var completedList []models.WorkOrder
	config.DB.Where("status = ? AND complete_time IS NOT NULL AND created_at IS NOT NULL", models.StatusCompleted).
		Find(&completedList)

	var totalDuration float64
	var completedCount int
	for _, o := range completedList {
		if o.CompleteTime != nil && !o.CreatedAt.IsZero() {
			totalDuration += o.CompleteTime.Sub(o.CreatedAt).Hours()
			completedCount++
		}
	}

	avgDuration := 0.0
	if completedCount > 0 {
		avgDuration = totalDuration / float64(completedCount)
	}

	var avgRating float64
	config.DB.Model(&models.OrderEvaluation{}).Select("COALESCE(AVG(rating), 0)").Scan(&avgRating)

	type DailyCount struct {
		Date  string `json:"date"`
		Count int64  `json:"count"`
	}

	now := time.Now()
	var dailyStats []DailyCount
	for i := 6; i >= 0; i-- {
		date := now.AddDate(0, 0, -i)
		start := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
		end := start.AddDate(0, 0, 1)

		var count int64
		config.DB.Model(&models.WorkOrder{}).
			Where("created_at >= ? AND created_at < ?", start, end).
			Count(&count)

		dailyStats = append(dailyStats, DailyCount{
			Date:  start.Format("2006-01-02"),
			Count: count,
		})
	}

	var technicians []models.User
	config.DB.Where("role = ?", models.RoleTechnician).Find(&technicians)

	type TechStat struct {
		ID          uint    `json:"id"`
		RealName    string  `json:"real_name"`
		TotalOrders int64   `json:"total_orders"`
		Completed   int64   `json:"completed"`
		AvgDuration float64 `json:"avg_duration"`
		AvgRating   float64 `json:"avg_rating"`
	}

	var techStats []TechStat
	for _, tech := range technicians {
		var total int64
		var completed int64
		config.DB.Model(&models.WorkOrder{}).Where("technician_id = ?", tech.ID).Count(&total)
		config.DB.Model(&models.WorkOrder{}).Where("technician_id = ? AND status = ?", tech.ID, models.StatusCompleted).Count(&completed)

		var avgDur float64
		rows := config.DB.Model(&models.WorkOrder{}).
			Select("COALESCE(AVG(JULIANDAY(complete_time) - JULIANDAY(created_at)) * 24, 0)").
			Where("technician_id = ? AND status = ? AND complete_time IS NOT NULL", tech.ID, models.StatusCompleted).
			Row()
		rows.Scan(&avgDur)

		var avgRat float64
		config.DB.Table("order_evaluations oe").
			Select("COALESCE(AVG(oe.rating), 0)").
			Joins("JOIN work_orders wo ON wo.id = oe.order_id").
			Where("wo.technician_id = ?", tech.ID).
			Scan(&avgRat)

		techStats = append(techStats, TechStat{
			ID:          tech.ID,
			RealName:    tech.RealName,
			TotalOrders: total,
			Completed:   completed,
			AvgDuration: avgDur,
			AvgRating:   avgRat,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"overview": gin.H{
			"total":        totalOrders,
			"pending":      pendingOrders,
			"processing":   processingOrders,
			"completed":    completedOrders,
			"rejected":     rejectedOrders,
			"overdue":      overdueOrders,
			"avg_duration": avgDuration,
			"avg_rating":   avgRating,
		},
		"status_stats":     statusStats,
		"type_stats":       typeStats,
		"area_stats":       areaStats,
		"daily_stats":      dailyStats,
		"technician_stats": techStats,
	})
}
