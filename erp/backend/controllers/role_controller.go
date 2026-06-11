package controllers

import (
	"erp/database"
	"erp/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CreateRoleRequest struct {
	Name        string `json:"name" binding:"required"`
	Code        string `json:"code" binding:"required"`
	Description string `json:"description"`
	Status      int    `json:"status"`
	MenuIDs     []uint `json:"menuIds"`
}

type UpdateRoleRequest struct {
	Name        string `json:"name"`
	Code        string `json:"code"`
	Description string `json:"description"`
	Status      int    `json:"status"`
	MenuIDs     []uint `json:"menuIds"`
}

func GetRoles(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	keyword := c.Query("keyword")

	var roles []models.Role
	var total int64

	query := database.DB.Model(&models.Role{})

	if keyword != "" {
		query = query.Where("name LIKE ? OR code LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	query.Count(&total)
	query.Preload("Menus").Offset((page - 1) * pageSize).Limit(pageSize).Find(&roles)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data": gin.H{
			"list":     roles,
			"total":    total,
			"page":     page,
			"pageSize": pageSize,
		},
	})
}

func GetRole(c *gin.Context) {
	id := c.Param("id")

	var role models.Role
	if err := database.DB.Preload("Menus").First(&role, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "角色不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data":    role,
	})
}

func CreateRole(c *gin.Context) {
	var req CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	var existingRole models.Role
	if err := database.DB.Where("name = ? OR code = ?", req.Name, req.Code).First(&existingRole).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"code": 409, "message": "角色名称或代码已存在"})
		return
	}

	role := models.Role{
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		Status:      req.Status,
	}

	if role.Status == 0 {
		role.Status = 1
	}

	if err := database.DB.Create(&role).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败"})
		return
	}

	if len(req.MenuIDs) > 0 {
		var menus []models.Menu
		database.DB.Find(&menus, req.MenuIDs)
		database.DB.Model(&role).Association("Menus").Replace(menus)
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "创建成功",
		"data":    role,
	})
}

func UpdateRole(c *gin.Context) {
	id := c.Param("id")

	var req UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	var role models.Role
	if err := database.DB.First(&role, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "角色不存在"})
		return
	}

	if req.Name != "" {
		role.Name = req.Name
	}
	if req.Code != "" {
		role.Code = req.Code
	}
	if req.Description != "" {
		role.Description = req.Description
	}
	if req.Status != 0 {
		role.Status = req.Status
	}

	if err := database.DB.Save(&role).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败"})
		return
	}

	if req.MenuIDs != nil {
		var menus []models.Menu
		database.DB.Find(&menus, req.MenuIDs)
		database.DB.Model(&role).Association("Menus").Replace(menus)
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "更新成功",
		"data":    role,
	})
}

func DeleteRole(c *gin.Context) {
	id := c.Param("id")

	var userCount int64
	database.DB.Model(&models.User{}).Where("role_id = ?", id).Count(&userCount)
	if userCount > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "该角色下还有用户，无法删除"})
		return
	}

	if err := database.DB.Delete(&models.Role{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "删除成功",
	})
}

func AssignMenus(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		MenuIDs []uint `json:"menuIds"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	var role models.Role
	if err := database.DB.First(&role, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "角色不存在"})
		return
	}

	var menus []models.Menu
	database.DB.Find(&menus, req.MenuIDs)
	database.DB.Model(&role).Association("Menus").Replace(menus)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "分配成功",
	})
}
