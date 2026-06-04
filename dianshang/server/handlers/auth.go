package handlers

import (
	"net/http"
	"server/models"
	"server/utils"
	"time"

	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}
	var user models.AdminUser
	models.DB.Where("username = ?", req.Username).First(&user)
	if user.Username == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "用户名或密码错误"})
		return
	}
	if !models.CheckPasswordHash(req.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "用户名或密码错误"})
		return
	}
	if user.Status == 0 {
		c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "账号已禁用"})
		return
	}
	token, err := utils.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "生成token失败"})
		return
	}
	models.DB.Model(&user).Update("last_login_at", time.Now())
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "登录成功", "data": gin.H{"token": token, "user": user}})
}

func GetInfo(c *gin.Context) {
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "未登录"})
		return
	}
	claim := claims.(*utils.Claims)
	var user models.AdminUser
	models.DB.Where("id = ?", claim.ID).First(&user)
	menus := GetMenuTreeByRole(user.Role)
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"user": user, "menus": menus}})
}

func Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "退出成功"})
}
