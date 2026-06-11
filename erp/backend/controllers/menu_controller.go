package controllers

import (
	"erp/database"
	"erp/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CreateMenuRequest struct {
	ParentID   *uint  `json:"parentId"`
	Name       string `json:"name" binding:"required"`
	Path       string `json:"path"`
	Component  string `json:"component"`
	Icon       string `json:"icon"`
	Sort       int    `json:"sort"`
	Type       string `json:"type"`
	Permission string `json:"permission"`
	Status     int    `json:"status"`
}

type UpdateMenuRequest struct {
	ParentID   *uint  `json:"parentId"`
	Name       string `json:"name"`
	Path       string `json:"path"`
	Component  string `json:"component"`
	Icon       string `json:"icon"`
	Sort       int    `json:"sort"`
	Type       string `json:"type"`
	Permission string `json:"permission"`
	Status     int    `json:"status"`
}

func buildMenuTree(menus []models.Menu, parentID *uint) []models.Menu {
	var tree []models.Menu
	for _, menu := range menus {
		if (menu.ParentID == nil && parentID == nil) ||
			(menu.ParentID != nil && parentID != nil && *menu.ParentID == *parentID) {
			menu.Children = buildMenuTree(menus, &menu.ID)
			tree = append(tree, menu)
		}
	}
	return tree
}

func GetMenus(c *gin.Context) {
	var menus []models.Menu
	database.DB.Order("sort ASC").Find(&menus)

	tree := buildMenuTree(menus, nil)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data":    tree,
	})
}

func GetMenu(c *gin.Context) {
	id := c.Param("id")

	var menu models.Menu
	if err := database.DB.First(&menu, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "菜单不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data":    menu,
	})
}

func CreateMenu(c *gin.Context) {
	var req CreateMenuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	menu := models.Menu{
		ParentID:   req.ParentID,
		Name:       req.Name,
		Path:       req.Path,
		Component:  req.Component,
		Icon:       req.Icon,
		Sort:       req.Sort,
		Type:       req.Type,
		Permission: req.Permission,
		Status:     req.Status,
	}

	if menu.Status == 0 {
		menu.Status = 1
	}
	if menu.Type == "" {
		menu.Type = "menu"
	}

	if err := database.DB.Create(&menu).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "创建成功",
		"data":    menu,
	})
}

func UpdateMenu(c *gin.Context) {
	id := c.Param("id")

	var req UpdateMenuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	var menu models.Menu
	if err := database.DB.First(&menu, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "菜单不存在"})
		return
	}

	if req.ParentID != nil {
		menu.ParentID = req.ParentID
	}
	if req.Name != "" {
		menu.Name = req.Name
	}
	if req.Path != "" {
		menu.Path = req.Path
	}
	if req.Component != "" {
		menu.Component = req.Component
	}
	if req.Icon != "" {
		menu.Icon = req.Icon
	}
	if req.Sort != 0 {
		menu.Sort = req.Sort
	}
	if req.Type != "" {
		menu.Type = req.Type
	}
	if req.Permission != "" {
		menu.Permission = req.Permission
	}
	if req.Status != 0 {
		menu.Status = req.Status
	}

	if err := database.DB.Save(&menu).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "更新成功",
		"data":    menu,
	})
}

func DeleteMenu(c *gin.Context) {
	id := c.Param("id")

	var childCount int64
	database.DB.Model(&models.Menu{}).Where("parent_id = ?", id).Count(&childCount)
	if childCount > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "该菜单下还有子菜单，无法删除"})
		return
	}

	if err := database.DB.Delete(&models.Menu{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "删除成功",
	})
}
