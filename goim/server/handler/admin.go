package handler

import (
	"goim/server/cache"
	"goim/server/service"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const (
	adminUsername = "admin"
	adminPassword = "123456"
	jwtSecret     = "goim-admin-secret-key"
)

type AdminLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AdminLoginResponse struct {
	Token string `json:"token"`
}

func AdminLogin(c *gin.Context) {
	var req AdminLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Username != adminUsername || req.Password != adminPassword {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": adminUsername,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, AdminLoginResponse{Token: tokenString})
}

func AdminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
			c.Abort()
			return
		}

		// Remove "Bearer " prefix if present
		if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
			tokenString = tokenString[7:]
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func GetAllUsers(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	users, total, err := service.GetAllUsers(pageSize, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 构建响应，包含在线状态
	var result []gin.H
	for _, user := range users {
		result = append(result, gin.H{
			"id":         user.ID,
			"username":   user.Username,
			"nickname":   user.Nickname,
			"avatar":     user.Avatar,
			"online":     cache.IsOnline(user.ID),
			"created_at": user.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"users": result,
		"pagination": gin.H{
			"page":       page,
			"page_size":  pageSize,
			"total":      total,
			"total_page": (total + pageSize - 1) / pageSize,
		},
	})
}

func GetAllMessages(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "20")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	messages, total, err := service.GetAllMessages(pageSize, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 构建响应，包含发送者和接收者信息
	var result []gin.H
	for _, msg := range messages {
		sender, _ := service.GetUserByID(msg.SenderID)
		var senderInfo gin.H
		if sender != nil {
			senderInfo = gin.H{
				"id":       sender.ID,
				"nickname": sender.Nickname,
				"avatar":   sender.Avatar,
			}
		}

		var receiverInfo gin.H
		if msg.ReceiverType == 0 {
			receiver, _ := service.GetUserByID(msg.ReceiverID)
			if receiver != nil {
				receiverInfo = gin.H{
					"id":       receiver.ID,
					"nickname": receiver.Nickname,
					"avatar":   receiver.Avatar,
					"type":     "user",
				}
			}
		} else {
			group, _ := service.GetGroupByID(msg.ReceiverID)
			if group != nil {
				receiverInfo = gin.H{
					"id":   group.ID,
					"name": group.Name,
					"type": "group",
				}
			}
		}

		result = append(result, gin.H{
			"id":            msg.ID,
			"sender":        senderInfo,
			"receiver":      receiverInfo,
			"receiver_type": msg.ReceiverType,
			"content":       msg.Content,
			"type":          msg.Type,
			"status":        msg.Status,
			"created_at":    msg.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"messages": result,
		"pagination": gin.H{
			"page":       page,
			"page_size":  pageSize,
			"total":      total,
			"total_page": (total + pageSize - 1) / pageSize,
		},
	})
}
