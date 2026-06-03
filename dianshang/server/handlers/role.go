package handlers

import (
	"server/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

var roles = []string{"super", "admin", "operator"}

func GetRoles(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": roles})
}

func GetRoleMenus(c *gin.Context) {
	role := c.Param("role")
	if !isValidRole(role) {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "角色不存在"})
		return
	}
	if role == "super" {
		var menus []models.Menu
		models.DB.Where("visible = 1").Order("sort").Find(&menus)
		tree := BuildMenuTree(menus, 0)
		c.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"tree": tree, "ids": getAllMenuIDs(menus)}})
		return
	}
	var roleMenus []models.RoleMenu
	models.DB.Where("role = ?", role).Find(&roleMenus)
	menuIDs := make([]uint, 0)
	for _, rm := range roleMenus {
		menuIDs = append(menuIDs, rm.MenuID)
	}
	var menus []models.Menu
	models.DB.Where("visible = 1").Order("sort").Find(&menus)
	tree := BuildMenuTreeWithChecked(menus, 0, menuIDs)
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"tree": tree, "ids": menuIDs}})
}

func getAllMenuIDs(menus []models.Menu) []uint {
	ids := make([]uint, 0)
	for _, menu := range menus {
		ids = append(ids, menu.ID)
	}
	return ids
}

func BuildMenuTreeWithChecked(menus []models.Menu, parentID uint, checkedIDs []uint) []models.Menu {
	tree := []models.Menu{}
	for _, menu := range menus {
		if menu.ParentID == parentID {
			menu.Children = BuildMenuTreeWithChecked(menus, menu.ID, checkedIDs)
			tree = append(tree, menu)
		}
	}
	return tree
}

func SetRoleMenus(c *gin.Context) {
	role := c.Param("role")
	if role == "super" {
		c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "超级管理员权限不可修改"})
		return
	}
	if !isValidRole(role) {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "角色不存在"})
		return
	}
	var req struct {
		MenuIDs []uint `json:"menu_ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}
	models.DB.Where("role = ?", role).Delete(&models.RoleMenu{})
	for _, menuID := range req.MenuIDs {
		var menu models.Menu
		if err := models.DB.Where("id = ?", menuID).First(&menu).Error; err == nil {
			roleMenu := models.RoleMenu{Role: role, MenuID: menuID}
			models.DB.Create(&roleMenu)
		}
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "权限设置成功"})
}

func isValidRole(role string) bool {
	for _, r := range roles {
		if r == role {
			return true
		}
	}
	return false
}