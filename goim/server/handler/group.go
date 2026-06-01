package handler

import (
	"goim/server/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CreateGroupRequest struct {
	Name    string `json:"name"`
	OwnerID string `json:"owner_id"`
}

type GroupMemberRequest struct {
	GroupID string `json:"group_id"`
	UserID  string `json:"user_id"`
}

func CreateGroup(c *gin.Context) {
	var req CreateGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	group, err := service.CreateGroup(req.Name, req.OwnerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":      group.ID,
		"name":    group.Name,
		"owner_id": group.OwnerID,
	})
}

func AddGroupMember(c *gin.Context) {
	var req GroupMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := service.AddGroupMember(req.GroupID, req.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "member added successfully"})
}

func RemoveGroupMember(c *gin.Context) {
	var req GroupMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := service.RemoveGroupMember(req.GroupID, req.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "member removed successfully"})
}

func GetGroups(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	groups, err := service.GetGroups(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var result []gin.H
	for _, group := range groups {
		result = append(result, gin.H{
			"id":      group.ID,
			"name":    group.Name,
			"avatar":  group.Avatar,
			"owner_id": group.OwnerID,
		})
	}

	c.JSON(http.StatusOK, result)
}

func GetGroupMembers(c *gin.Context) {
	groupID := c.Param("groupID")
	if groupID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "group_id is required"})
		return
	}

	members, err := service.GetGroupMembers(groupID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var result []gin.H
	for _, member := range members {
		user, _ := service.GetUserByID(member.UserID)
		result = append(result, gin.H{
			"id":       member.ID,
			"group_id": member.GroupID,
			"user_id":  member.UserID,
			"role":     member.Role,
			"username": user.Username,
			"nickname": user.Nickname,
		})
	}

	c.JSON(http.StatusOK, result)
}