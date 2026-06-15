package handlers

import (
	"net/http"

	"fangwu-server/models"
	"fangwu-server/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, http.StatusBadRequest, "请输入用户名和密码")
		return
	}

	var admin models.Admin
	if err := models.DB.Where("username = ?", req.Username).First(&admin).Error; err != nil {
		utils.Fail(c, http.StatusUnauthorized, "用户名或密码错误")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(req.Password)); err != nil {
		utils.Fail(c, http.StatusUnauthorized, "用户名或密码错误")
		return
	}

	token, err := utils.GenerateToken(admin.ID, admin.Username)
	if err != nil {
		utils.Fail(c, http.StatusInternalServerError, "生成令牌失败")
		return
	}

	utils.Success(c, gin.H{
		"token":    token,
		"username": admin.Username,
	})
}

func GetProfile(c *gin.Context) {
	adminID := c.GetUint("admin_id")
	var admin models.Admin
	if err := models.DB.First(&admin, adminID).Error; err != nil {
		utils.Fail(c, http.StatusNotFound, "用户不存在")
		return
	}
	utils.Success(c, gin.H{
		"id":       admin.ID,
		"username": admin.Username,
	})
}

func ChangePassword(c *gin.Context) {
	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, http.StatusBadRequest, "请填写完整信息，新密码至少6位")
		return
	}

	adminID := c.GetUint("admin_id")
	var admin models.Admin
	if err := models.DB.First(&admin, adminID).Error; err != nil {
		utils.Fail(c, http.StatusNotFound, "用户不存在")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(req.OldPassword)); err != nil {
		utils.Fail(c, http.StatusBadRequest, "原密码错误")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		utils.Fail(c, http.StatusInternalServerError, "密码加密失败")
		return
	}

	models.DB.Model(&admin).Update("password", string(hashedPassword))
	utils.Success(c, nil)
}
