package controllers

import (
	"erp/database"
	"erp/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CreateSupplierRequest struct {
	Name    string `json:"name" binding:"required"`
	Contact string `json:"contact"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
	Remark  string `json:"remark"`
}

type UpdateSupplierRequest struct {
	Name    string `json:"name"`
	Contact string `json:"contact"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
	Remark  string `json:"remark"`
}

func GetSuppliers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	name := c.Query("name")

	var suppliers []models.Supplier
	var total int64

	query := database.DB.Model(&models.Supplier{})

	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}

	query.Count(&total)
	query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&suppliers)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data": gin.H{
			"list":     suppliers,
			"total":    total,
			"page":     page,
			"pageSize": pageSize,
		},
	})
}

func GetSupplier(c *gin.Context) {
	id := c.Param("id")

	var supplier models.Supplier
	if err := database.DB.First(&supplier, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "供应商不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data":    supplier,
	})
}

func CreateSupplier(c *gin.Context) {
	var req CreateSupplierRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	supplier := models.Supplier{
		Name:    req.Name,
		Contact: req.Contact,
		Phone:   req.Phone,
		Address: req.Address,
		Remark:  req.Remark,
	}

	if err := database.DB.Create(&supplier).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "创建成功",
		"data":    supplier,
	})
}

func UpdateSupplier(c *gin.Context) {
	id := c.Param("id")

	var req UpdateSupplierRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	var supplier models.Supplier
	if err := database.DB.First(&supplier, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "供应商不存在"})
		return
	}

	if req.Name != "" {
		supplier.Name = req.Name
	}
	if req.Contact != "" {
		supplier.Contact = req.Contact
	}
	if req.Phone != "" {
		supplier.Phone = req.Phone
	}
	if req.Address != "" {
		supplier.Address = req.Address
	}
	if req.Remark != "" {
		supplier.Remark = req.Remark
	}

	if err := database.DB.Save(&supplier).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "更新成功",
		"data":    supplier,
	})
}

func DeleteSupplier(c *gin.Context) {
	id := c.Param("id")

	if err := database.DB.Delete(&models.Supplier{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "删除成功",
	})
}

