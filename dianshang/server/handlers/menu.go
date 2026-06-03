package handlers

import (
	"server/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetMenus(c *gin.Context) {
	var menus []models.Menu
	models.DB.Where("visible = 1").Order("sort").Find(&menus)
	tree := BuildMenuTree(menus, 0)
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": tree})
}

func BuildMenuTree(menus []models.Menu, parentID uint) []models.Menu {
	tree := []models.Menu{}
	for _, menu := range menus {
		if menu.ParentID == parentID {
			menu.Children = BuildMenuTree(menus, menu.ID)
			tree = append(tree, menu)
		}
	}
	return tree
}

func GetMenuTreeByRole(role string) []models.Menu {
	if role == "super" {
		var menus []models.Menu
		models.DB.Where("visible = 1").Order("sort").Find(&menus)
		return BuildMenuTree(menus, 0)
	}
	var roleMenus []models.RoleMenu
	models.DB.Where("role = ?", role).Find(&roleMenus)
	menuIDs := make([]uint, 0)
	for _, rm := range roleMenus {
		menuIDs = append(menuIDs, rm.MenuID)
	}
	var menus []models.Menu
	models.DB.Where("id IN (?) AND visible = 1", menuIDs).Order("sort").Find(&menus)
	return BuildMenuTree(menus, 0)
}

func CreateMenu(c *gin.Context) {
	var menu models.Menu
	if err := c.ShouldBindJSON(&menu); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}
	if menu.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "菜单名称不能为空"})
		return
	}
	if menu.ParentID != 0 {
		var parent models.Menu
		if err := models.DB.Where("id = ?", menu.ParentID).First(&parent).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "父菜单不存在"})
			return
		}
		if parent.ParentID != 0 {
			var grandParent models.Menu
			if err := models.DB.Where("id = ?", parent.ParentID).First(&grandParent).Error; err == nil {
				c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "菜单最多支持三级"})
				return
			}
		}
	}
	models.DB.Create(&menu)
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "创建成功", "data": menu})
}

func UpdateMenu(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var menu models.Menu
	if err := models.DB.Where("id = ?", id).First(&menu).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "菜单不存在"})
		return
	}
	var req struct {
		Name    string `json:"name"`
		Path    string `json:"path"`
		Icon    string `json:"icon"`
		Sort    int    `json:"sort"`
		Visible int    `json:"visible"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}
	if req.Name != "" {
		menu.Name = req.Name
	}
	if req.Path != "" {
		menu.Path = req.Path
	}
	if req.Icon != "" {
		menu.Icon = req.Icon
	}
	menu.Sort = req.Sort
	menu.Visible = req.Visible
	models.DB.Save(&menu)
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "更新成功", "data": menu})
}

func DeleteMenu(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var menu models.Menu
	if err := models.DB.Where("id = ?", id).First(&menu).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "菜单不存在"})
		return
	}
	DeleteMenuWithChildren(uint(id))
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "删除成功"})
}

func DeleteMenuWithChildren(id uint) {
	var children []models.Menu
	models.DB.Where("parent_id = ?", id).Find(&children)
	for _, child := range children {
		DeleteMenuWithChildren(child.ID)
	}
	models.DB.Delete(&models.Menu{}, id)
	models.DB.Where("menu_id = ?", id).Delete(&models.RoleMenu{})
}