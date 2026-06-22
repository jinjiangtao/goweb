package handlers

import (
	"errors"
	"net/http"
	"time"
	"weixiu/config"
	"weixiu/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AssignOrderRequest struct {
	TechnicianID uint `json:"technician_id" binding:"required"`
}

type ProcessOrderRequest struct {
	Content string `json:"content" binding:"required"`
	Images  string `json:"images"`
	Status  string `json:"status"`
}

type RejectOrderRequest struct {
	Reason string `json:"reason" binding:"required"`
}

type EvaluateOrderRequest struct {
	Rating  int    `json:"rating" binding:"required,min=1,max=5"`
	Content string `json:"content"`
}

func AssignWorkOrder(c *gin.Context) {
	id := c.Param("id")
	var req AssignOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误", "detail": err.Error()})
		return
	}

	var tech models.User
	if err := config.DB.Where("id = ? AND role = ?", req.TechnicianID, models.RoleTechnician).First(&tech).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "运维人员不存在"})
		return
	}

	result := config.DB.Model(&models.WorkOrder{}).
		Where("id = ? AND status IN ?", id, []models.OrderStatus{models.StatusPending, models.StatusRejected}).
		Updates(map[string]interface{}{
			"technician_id": req.TechnicianID,
			"status":        models.StatusAssigned,
			"assign_time":   time.Now(),
		})

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "派单失败"})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "工单状态不允许派单或工单不存在"})
		return
	}

	var order models.WorkOrder
	config.DB.Preload("Technician").First(&order, id)
	c.JSON(http.StatusOK, gin.H{
		"message": "派单成功",
		"order":   order,
	})
}

func AutoAssignWorkOrder(c *gin.Context) {
	id := c.Param("id")

	var order models.WorkOrder
	if err := config.DB.First(&order, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "工单不存在"})
		return
	}

	if order.Status != models.StatusPending && order.Status != models.StatusRejected {
		c.JSON(http.StatusBadRequest, gin.H{"error": "工单状态不允许自动派单"})
		return
	}

	var technicians []models.User
	query := config.DB.Where("role = ?", models.RoleTechnician)
	if order.Area != "" {
		query = query.Where("area = ? OR area = ?", order.Area, "总部")
	}
	query.Find(&technicians)

	if len(technicians) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "未找到可用的运维人员"})
		return
	}

	type techLoad struct {
		ID    uint
		Count int
	}
	var loads []techLoad
	config.DB.Table("work_orders").
		Select("technician_id as id, COUNT(*) as count").
		Where("technician_id IN ? AND status IN ?", func() []uint {
			ids := make([]uint, len(technicians))
			for i, t := range technicians {
				ids[i] = t.ID
			}
			return ids
		}(), []models.OrderStatus{models.StatusAssigned, models.StatusProcessing}).
		Group("technician_id").
		Scan(&loads)

	loadMap := make(map[uint]int)
	for _, t := range technicians {
		loadMap[t.ID] = 0
	}
	for _, l := range loads {
		loadMap[l.ID] = l.Count
	}

	var bestTech models.User
	minLoad := -1
	for _, t := range technicians {
		if minLoad == -1 || loadMap[t.ID] < minLoad {
			minLoad = loadMap[t.ID]
			bestTech = t
		}
	}

	now := time.Now()
	config.DB.Model(&order).Updates(map[string]interface{}{
		"technician_id": bestTech.ID,
		"status":        models.StatusAssigned,
		"assign_time":   now,
	})
	order.Technician = &bestTech

	c.JSON(http.StatusOK, gin.H{
		"message": "自动派单成功",
		"order":   order,
	})
}

func StartProcessOrder(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("user_id")

	result := config.DB.Model(&models.WorkOrder{}).
		Where("id = ? AND technician_id = ? AND status = ?", id, userID, models.StatusAssigned).
		Updates(map[string]interface{}{
			"status":             models.StatusProcessing,
			"start_process_time": time.Now(),
		})

	if result.RowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "工单状态不允许开始处理或权限不足"})
		return
	}

	var order models.WorkOrder
	config.DB.First(&order, id)
	c.JSON(http.StatusOK, gin.H{"message": "开始处理工单", "order": order})
}

