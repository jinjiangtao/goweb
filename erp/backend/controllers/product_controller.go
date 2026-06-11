package controllers

import (
	"erp/database"
	"erp/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CreateProductRequest struct {
	Name  string  `json:"name" binding:"required"`
	Code  string  `json:"code" binding:"required"`
	Price float64 `json:"price"`
	Spec  string  `json:"spec"`
}

type UpdateProductRequest struct {
	Name  string  `json:"name"`
	Code  string  `json:"code"`
	Price float64 `json:"price"`
	Spec  string  `json:"spec"`
}

func GetProducts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	name := c.Query("name")

	var products []models.Product
	var total int64

	query := database.DB.Model(&models.Product{})

	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}

	query.Count(&total)
	query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&products)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data": gin.H{
			"list":     products,
			"total":    total,
			"page":     page,
			"pageSize": pageSize,
		},
	})
}

func GetProduct(c *gin.Context) {
	id := c.Param("id")

	var product models.Product
	if err := database.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "产品不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data":    product,
	})
}

func CreateProduct(c *gin.Context) {
	var req CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	var existingProduct models.Product
	if err := database.DB.Where("code = ?", req.Code).First(&existingProduct).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"code": 409, "message": "产品编码已存在"})
		return
	}

	product := models.Product{
		Name:  req.Name,
		Code:  req.Code,
		Price: req.Price,
		Spec:  req.Spec,
	}

	if err := database.DB.Create(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "创建成功",
		"data":    product,
	})
}

func UpdateProduct(c *gin.Context) {
	id := c.Param("id")

	var req UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	var product models.Product
	if err := database.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "产品不存在"})
		return
	}

	if req.Name != "" {
		product.Name = req.Name
	}
	if req.Code != "" {
		product.Code = req.Code
	}
	if req.Price != 0 {
		product.Price = req.Price
	}
	if req.Spec != "" {
		product.Spec = req.Spec
	}

	if err := database.DB.Save(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "更新成功",
		"data":    product,
	})
}

func DeleteProduct(c *gin.Context) {
	id := c.Param("id")

	if err := database.DB.Delete(&models.Product{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "删除成功",
	})
}
