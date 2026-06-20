package handlers

import (
	"logistics-server/database"
	"logistics-server/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CreateOrderRequest struct {
	OrderNo     string  `json:"order_no" binding:"required"`
	Sender      string  `json:"sender"`
	Receiver    string  `json:"receiver"`
	Origin      string  `json:"origin"`
	Destination string  `json:"destination"`
	GoodsName   string  `json:"goods_name"`
	Weight      float64 `json:"weight"`
}

func refreshOrderStatus(orderID uint) {
	var order models.LogisticsOrder
	if err := database.DB.First(&order, orderID).Error; err != nil {
		return
	}

	var nodes []models.LogisticsNode
	database.DB.Where("order_id = ?", orderID).Order("occurred_at ASC, id ASC").Find(&nodes)

	if len(nodes) == 0 {
		order.Status = models.OrderStatusPending
		database.DB.Save(&order)
		return
	}

	hasAbnormal := false
	hasDelivered := false
	hasArrived := false
	hasInTransit := false
	hasPicked := false

	for _, node := range nodes {
		if node.IsAbnormal {
			hasAbnormal = true
		}
		switch node.Status {
		case models.NodeStatusDelivered:
			hasDelivered = true
		case models.NodeStatusArrived:
			hasArrived = true
		case models.NodeStatusInTransit:
			hasInTransit = true
		case models.NodeStatusPicked:
			hasPicked = true
		}
	}

	newStatus := models.OrderStatusPending
	if hasDelivered {
		newStatus = models.OrderStatusDelivered
	} else if hasAbnormal {
		newStatus = models.OrderStatusException
	} else if hasArrived {
		newStatus = models.OrderStatusArrived
	} else if hasInTransit {
		newStatus = models.OrderStatusInTransit
	} else if hasPicked {
		newStatus = models.OrderStatusPicked
	}

	if order.Status != newStatus {
		order.Status = newStatus
		database.DB.Save(&order)
	}
}

func CreateOrder(c *gin.Context) {
	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order := models.LogisticsOrder{
		OrderNo:     req.OrderNo,
		Sender:      req.Sender,
		Receiver:    req.Receiver,
		Origin:      req.Origin,
		Destination: req.Destination,
		GoodsName:   req.GoodsName,
		Weight:      req.Weight,
		Status:      models.OrderStatusPending,
	}

	if err := database.DB.Create(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "订单创建成功",
		"data":    order,
	})
}

func GetOrders(c *gin.Context) {
	var orders []models.LogisticsOrder
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	keyword := c.Query("keyword")

	offset := (page - 1) * pageSize

	query := database.DB.Model(&models.LogisticsOrder{})
	if keyword != "" {
		query = query.Where("order_no LIKE ? OR sender LIKE ? OR receiver LIKE ? OR origin LIKE ? OR destination LIKE ? OR goods_name LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	var total int64
	query.Count(&total)

	dbQuery := database.DB.Model(&models.LogisticsOrder{})
	if keyword != "" {
		dbQuery = dbQuery.Where("order_no LIKE ? OR sender LIKE ? OR receiver LIKE ? OR origin LIKE ? OR destination LIKE ? OR goods_name LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	if err := dbQuery.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total": total,
		"page":  page,
		"size":  pageSize,
		"data":  orders,
	})
}

func GetOrder(c *gin.Context) {
	id := c.Param("id")
	var order models.LogisticsOrder

	if err := database.DB.Preload("Nodes", func(db *gorm.DB) *gorm.DB {
		return db.Order("occurred_at ASC, id ASC")
	}).First(&order, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "订单不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": order})
}

func SearchOrder(c *gin.Context) {
	orderNo := c.Query("order_no")
	if orderNo == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请输入订单号"})
		return
	}

	var order models.LogisticsOrder
	if err := database.DB.Where("order_no = ?", orderNo).Preload("Nodes", func(db *gorm.DB) *gorm.DB {
		return db.Order("occurred_at ASC, id ASC")
	}).First(&order).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "订单不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": order})
}

