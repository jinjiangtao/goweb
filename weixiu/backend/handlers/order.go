package handlers

import (
	"net/http"
	"strconv"
	"time"
	"weixiu/config"
	"weixiu/models"

	"github.com/gin-gonic/gin"
)

type CreateOrderRequest struct {
	Title        string               `json:"title" binding:"required,max=200"`
	Description  string               `json:"description" binding:"required,max=1000"`
	FaultImages  string               `json:"fault_images"`
	DeviceType   models.DeviceType    `json:"device_type" binding:"required"`
	DeviceName   string               `json:"device_name"`
	DeviceID     *uint                `json:"device_id"`
	Area         string               `json:"area" binding:"required"`
	Location     string               `json:"location"`
	ContactName  string               `json:"contact_name"`
	ContactPhone string               `json:"contact_phone"`
	Priority     models.OrderPriority `json:"priority"`
}

func GenerateOrderNo() string {
	now := time.Now()
	return "WO" + now.Format("20060102150405") + strconv.Itoa(int(now.UnixNano()%1000))
}

func CreateWorkOrder(c *gin.Context) {
	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误", "detail": err.Error()})
		return
	}

	userID, _ := c.Get("user_id")

	order := models.WorkOrder{
		OrderNo:      GenerateOrderNo(),
		Title:        req.Title,
		Description:  req.Description,
		FaultImages:  req.FaultImages,
		DeviceType:   req.DeviceType,
		DeviceName:   req.DeviceName,
		DeviceID:     req.DeviceID,
		Area:         req.Area,
		Location:     req.Location,
		ContactName:  req.ContactName,
		ContactPhone: req.ContactPhone,
		Priority:     req.Priority,
		Status:       models.StatusPending,
		UserID:       userID.(uint),
	}

	if order.Priority == "" {
		order.Priority = models.PriorityMedium
	}

	var expectHours time.Duration
	switch order.Priority {
	case models.PriorityUrgent:
		expectHours = 2
	case models.PriorityHigh:
		expectHours = 8
	case models.PriorityMedium:
		expectHours = 24
	default:
		expectHours = 48
	}
	expectTime := time.Now().Add(expectHours * time.Hour)
	order.ExpectTime = &expectTime

	if err := config.DB.Create(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "工单创建失败"})
		return
	}

	config.DB.Preload("User").First(&order, order.ID)

	c.JSON(http.StatusOK, gin.H{
		"message": "工单提交成功",
		"order":   order,
	})
}

func GetWorkOrders(c *gin.Context) {
	status := c.Query("status")
	deviceType := c.Query("device_type")
	area := c.Query("area")
	priority := c.Query("priority")
	userIDStr := c.Query("user_id")
	technicianIDStr := c.Query("technician_id")
	keyword := c.Query("keyword")

	userID, _ := c.Get("user_id")
	userRole, _ := c.Get("role")

	var orders []models.WorkOrder
	query := config.DB.Preload("User").Preload("Technician").Order("created_at DESC")

	if keyword != "" {
		query = query.Where("title LIKE ? OR description LIKE ? OR order_no LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if deviceType != "" {
		query = query.Where("device_type = ?", deviceType)
	}
	if area != "" {
		query = query.Where("area = ?", area)
	}
	if priority != "" {
		query = query.Where("priority = ?", priority)
	}
	if userIDStr != "" {
		query = query.Where("user_id = ?", userIDStr)
	}
	if technicianIDStr != "" {
		query = query.Where("technician_id = ?", technicianIDStr)
	}

	if userRole == "user" {
		query = query.Where("user_id = ?", userID)
	} else if userRole == "technician" {
		query = query.Where("technician_id = ? OR status IN ?", userID, []models.OrderStatus{models.StatusPending, models.StatusAssigned})
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var total int64
	query.Model(&models.WorkOrder{}).Count(&total)
	query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&orders)

	c.JSON(http.StatusOK, gin.H{
		"list":      orders,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func GetWorkOrder(c *gin.Context) {
	id := c.Param("id")
	var order models.WorkOrder
	if err := config.DB.Preload("User").Preload("Technician").First(&order, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "工单不存在"})
		return
	}

	var logs []models.RepairLog
	config.DB.Preload("Technician").Where("order_id = ?", order.ID).Order("created_at ASC").Find(&logs)

	var evaluation *models.OrderEvaluation
	config.DB.Where("order_id = ?", order.ID).First(&evaluation)

	c.JSON(http.StatusOK, gin.H{
		"order":      order,
		"logs":       logs,
		"evaluation": evaluation,
	})
}
