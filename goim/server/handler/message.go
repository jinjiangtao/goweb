package handler

import (
	"goim/server/cache"
	"goim/server/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type MarkReadRequest struct {
	UserID    string `json:"user_id"`
	MessageID string `json:"message_id"`
}

func GetMessageHistory(c *gin.Context) {
	userID := c.Query("user_id")
	targetID := c.Query("target_id")
	targetType, _ := strconv.Atoi(c.Query("target_type"))
	limit, _ := strconv.Atoi(c.Query("limit"))
	offset, _ := strconv.Atoi(c.Query("offset"))

	if userID == "" || targetID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id and target_id are required"})
		return
	}

	if limit == 0 {
		limit = 20
	}

	messages, err := service.GetMessageHistory(userID, targetID, targetType, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, messages)
}

func MarkMessageRead(c *gin.Context) {
	var req MarkReadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := service.MarkMessagesRead(req.UserID, req.MessageID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "message marked as read"})
}

func GetUnreadCounts(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	counts := cache.GetAllUnreadCounts(userID)
	c.JSON(http.StatusOK, counts)
}