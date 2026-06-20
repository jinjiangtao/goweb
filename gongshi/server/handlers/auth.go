package handlers

import (
	"gongshi/database"
	"gongshi/models"
	"net/http"
	"strconv"

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

	var user models.User
	if err := database.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	if user.Password != req.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "登录成功",
		"user":    user,
		"token":   "mock-token-" + strconv.Itoa(int(user.ID)),
	})
}

func GetUsers(c *gin.Context) {
	var users []models.User
	database.DB.Where("role = ?", "employee").Find(&users)
	c.JSON(http.StatusOK, users)
}

func GetCurrentUser(c *gin.Context) {
	userID, _ := strconv.Atoi(c.GetHeader("X-User-ID"))
	if userID == 0 {
		userID = 2
	}

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}
	c.JSON(http.StatusOK, user)
}
