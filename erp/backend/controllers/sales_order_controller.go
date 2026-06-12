package controllers

import (
	"erp/database"
	"erp/models"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type SalesOrderItemRequest struct {
	ProductID uint    `json:"productId" binding:"required"`
	Quantity  float64 `json:"quantity" binding:"required"`
	UnitPrice float64 `json:"unitPrice" binding:"required"`
}

type CreateSalesOrderRequest struct {
	CustomerID   uint                     `json:"customerId" binding:"required"`
	Items        []SalesOrderItemRequest  `json:"items" binding:"required,min=1,dive"`
	DeliveryDate *string                  `json:"deliveryDate"`
	Remark       string                   `json:"remark"`
}

type UpdateSalesOrderRequest struct {
	CustomerID   uint                     `json:"customerId"`
	Items        []SalesOrderItemRequest  `json:"items"`
	DeliveryDate *string                  `json:"deliveryDate"`
	Remark       string                   `json:"remark"`
}

type UpdateSalesOrderStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

func generateSalesOrderNo() string {
	now := time.Now()
	dateStr := now.Format("20060102")
	
	var count int64
	database.DB.Model(&models.SalesOrder{}).
		Where("created_at >= ? AND created_at < ?", 
			time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()),
			time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())).
		Count(&count)
	
	seq := count + 1
	return fmt.Sprintf("SO%s%04d", dateStr, seq)
}

func GetSalesOrders(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	customerID := c.Query("customerId")
	status := c.Query("status")
	orderNo := c.Query("orderNo")

	var orders []models.SalesOrder
	var total int64

	query := database.DB.Model(&models.SalesOrder{}).Preload("Customer").Preload("CreatedBy")

	if customerID != "" {
		query = query.Where("customer_id = ?", customerID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if orderNo != "" {
		query = query.Where("order_no LIKE ?", "%"+orderNo+"%")
	}

	query.Count(&total)
	query.Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&orders)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data": gin.H{
			"list":     orders,
			"total":    total,
			"page":     page,
			"pageSize": pageSize,
		},
	})
}

func GetSalesOrder(c *gin.Context) {
	id := c.Param("id")

	var order models.SalesOrder
	if err := database.DB.Preload("Customer").
		Preload("CreatedBy").
		Preload("Items").
		Preload("Items.Product").
		First(&order, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "销售订单不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data":    order,
	})
}

func CreateSalesOrder(c *gin.Context) {
	var req CreateSalesOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	userID, _ := c.Get("userId")

	var customer models.Customer
	if err := database.DB.First(&customer, req.CustomerID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "客户不存在"})
		return
	}

	orderNo := generateSalesOrderNo()

	totalAmount := 0.0
	var items []models.SalesOrderItem
	for _, itemReq := range req.Items {
		var product models.Product
		if err := database.DB.First(&product, itemReq.ProductID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "产品不存在"})
			return
		}
		
		amount := itemReq.Quantity * itemReq.UnitPrice
		totalAmount += amount
		
		items = append(items, models.SalesOrderItem{
			ProductID: itemReq.ProductID,
			Quantity:  itemReq.Quantity,
			UnitPrice: itemReq.UnitPrice,
			Amount:    amount,
		})
	}

	var deliveryDate *time.Time
	if req.DeliveryDate != nil {
		if t, err := time.Parse("2006-01-02", *req.DeliveryDate); err == nil {
			deliveryDate = &t
		}
	}

	order := models.SalesOrder{
		OrderNo:      orderNo,
		CustomerID:   req.CustomerID,
		Items:        items,
		TotalAmount:  totalAmount,
		Status:       models.SalesOrderStatusDraft,
		DeliveryDate: deliveryDate,
		Remark:       req.Remark,
		CreatedByID:  userID.(uint),
	}

	if err := database.DB.Create(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "创建成功",
		"data":    order,
	})
}

func UpdateSalesOrder(c *gin.Context) {
	id := c.Param("id")

	var req UpdateSalesOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	var order models.SalesOrder
	if err := database.DB.Preload("Items").First(&order, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "销售订单不存在"})
		return
	}

	if order.Status != models.SalesOrderStatusDraft {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "只有草稿状态可以编辑"})
		return
	}

	if req.CustomerID != 0 {
		order.CustomerID = req.CustomerID
	}

	if req.Items != nil && len(req.Items) > 0 {
		totalAmount := 0.0
		var items []models.SalesOrderItem
		
		for _, itemReq := range req.Items {
			var product models.Product
			if err := database.DB.First(&product, itemReq.ProductID).Error; err != nil {
				c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "产品不存在"})
				return
			}
			
			amount := itemReq.Quantity * itemReq.UnitPrice
			totalAmount += amount
			
			items = append(items, models.SalesOrderItem{
				ProductID: itemReq.ProductID,
				Quantity:  itemReq.Quantity,
				UnitPrice: itemReq.UnitPrice,
				Amount:    amount,
			})
		}
		
		database.DB.Delete(&models.SalesOrderItem{}, "sales_order_id = ?", order.ID)
		order.Items = items
		order.TotalAmount = totalAmount
	}

	if req.DeliveryDate != nil {
		if t, err := time.Parse("2006-01-02", *req.DeliveryDate); err == nil {
			order.DeliveryDate = &t
		}
	}

	if req.Remark != "" {
		order.Remark = req.Remark
	}

	if err := database.DB.Save(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "更新成功",
		"data":    order,
	})
}

func UpdateSalesOrderStatus(c *gin.Context) {
	id := c.Param("id")

	var req UpdateSalesOrderStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	var order models.SalesOrder
	if err := database.DB.First(&order, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "销售订单不存在"})
		return
	}

	newStatus := models.SalesOrderStatus(req.Status)
	validTransitions := map[models.SalesOrderStatus][]models.SalesOrderStatus{
		models.SalesOrderStatusDraft: {
			models.SalesOrderStatusApproved,
			models.SalesOrderStatusCancelled,
		},
		models.SalesOrderStatusApproved: {
			models.SalesOrderStatusPartialShip,
			models.SalesOrderStatusCompleted,
			models.SalesOrderStatusCancelled,
		},
		models.SalesOrderStatusPartialShip: {
			models.SalesOrderStatusPartialShip,
			models.SalesOrderStatusCompleted,
			models.SalesOrderStatusCancelled,
		},
	}

	valid := false
	if order.Status == newStatus {
		valid = true
	} else if allowedStatuses, ok := validTransitions[order.Status]; ok {
		for _, s := range allowedStatuses {
			if s == newStatus {
				valid = true
				break
			}
		}
	}

	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "状态转换无效"})
		return
	}

	order.Status = newStatus
	if err := database.DB.Save(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "更新成功",
		"data":    order,
	})
}

func DeleteSalesOrder(c *gin.Context) {
	id := c.Param("id")

	var order models.SalesOrder
	if err := database.DB.First(&order, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "销售订单不存在"})
		return
	}

	if order.Status != models.SalesOrderStatusDraft {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "只有草稿状态可以删除"})
		return
	}

	database.DB.Delete(&models.SalesOrderItem{}, "sales_order_id = ?", order.ID)
	if err := database.DB.Delete(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "删除成功",
	})
}
