package handlers

import (
	"net/http"

	"student-signup-server/models"
	"student-signup-server/utils"

	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	admin, err := models.GetAdminByUsername(req.Username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	if admin.Status == models.StatusDisabled {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户已被禁用"})
		return
	}

	if !admin.CheckPassword(req.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	go UpdateLastLoginTime(admin.ID)

	token, err := utils.GenerateToken(admin.ID, admin.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成token失败"})
		return
	}

	menus, err := models.GetMenuByRole(admin.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取菜单失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":   token,
		"message": "登录成功",
		"user": gin.H{
			"id":       admin.ID,
			"username": admin.Username,
			"nickname": admin.Nickname,
			"role":     admin.Role,
		},
		"menus": menus,
	})
}
