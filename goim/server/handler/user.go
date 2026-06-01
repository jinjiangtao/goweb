package handler

import (
	"goim/server/cache"
	"goim/server/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AddFriendRequest struct {
	UserID   string `json:"user_id"`
	FriendID string `json:"friend_id"`
}

func AddFriend(c *gin.Context) {
	var req AddFriendRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := service.AddFriend(req.UserID, req.FriendID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "friend added successfully"})
}

func GetFriends(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	friends, err := service.GetFriends(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var result []gin.H
	for _, friend := range friends {
		result = append(result, gin.H{
			"id":       friend.ID,
			"username": friend.Username,
			"nickname": friend.Nickname,
			"avatar":   friend.Avatar,
			"online":   cache.IsOnline(friend.ID),
		})
	}

	c.JSON(http.StatusOK, result)
}

func GetOnlineStatus(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"online": cache.IsOnline(userID)})
}