package handlers

import (
	"net/http"
	"time"
	"weixiu/config"
	"weixiu/middleware"
	"weixiu/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6"`
	RealName string `json:"real_name"`
	Role     string `json:"role" binding:"required,oneof=admin user technician"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Area     string `json:"area"`
	Skills   string `json:"skills"`
}

func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误", "detail": err.Error()})
		return
	}

	var user models.User
	if err := config.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &middleware.Claims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     string(user.Role),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "weixiu-oms",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(config.JWTSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成令牌失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
		"user": gin.H{
			"id":        user.ID,
			"username":  user.Username,
			"real_name": user.RealName,
			"role":      user.Role,
			"phone":     user.Phone,
			"email":     user.Email,
			"area":      user.Area,
		},
	})
}

func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误", "detail": err.Error()})
		return
	}

	var existingUser models.User
	if config.DB.Where("username = ?", req.Username).First(&existingUser).Error == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户名已存在"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "密码加密失败"})
		return
	}

	user := models.User{
		Username: req.Username,
		Password: string(hashedPassword),
		RealName: req.RealName,
		Role:     models.UserRole(req.Role),
		Phone:    req.Phone,
		Email:    req.Email,
		Area:     req.Area,
		Skills:   req.Skills,
	}

	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "用户创建失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "用户注册成功",
		"user": gin.H{
			"id":        user.ID,
			"username":  user.Username,
			"real_name": user.RealName,
			"role":      user.Role,
		},
	})
}

func GetCurrentUser(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":        user.ID,
		"username":  user.Username,
		"real_name": user.RealName,
		"role":      user.Role,
		"phone":     user.Phone,
		"email":     user.Email,
		"area":      user.Area,
		"skills":    user.Skills,
	})
}

func GetUsers(c *gin.Context) {
	role := c.Query("role")
	var users []models.User
	query := config.DB.Select("id, username, real_name, role, phone, email, area, skills, created_at")
	if role != "" {
		query = query.Where("role = ?", role)
	}
	query.Find(&users)
	c.JSON(http.StatusOK, users)
}

func GetTechnicians(c *gin.Context) {
	area := c.Query("area")
	var technicians []models.User
	query := config.DB.Select("id, username, real_name, role, phone, email, area, skills, created_at").Where("role = ?", models.RoleTechnician)
	if area != "" {
		query = query.Where("area = ?", area)
	}
	query.Find(&technicians)
	c.JSON(http.StatusOK, technicians)
}

func InitDefaultUsers() {
	var count int64
	config.DB.Model(&models.User{}).Count(&count)
	if count > 0 {
		return
	}

	users := []models.User{
		{
			Username: "admin",
			RealName: "系统管理员",
			Role:     models.RoleAdmin,
			Phone:    "13800000001",
			Area:     "总部",
		},
		{
			Username: "user1",
			RealName: "报修用户张三",
			Role:     models.RoleUser,
			Phone:    "13800000002",
			Area:     "教学楼A区",
		},
		{
			Username: "tech1",
			RealName: "运维工程师李四",
			Role:     models.RoleTechnician,
			Phone:    "13800000003",
			Area:     "教学楼A区",
			Skills:   "电脑,打印机,网络设备",
		},
		{
			Username: "tech2",
			RealName: "运维工程师王五",
			Role:     models.RoleTechnician,
			Phone:    "13800000004",
			Area:     "教学楼B区",
			Skills:   "投影仪,空调,其他设备",
		},
	}

	for i := range users {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
		users[i].Password = string(hashedPassword)
		config.DB.Create(&users[i])
	}
}
