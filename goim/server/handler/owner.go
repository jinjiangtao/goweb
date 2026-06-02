package handler

import (
	"goim/server/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CreateOwnerRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Avatar      string `json:"avatar"`
}

type JoinOwnerRequest struct {
	OwnerID string `json:"owner_id"`
	UserID  string `json:"user_id"`
}

type RemoveOwnerMemberRequest struct {
	OwnerID string `json:"owner_id"`
	UserID  string `json:"user_id"`
}

func CreateOwner(c *gin.Context) {
	var req CreateOwnerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	owner, err := service.CreateOwner(req.Name, req.Description, req.Avatar)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":          owner.ID,
		"name":        owner.Name,
		"description": owner.Description,
		"avatar":      owner.Avatar,
	})
}

func GetAllOwners(c *gin.Context) {
	owners, err := service.GetAllOwners()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var result []gin.H
	for _, owner := range owners {
		result = append(result, gin.H{
			"id":          owner.ID,
			"name":        owner.Name,
			"description": owner.Description,
			"avatar":      owner.Avatar,
			"created_at":  owner.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, result)
}

func GetOwnerByID(c *gin.Context) {
	ownerID := c.Param("ownerID")
	if ownerID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "owner_id is required"})
		return
	}

	owner, err := service.GetOwnerByID(ownerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if owner == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "owner not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":          owner.ID,
		"name":        owner.Name,
		"description": owner.Description,
		"avatar":      owner.Avatar,
	})
}

func DeleteOwner(c *gin.Context) {
	ownerID := c.Param("ownerID")
	if ownerID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "owner_id is required"})
		return
	}

	err := service.DeleteOwner(ownerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "owner deleted successfully"})
}

func JoinOwner(c *gin.Context) {
	var req JoinOwnerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := service.JoinOwner(req.OwnerID, req.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "joined successfully"})
}

func LeaveOwner(c *gin.Context) {
	var req JoinOwnerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := service.LeaveOwner(req.OwnerID, req.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "left successfully"})
}

func GetOwnerMembers(c *gin.Context) {
	ownerID := c.Param("ownerID")
	if ownerID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "owner_id is required"})
		return
	}

	members, err := service.GetOwnerMembers(ownerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var result []gin.H
	for _, member := range members {
		user, _ := service.GetUserByID(member.UserID)
		result = append(result, gin.H{
			"id":        member.ID,
			"owner_id":  member.OwnerID,
			"user_id":   member.UserID,
			"username":  user.Username,
			"nickname":  user.Nickname,
			"joined_at": member.JoinedAt,
		})
	}

	c.JSON(http.StatusOK, result)
}

func RemoveOwnerMember(c *gin.Context) {
	var req RemoveOwnerMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := service.LeaveOwner(req.OwnerID, req.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "member removed successfully"})
}

func GetOwnersByUserID(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	owners, err := service.GetOwnersByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var result []gin.H
	for _, owner := range owners {
		result = append(result, gin.H{
			"id":          owner.ID,
			"name":        owner.Name,
			"description": owner.Description,
			"avatar":      owner.Avatar,
		})
	}

	c.JSON(http.StatusOK, result)
}