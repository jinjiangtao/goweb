package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"student-signup-server/models"
)

type CreateUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Nickname string `json:"nickname"`
	Role     string `json:"role"`
}

type UpdateUserRequest struct {
	Nickname string `json:"nickname"`
	Role     string `json:"role"`
	Status   int    `json:"status"`
}

type ResetPasswordRequest struct {
	Password string `json:"password" binding:"required"`
}

func GetUsers(c *gin.Context) {
	username := c.Query("username")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	admins, total, err := models.GetAdmins(username, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取用户列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  admins,
		"total": total,
	})
}

func CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	_, err := models.GetAdminByUsername(req.Username)
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户名已存在"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "密码加密失败"})
		return
	}

	if req.Role != models.RoleSuperAdmin && req.Role != models.RoleAdmin {
		req.Role = models.RoleAdmin
	}

	admin := models.Admin{
		Username: req.Username,
		Password: string(hashedPassword),
		Nickname: req.Nickname,
		Role:     req.Role,
		Status:   models.StatusEnabled,
	}

	err = models.CreateAdmin(&admin)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建用户失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "创建成功"})
}

func UpdateUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req UpdateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	admin, err := models.GetAdminByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	if req.Nickname != "" {
		admin.Nickname = req.Nickname
	}
	if req.Role != "" && (req.Role == models.RoleSuperAdmin || req.Role == models.RoleAdmin) {
		admin.Role = req.Role
	}
	if req.Status == models.StatusEnabled || req.Status == models.StatusDisabled {
		admin.Status = req.Status
	}

	err = models.UpdateAdmin(admin)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新用户失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}

func DeleteUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	currentAdminID := c.GetUint("admin_id")

	if uint(id) == currentAdminID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "不能删除自己"})
		return
	}

	err := models.DeleteAdmin(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除用户失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

func ResetPassword(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req ResetPasswordRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	admin, err := models.GetAdminByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "密码加密失败"})
		return
	}

	admin.Password = string(hashedPassword)
	err = models.UpdateAdmin(admin)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "重置密码失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "重置密码成功"})
}

func GetUserInfo(c *gin.Context) {
	adminID := c.GetUint("admin_id")
	admin, err := models.GetAdminByID(adminID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取用户信息失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":       admin.ID,
		"username": admin.Username,
		"nickname": admin.Nickname,
		"role":     admin.Role,
	})
}

func UpdateLastLoginTime(adminID uint) {
	admin, err := models.GetAdminByID(adminID)
	if err != nil {
		return
	}
	admin.LastLoginTime = time.Now()
	models.UpdateAdmin(admin)
}