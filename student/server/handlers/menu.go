package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"student-signup-server/models"
)

type CreateMenuRequest struct {
	Name     string `json:"name" binding:"required"`
	Path     string `json:"path"`
	Icon     string `json:"icon"`
	Sort     int    `json:"sort"`
	ParentID uint   `json:"parent_id"`
	Visible  int    `json:"visible"`
}

type UpdateMenuRequest struct {
	Name     string `json:"name"`
	Path     string `json:"path"`
	Icon     string `json:"icon"`
	Sort     int    `json:"sort"`
	ParentID uint   `json:"parent_id"`
	Visible  int    `json:"visible"`
}

func GetMenus(c *gin.Context) {
	menus, err := models.GetAllMenus()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取菜单列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": menus,
	})
}

func GetMenuTree(c *gin.Context) {
	adminID := c.GetUint("admin_id")
	admin, err := models.GetAdminByID(adminID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取菜单失败"})
		return
	}

	menus, err := models.GetMenuByRole(admin.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取菜单失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": menus,
	})
}

func GetParentMenus(c *gin.Context) {
	menus, err := models.GetParentMenus()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取父级菜单失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": menus,
	})
}

func CreateMenu(c *gin.Context) {
	var req CreateMenuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	if req.Visible != 0 && req.Visible != 1 {
		req.Visible = 1
	}

	menu := models.Menu{
		Name:     req.Name,
		Path:     req.Path,
		Icon:     req.Icon,
		Sort:     req.Sort,
		ParentID: req.ParentID,
		Visible:  req.Visible,
	}

	err := models.CreateMenu(&menu)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建菜单失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "创建成功"})
}

func UpdateMenu(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req UpdateMenuRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	menu, err := models.GetMenuByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "菜单不存在"})
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
	if req.Sort != 0 {
		menu.Sort = req.Sort
	}
	if req.ParentID != 0 {
		menu.ParentID = req.ParentID
	}
	if req.Visible == 0 || req.Visible == 1 {
		menu.Visible = req.Visible
	}

	err = models.UpdateMenu(menu)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新菜单失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}

func DeleteMenu(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	err := models.DeleteMenu(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除菜单失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}