func UpdateOrder(c *gin.Context) {
	id := c.Param("id")
	var order models.LogisticsOrder

	if err := database.DB.First(&order, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "订单不存在"})
		return
	}

	var req struct {
		Sender      string  `json:"sender"`
		Receiver    string  `json:"receiver"`
		Origin      string  `json:"origin"`
		Destination string  `json:"destination"`
		GoodsName   string  `json:"goods_name"`
		Weight      float64 `json:"weight"`
		Status      string  `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Sender != "" {
		order.Sender = req.Sender
	}
	if req.Receiver != "" {
		order.Receiver = req.Receiver
	}
	if req.Origin != "" {
		order.Origin = req.Origin
	}
	if req.Destination != "" {
		order.Destination = req.Destination
	}
	if req.GoodsName != "" {
		order.GoodsName = req.GoodsName
	}
	if req.Weight != 0 {
		order.Weight = req.Weight
	}
	if req.Status != "" {
		order.Status = req.Status
	}

	if err := database.DB.Save(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "订单更新成功",
		"data":    order,
	})
}

func DeleteOrder(c *gin.Context) {
	id := c.Param("id")

	if err := database.DB.Delete(&models.LogisticsOrder{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	database.DB.Where("order_id = ?", id).Delete(&models.LogisticsNode{})

	c.JSON(http.StatusOK, gin.H{"message": "订单删除成功"})
}

type CreateNodeRequest struct {
	NodeName       string `json:"node_name" binding:"required"`
	Status         string `json:"status"`
	Location       string `json:"location"`
	Description    string `json:"description"`
	Operator       string `json:"operator"`
	IsAbnormal     bool   `json:"is_abnormal"`
	AbnormalReason string `json:"abnormal_reason"`
	OccurredAt     string `json:"occurred_at"`
}

func AddNode(c *gin.Context) {
	orderID := c.Param("id")
	var req CreateNodeRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var order models.LogisticsOrder
	if err := database.DB.First(&order, orderID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "订单不存在"})
		return
	}

	occurredAt := time.Now()
	if req.OccurredAt != "" {
		if t, err := time.Parse(time.RFC3339, req.OccurredAt); err == nil {
			occurredAt = t
		}
	}

	status := req.Status
	if status == "" {
		status = models.NodeStatusInTransit
	}

	node := models.LogisticsNode{
		OrderID:        order.ID,
		NodeName:       req.NodeName,
		Status:         status,
		Location:       req.Location,
		Description:    req.Description,
		Operator:       req.Operator,
		IsAbnormal:     req.IsAbnormal,
		AbnormalReason: req.AbnormalReason,
		OccurredAt:     occurredAt,
	}

	if err := database.DB.Create(&node).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if req.IsAbnormal && order.Status != models.OrderStatusException {
		order.Status = models.OrderStatusException
		database.DB.Save(&order)
	} else if !req.IsAbnormal && status == models.NodeStatusDelivered {
		order.Status = models.OrderStatusDelivered
		database.DB.Save(&order)
	} else if status == models.NodeStatusArrived && order.Status != models.OrderStatusException && order.Status != models.OrderStatusDelivered {
		order.Status = models.OrderStatusArrived
		database.DB.Save(&order)
	} else if status == models.NodeStatusInTransit && order.Status != models.OrderStatusException && order.Status != models.OrderStatusArrived && order.Status != models.OrderStatusDelivered {
		order.Status = models.OrderStatusInTransit
		database.DB.Save(&order)
	} else if status == models.NodeStatusPicked && order.Status == models.OrderStatusPending {
		order.Status = models.OrderStatusPicked
		database.DB.Save(&order)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "节点添加成功",
		"data":    node,
	})
}

func GetNodes(c *gin.Context) {
	orderID := c.Param("id")
	var nodes []models.LogisticsNode

	if err := database.DB.Where("order_id = ?", orderID).Order("occurred_at ASC, id ASC").Find(&nodes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": nodes})
}

func UpdateNode(c *gin.Context) {
	id := c.Param("node_id")
	var node models.LogisticsNode

	if err := database.DB.First(&node, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "节点不存在"})
		return
	}

	var req struct {
		NodeName       string `json:"node_name"`
		Status         string `json:"status"`
		Location       string `json:"location"`
		Description    string `json:"description"`
		Operator       string `json:"operator"`
		IsAbnormal     *bool  `json:"is_abnormal"`
		AbnormalReason string `json:"abnormal_reason"`
		OccurredAt     string `json:"occurred_at"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.NodeName != "" {
		node.NodeName = req.NodeName
	}
	if req.Status != "" {
		node.Status = req.Status
	}
	if req.Location != "" {
		node.Location = req.Location
	}
	if req.Description != "" {
		node.Description = req.Description
	}
	if req.Operator != "" {
		node.Operator = req.Operator
	}
	if req.IsAbnormal != nil {
		node.IsAbnormal = *req.IsAbnormal
	}
	if req.AbnormalReason != "" {
		node.AbnormalReason = req.AbnormalReason
	}
	if req.OccurredAt != "" {
		if t, err := time.Parse(time.RFC3339, req.OccurredAt); err == nil {
			node.OccurredAt = t
		}
	}

	if err := database.DB.Save(&node).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	refreshOrderStatus(node.OrderID)

	c.JSON(http.StatusOK, gin.H{
		"message": "节点更新成功",
		"data":    node,
	})
}

func DeleteNode(c *gin.Context) {
	id := c.Param("node_id")

	var node models.LogisticsNode
	if err := database.DB.First(&node, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "节点不存在"})
		return
	}

	orderID := node.OrderID

	if err := database.DB.Delete(&node).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	refreshOrderStatus(orderID)

	c.JSON(http.StatusOK, gin.H{"message": "节点删除成功"})
}

func GetProgress(c *gin.Context) {
	orderID := c.Param("id")
	var order models.LogisticsOrder

	if err := database.DB.First(&order, orderID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "订单不存在"})
		return
	}

	var nodes []models.LogisticsNode
	database.DB.Where("order_id = ?", orderID).Order("occurred_at ASC, id ASC").Find(&nodes)

	totalStages := 5
	currentStage := 0

	switch order.Status {
	case models.OrderStatusPending:
		currentStage = 0
	case models.OrderStatusPicked:
		currentStage = 1
	case models.OrderStatusInTransit:
		currentStage = 2
	case models.OrderStatusArrived:
		currentStage = 3
	case models.OrderStatusDelivered:
		currentStage = 5
	case models.OrderStatusException:
		currentStage = 2
	}

	progress := float64(currentStage) / float64(totalStages) * 100
	if progress > 100 {
		progress = 100
	}

	c.JSON(http.StatusOK, gin.H{
		"progress":     progress,
		"currentStage": currentStage,
		"totalStages":  totalStages,
		"status":       order.Status,
		"node_count":   len(nodes),
	})
}
