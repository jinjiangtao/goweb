package handlers

import (
	"net/http"

	"zhishiku-server/middleware"
	"zhishiku-server/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func GetCategories(c *gin.Context) {
	var categories []models.Category
	if err := db.Order("sort_order ASC, id ASC").Find(&categories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取分类列表失败"})
		return
	}
	c.JSON(http.StatusOK, categories)
}

type CategoryRequest struct {
	Name      string `json:"name" binding:"required"`
	SortOrder int    `json:"sort_order"`
}

func CreateCategory(c *gin.Context) {
	var req CategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	category := models.Category{
		Name:      req.Name,
		SortOrder: req.SortOrder,
	}

	if err := db.Create(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建分类失败"})
		return
	}

	c.JSON(http.StatusCreated, category)
}

func UpdateCategory(c *gin.Context) {
	id := c.Param("id")

	var category models.Category
	if err := db.First(&category, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "分类不存在"})
		return
	}

	var req CategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	category.Name = req.Name
	category.SortOrder = req.SortOrder

	if err := db.Save(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新分类失败"})
		return
	}

	c.JSON(http.StatusOK, category)
}

func DeleteCategory(c *gin.Context) {
	id := c.Param("id")

	var category models.Category
	if err := db.First(&category, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "分类不存在"})
		return
	}

	var articleCount int64
	db.Model(&models.Article{}).Where("category_id = ? AND deleted_at IS NULL", id).Count(&articleCount)
	if articleCount > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "该分类下已有文章，无法删除"})
		return
	}

	if err := db.Delete(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除分类失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

func InitDefaultAdmin() {
	var admin models.User
	if err := db.Where("username = ?", "admin").First(&admin).Error; err != nil {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		db.Create(&models.User{
			Username: "admin",
			Password: string(hashedPassword),
			Nickname: "管理员",
			Role:     "admin",
		})
	}

	var count int64
	db.Model(&models.Category{}).Count(&count)
	if count == 0 {
		db.Create(&models.Category{Name: "技术文档", SortOrder: 1})
		db.Create(&models.Category{Name: "产品规范", SortOrder: 2})
		db.Create(&models.Category{Name: "团队经验", SortOrder: 3})
	}
}

func CheckCategoryPermission() gin.HandlerFunc {
	return func(c *gin.Context) {
		role := middleware.GetRole(c)
		if role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "需要管理员权限"})
			c.Abort()
			return
		}
		c.Next()
	}
}