func ProcessWorkOrder(c *gin.Context) {
	id := c.Param("id")
	var req ProcessOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误", "detail": err.Error()})
		return
	}

	userID, _ := c.Get("user_id")
	uid := userID.(uint)

	var order models.WorkOrder
	if err := config.DB.First(&order, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "工单不存在"})
		return
	}

	if order.TechnicianID == nil || *order.TechnicianID != uid {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权处理此工单"})
		return
	}

	if order.Status != models.StatusProcessing {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请先开始处理工单"})
		return
	}

	newStatus := models.OrderStatus(req.Status)
	if newStatus == "" {
		newStatus = models.StatusProcessing
	}

	log := models.RepairLog{
		OrderID:      order.ID,
		TechnicianID: uid,
		Content:      req.Content,
		Images:       req.Images,
		Status:       newStatus,
	}
	config.DB.Create(&log)

	updates := map[string]interface{}{
		"status": newStatus,
	}
	if newStatus == models.StatusCompleted {
		now := time.Now()
		updates["complete_time"] = now
	}
	config.DB.Model(&order).Updates(updates)

	c.JSON(http.StatusOK, gin.H{
		"message": "处理记录已保存",
		"log":     log,
	})
}

func CompleteWorkOrder(c *gin.Context) {
	id := c.Param("id")
	var req ProcessOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误", "detail": err.Error()})
		return
	}

	userID, _ := c.Get("user_id")
	uid := userID.(uint)

	var order models.WorkOrder
	if err := config.DB.First(&order, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "工单不存在"})
		return
	}

	if order.TechnicianID == nil || *order.TechnicianID != uid {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权处理此工单"})
		return
	}

	if order.Status != models.StatusProcessing {
		c.JSON(http.StatusBadRequest, gin.H{"error": "工单当前状态不允许完成"})
		return
	}

	now := time.Now()
	log := models.RepairLog{
		OrderID:      order.ID,
		TechnicianID: uid,
		Content:      req.Content,
		Images:       req.Images,
		Status:       models.StatusCompleted,
	}
	config.DB.Create(&log)

	config.DB.Model(&order).Updates(map[string]interface{}{
		"status":        models.StatusCompleted,
		"complete_time": now,
	})

	c.JSON(http.StatusOK, gin.H{"message": "工单已完成"})
}

func RejectWorkOrder(c *gin.Context) {
	id := c.Param("id")
	var req RejectOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误", "detail": err.Error()})
		return
	}

	result := config.DB.Model(&models.WorkOrder{}).
		Where("id = ? AND status != ? AND status != ?", id, models.StatusCompleted, models.StatusCancelled).
		Updates(map[string]interface{}{
			"status":        models.StatusRejected,
			"reject_reason": req.Reason,
		})

	if result.RowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "工单状态不允许驳回或工单不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "工单已驳回"})
}

func CancelWorkOrder(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("user_id")

	result := config.DB.Model(&models.WorkOrder{}).
		Where("id = ? AND user_id = ? AND status IN ?", id, userID, []models.OrderStatus{models.StatusPending, models.StatusAssigned}).
		Update("status", models.StatusCancelled)

	if result.RowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "工单状态不允许取消或权限不足"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "工单已取消"})
}

func EvaluateWorkOrder(c *gin.Context) {
	id := c.Param("id")
	var req EvaluateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误", "detail": err.Error()})
		return
	}

	userID, _ := c.Get("user_id")

	var order models.WorkOrder
	if err := config.DB.First(&order, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "工单不存在"})
		return
	}

	if order.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权评价此工单"})
		return
	}

	if order.Status != models.StatusCompleted {
		c.JSON(http.StatusBadRequest, gin.H{"error": "只能评价已完成的工单"})
		return
	}

	var existing models.OrderEvaluation
	err := config.DB.Where("order_id = ?", order.ID).First(&existing).Error
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "该工单已评价"})
		return
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
		return
	}

	evaluation := models.OrderEvaluation{
		OrderID: order.ID,
		UserID:  userID.(uint),
		Rating:  req.Rating,
		Content: req.Content,
	}
	config.DB.Create(&evaluation)

	c.JSON(http.StatusOK, gin.H{
		"message":    "评价成功",
		"evaluation": evaluation,
	})
}

func CheckOverdueOrders() {
	var orders []models.WorkOrder
	config.DB.Where("status IN ? AND is_overdue = ? AND expect_time < ?",
		[]models.OrderStatus{models.StatusPending, models.StatusAssigned, models.StatusProcessing},
		false,
		time.Now(),
	).Find(&orders)

	for _, order := range orders {
		config.DB.Model(&order).Update("is_overdue", true)
	}
}
