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

type PurchaseOrderItemRequest struct {
	ProductID uint    `json:"productId" binding:"required"`
	Quantity  float64 `json:"quantity" binding:"required"`
	UnitPrice float64 `json:"unitPrice" binding:"required"`
}

type CreatePurchaseOrderRequest struct {
	SupplierID   uint                       `json:"supplierId" binding:"required"`
	Items        []PurchaseOrderItemRequest `json:"items" binding:"required,min=1,dive"`
	ExpectedDate *string                    `json:"expectedDate"`
	Remark       string                     `json:"remark"`
}

type UpdatePurchaseOrderRequest struct {
	SupplierID   uint                       `json:"supplierId"`
	Items        []PurchaseOrderItemRequest `json:"items"`
	ExpectedDate *string                    `json:"expectedDate"`
	Remark       string                     `json:"remark"`
}

type UpdatePurchaseOrderStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

func generatePurchaseOrderNo() string {
	now := time.Now()
	dateStr := now.Format("20060102")
	
	var count int64
	database.DB.Model(&models.PurchaseOrder{}).
		Where("created_at >= ? AND created_at < ?", 
			time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()),
			time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())).
		Count(&count)
	
	seq := count + 1
	return fmt.Sprintf("PO%s%04d", dateStr, seq)
}

func GetPurchaseOrders(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	supplierID := c.Query("supplierId")
	status := c.Query("status")
	orderNo := c.Query("orderNo")

	var orders []models.PurchaseOrder
	var total int64

	query := database.DB.Model(&models.PurchaseOrder{}).Preload("Supplier").Preload("CreatedBy")

	if supplierID != "" {
		query = query.Where("supplier_id = ?", supplierID)
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

func GetPurchaseOrder(c *gin.Context) {
	id := c.Param("id")

	var order models.PurchaseOrder
	if err := database.DB.Preload("Supplier").
		Preload("CreatedBy").
		Preload("Items").
		Preload("Items.Product").
		First(&order, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "采购订单不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data":    order,
	})
}

func CreatePurchaseOrder(c *gin.Context) {
	var req CreatePurchaseOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	userID, _ := c.Get("userId")

	var supplier models.Supplier
	if err := database.DB.First(&supplier, req.SupplierID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "供应商不存在"})
		return
	}

	orderNo := generatePurchaseOrderNo()

	totalAmount := 0.0
	var items []models.PurchaseOrderItem
	for _, itemReq := range req.Items {
		var product models.Product
		if err := database.DB.First(&product, itemReq.ProductID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "产品不存在"})
			return
		}
		
		amount := itemReq.Quantity * itemReq.UnitPrice
		totalAmount += amount
		
		items = append(items, models.PurchaseOrderItem{
			ProductID: itemReq.ProductID,
			Quantity:  itemReq.Quantity,
			UnitPrice: itemReq.UnitPrice,
			Amount:    amount,
		})
	}

	var expectedDate *time.Time
	if req.ExpectedDate != nil {
		if t, err := time.Parse("2006-01-02", *req.ExpectedDate); err == nil {
			expectedDate = &t
		}
	}

	order := models.PurchaseOrder{
		OrderNo:      orderNo,
		SupplierID:   req.SupplierID,
		Items:        items,
		TotalAmount:  totalAmount,
		Status:       models.PurchaseOrderStatusDraft,
		ExpectedDate: expectedDate,
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

func UpdatePurchaseOrder(c *gin.Context) {
	id := c.Param("id")

	var req UpdatePurchaseOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	var order models.PurchaseOrder
	if err := database.DB.Preload("Items").First(&order, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "采购订单不存在"})
		return
	}

	if order.Status != models.PurchaseOrderStatusDraft {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "只有草稿状态可以编辑"})
		return
	}

	if req.SupplierID != 0 {
		order.SupplierID = req.SupplierID
	}

	if req.Items != nil && len(req.Items) > 0 {
		totalAmount := 0.0
		var items []models.PurchaseOrderItem
		
		for _, itemReq := range req.Items {
			var product models.Product
			if err := database.DB.First(&product, itemReq.ProductID).Error; err != nil {
				c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "产品不存在"})
				return
			}
			
			amount := itemReq.Quantity * itemReq.UnitPrice
			totalAmount += amount
			
			items = append(items, models.PurchaseOrderItem{
				ProductID: itemReq.ProductID,
				Quantity:  itemReq.Quantity,
				UnitPrice: itemReq.UnitPrice,
				Amount:    amount,
			})
		}
		
		database.DB.Delete(&models.PurchaseOrderItem{}, "purchase_order_id = ?", order.ID)
		order.Items = items
		order.TotalAmount = totalAmount
	}

	if req.ExpectedDate != nil {
		if t, err := time.Parse("2006-01-02", *req.ExpectedDate); err == nil {
			order.ExpectedDate = &t
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

func UpdatePurchaseOrderStatus(c *gin.Context) {
	id := c.Param("id")

	var req UpdatePurchaseOrderStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	var order models.PurchaseOrder
	if err := database.DB.First(&order, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "采购订单不存在"})
		return
	}

	newStatus := models.PurchaseOrderStatus(req.Status)
	validTransitions := map[models.PurchaseOrderStatus][]models.PurchaseOrderStatus{
		models.PurchaseOrderStatusDraft: {
			models.PurchaseOrderStatusApproved,
			models.PurchaseOrderStatusCancelled,
		},
		models.PurchaseOrderStatusApproved: {
			models.PurchaseOrderStatusPartialReceive,
			models.PurchaseOrderStatusCompleted,
			models.PurchaseOrderStatusCancelled,
		},
		models.PurchaseOrderStatusPartialReceive: {
			models.PurchaseOrderStatusPartialReceive,
			models.PurchaseOrderStatusCompleted,
			models.PurchaseOrderStatusCancelled,
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

func DeletePurchaseOrder(c *gin.Context) {
	id := c.Param("id")

	var order models.PurchaseOrder
	if err := database.DB.First(&order, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "采购订单不存在"})
		return
	}

	if order.Status != models.PurchaseOrderStatusDraft {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "只有草稿状态可以删除"})
		return
	}

	database.DB.Delete(&models.PurchaseOrderItem{}, "purchase_order_id = ?", order.ID)
	if err := database.DB.Delete(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "删除成功",
	})
}
