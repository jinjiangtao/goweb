package handlers

import (
	"biji/database"
	"biji/models"
	"biji/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CreateTagRequest struct {
	Name     string  `json:"name" binding:"required"`
	Color    string  `json:"color"`
	ParentID *string `json:"parentId"`
	GroupID  *string `json:"groupId"`
	Sort     int     `json:"sort"`
}

type UpdateTagRequest struct {
	Name     string  `json:"name"`
	Color    string  `json:"color"`
	ParentID *string `json:"parentId"`
	GroupID  *string `json:"groupId"`
	Sort     *int    `json:"sort"`
}

type CreateTagGroupRequest struct {
	Name string `json:"name" binding:"required"`
	Sort int    `json:"sort"`
}

type UpdateTagGroupRequest struct {
	Name string `json:"name"`
	Sort *int   `json:"sort"`
}

type ReorderTagsRequest struct {
	TagIDs []string `json:"tagIds" binding:"required"`
}

func CreateTag(c *gin.Context) {
	var req CreateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Color == "" {
		req.Color = "#409EFF"
	}

	var maxSort int
	query := database.DB.Model(&models.Tag{})
	if req.ParentID != nil {
		query = query.Where("parent_id = ?", *req.ParentID)
	} else {
		query = query.Where("parent_id IS NULL")
	}
	query.Select("COALESCE(MAX(sort), 0)").Scan(&maxSort)

	tag := models.Tag{
		ID:       utils.GenerateID(),
		Name:     req.Name,
		Color:    req.Color,
		ParentID: req.ParentID,
		GroupID:  req.GroupID,
		Sort:     maxSort + 1,
	}

	if err := database.DB.Create(&tag).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, tag)
}

func GetTags(c *gin.Context) {
	flat := c.DefaultQuery("flat", "false")

	if flat == "true" {
		var tags []models.Tag
		database.DB.Order("sort ASC, created_at ASC").Find(&tags)
		c.JSON(http.StatusOK, tags)
		return
	}

	var rootTags []models.Tag
	database.DB.Where("parent_id IS NULL").Order("sort ASC, created_at ASC").Find(&rootTags)

	for i := range rootTags {
		loadChildren(&rootTags[i])
	}

	c.JSON(http.StatusOK, rootTags)
}

func loadChildren(tag *models.Tag) {
	var children []models.Tag
	database.DB.Where("parent_id = ?", tag.ID).Order("sort ASC, created_at ASC").Find(&children)
	for i := range children {
		loadChildren(&children[i])
	}
	tag.Children = children
}

func GetTag(c *gin.Context) {
	id := c.Param("id")
	var tag models.Tag
	if err := database.DB.First(&tag, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tag not found"})
		return
	}
	loadChildren(&tag)
	c.JSON(http.StatusOK, tag)
}

func UpdateTag(c *gin.Context) {
	id := c.Param("id")
	var tag models.Tag
	if err := database.DB.First(&tag, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tag not found"})
		return
	}

	var req UpdateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Name != "" {
		tag.Name = req.Name
	}
	if req.Color != "" {
		tag.Color = req.Color
	}
	if req.ParentID != nil {
		tag.ParentID = req.ParentID
	}
	if req.GroupID != nil {
		tag.GroupID = req.GroupID
	}
	if req.Sort != nil {
		tag.Sort = *req.Sort
	}

	if err := database.DB.Save(&tag).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tag)
}

func DeleteTag(c *gin.Context) {
	id := c.Param("id")
	var tag models.Tag
	if err := database.DB.First(&tag, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tag not found"})
		return
	}

	deleteTagAndChildren(id)

	c.JSON(http.StatusOK, gin.H{"message": "Tag deleted"})
}

func deleteTagAndChildren(tagID string) {
	var children []models.Tag
	database.DB.Where("parent_id = ?", tagID).Find(&children)
	for _, child := range children {
		deleteTagAndChildren(child.ID)
	}

	var tag models.Tag
	database.DB.First(&tag, "id = ?", tagID)
	database.DB.Model(&tag).Association("Notes").Clear()
	database.DB.Delete(&tag)
}

func ReorderTags(c *gin.Context) {
	var req ReorderTagsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for i, tagID := range req.TagIDs {
		database.DB.Model(&models.Tag{}).Where("id = ?", tagID).Update("sort", i+1)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tags reordered"})
}

func CreateTagGroup(c *gin.Context) {
	var req CreateTagGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var maxSort int
	database.DB.Model(&models.TagGroup{}).Select("COALESCE(MAX(sort), 0)").Scan(&maxSort)

	group := models.TagGroup{
		ID:   utils.GenerateID(),
		Name: req.Name,
		Sort: maxSort + 1,
	}

	if err := database.DB.Create(&group).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, group)
}

func GetTagGroups(c *gin.Context) {
	var groups []models.TagGroup
	database.DB.Order("sort ASC, created_at ASC").Preload("Tags").Find(&groups)
	for i := range groups {
		for j := range groups[i].Tags {
			groups[i].Tags[j].Children = nil
		}
	}
	c.JSON(http.StatusOK, groups)
}

func UpdateTagGroup(c *gin.Context) {
	id := c.Param("id")
	var group models.TagGroup
	if err := database.DB.First(&group, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tag group not found"})
		return
	}

	var req UpdateTagGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Name != "" {
		group.Name = req.Name
	}
	if req.Sort != nil {
		group.Sort = *req.Sort
	}

	if err := database.DB.Save(&group).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, group)
}

func DeleteTagGroup(c *gin.Context) {
	id := c.Param("id")
	var group models.TagGroup
	if err := database.DB.First(&group, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tag group not found"})
		return
	}

	database.DB.Model(&models.Tag{}).Where("group_id = ?", id).Update("group_id", nil)
	database.DB.Delete(&group)

	c.JSON(http.StatusOK, gin.H{"message": "Tag group deleted"})
}

func GetTagNoteCount(c *gin.Context) {
	id := c.Param("id")
	var count int64
	database.DB.Table("note_tags").Where("tag_id = ?", id).Count(&count)
	c.JSON(http.StatusOK, gin.H{"count": count, "tagId": id})
}

func init() {
	_ = strconv.Itoa
}
