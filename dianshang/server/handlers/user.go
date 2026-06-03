package handlers

import (
	"server/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetUsers(c *gin.Context) {
	username := c.Query("username")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	offset := (page - 1) * size

	var users []models.AdminUser
	query := models.DB.Model(&models.AdminUser{})
	if username != "" {
		query = query.Where("username LIKE ?", "%"+username+"%")
	}
	var total int
	query.Count(&total)
	query.Order("created_at DESC").Offset(offset).Limit(size).Find(&users)

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"list": users, "total": total}})
}

func CreateUser(c *gin.Context) {
	var user models.AdminUser
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}
	if user.Username == "" || user.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "用户名和密码不能为空"})
		return
	}
	hashedPassword, err := models.GeneratePasswordHash(user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "密码加密失败"})
		return
	}
	user.Password = hashedPassword
	if err := models.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "创建成功", "data": user})
}

func UpdateUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var user models.AdminUser
	if err := models.DB.Where("id = ?", id).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "用户不存在"})
		return
	}
	if user.Role == "super" {
		c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "超级管理员不可修改"})
		return
	}
	var req struct {
		Nickname string `json:"nickname"`
		Role     string `json:"role"`
		Status   int    `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}
	if req.Nickname != "" {
		user.Nickname = req.Nickname
	}
	if req.Role != "" {
		user.Role = req.Role
	}
	user.Status = req.Status
	models.DB.Save(&user)
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "更新成功", "data": user})
}

func DeleteUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var user models.AdminUser
	if err := models.DB.Where("id = ?", id).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "用户不存在"})
		return
	}
	if user.Role == "super" {
		c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "超级管理员不可删除"})
		return
	}
	models.DB.Delete(&user)
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "删除成功"})
}

func UpdateUserStatus(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var user models.AdminUser
	if err := models.DB.Where("id = ?", id).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "用户不存在"})
		return
	}
	if user.Role == "super" {
		c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "超级管理员状态不可修改"})
		return
	}
	user.Status = 1 - user.Status
	models.DB.Save(&user)
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "状态更新成功", "data": user})
}

func ResetPassword(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var user models.AdminUser
	if err := models.DB.Where("id = ?", id).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "用户不存在"})
		return
	}
	hashedPassword, err := models.GeneratePasswordHash("123456")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "密码加密失败"})
		return
	}
	user.Password = hashedPassword
	models.DB.Save(&user)
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "密码重置成功"})
}