package handlers

import (
	"math"
	"net/http"
	"strconv"
	"time"

	"zhishiku-server/middleware"
	"zhishiku-server/models"

	"github.com/gin-gonic/gin"
)

type ArticleRequest struct {
	Title      string `json:"title" binding:"required"`
	Content    string `json:"content" binding:"required"`
	CategoryID uint   `json:"category_id" binding:"required"`
	Status     string `json:"status" binding:"required,oneof=draft published"`
}

func GetArticles(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	keyword := c.Query("keyword")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	query := db.Model(&models.Article{}).Where("status = ?", "published")

	if keyword != "" {
		query = query.Where("title LIKE ?", "%"+keyword+"%")
	}

	var total int64
	query.Count(&total)

	var articles []models.Article
	offset := (page - 1) * pageSize
	if err := query.Preload("Category").Preload("Author").
		Order("published_at DESC, created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&articles).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取文章列表失败"})
		return
	}

	var list []models.ArticleListResponse
	for _, a := range articles {
		item := models.ArticleListResponse{
			ID:           a.ID,
			Title:        a.Title,
			CategoryID:   a.CategoryID,
			CategoryName: a.Category.Name,
			AuthorID:     a.AuthorID,
			AuthorName:   a.Author.Nickname,
			Status:       a.Status,
			PublishedAt:  a.PublishedAt,
			CreatedAt:    a.CreatedAt,
			UpdatedAt:    a.UpdatedAt,
		}
		list = append(list, item)
	}

	c.JSON(http.StatusOK, gin.H{
		"list":        list,
		"total":       total,
		"page":        page,
		"page_size":   pageSize,
		"total_pages": int(math.Ceil(float64(total) / float64(pageSize))),
	})
}

func GetArticle(c *gin.Context) {
	id := c.Param("id")

	var article models.Article
	if err := db.Preload("Category").Preload("Author").First(&article, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":            article.ID,
		"title":         article.Title,
		"content":       article.Content,
		"category_id":   article.CategoryID,
		"category_name": article.Category.Name,
		"author_id":     article.AuthorID,
		"author_name":   article.Author.Nickname,
		"status":        article.Status,
		"published_at":  article.PublishedAt,
		"created_at":    article.CreatedAt,
		"updated_at":    article.UpdatedAt,
	})
}

func CreateArticle(c *gin.Context) {
	var req ArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	userID := middleware.GetUserID(c)

	article := models.Article{
		Title:      req.Title,
		Content:    req.Content,
		CategoryID: req.CategoryID,
		AuthorID:   userID,
		Status:     req.Status,
	}

	if req.Status == "published" {
		now := time.Now()
		article.PublishedAt = &now
	}

	if err := db.Create(&article).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建文章失败"})
		return
	}

	var maxVer int
	db.Model(&models.ArticleVersion{}).Where("article_id = ?", article.ID).Select("COALESCE(MAX(version_number), 0)").Scan(&maxVer)
	version := models.ArticleVersion{
		ArticleID:       article.ID,
		VersionNumber:   maxVer + 1,
		TitleSnapshot:   article.Title,
		ContentSnapshot: article.Content,
		CategoryIDSnap:  article.CategoryID,
		StatusSnapshot:  article.Status,
		CreatorID:       userID,
	}
	db.Create(&version)

	c.JSON(http.StatusCreated, gin.H{"id": article.ID, "message": "创建成功"})
}

func UpdateArticle(c *gin.Context) {
	id := c.Param("id")

	var article models.Article
	if err := db.First(&article, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
		return
	}

	userID := middleware.GetUserID(c)
	role := middleware.GetRole(c)

	if article.AuthorID != userID && role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权编辑此文章"})
		return
	}

	var req ArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	wasDraft := article.Status == "draft"

	article.Title = req.Title
	article.Content = req.Content
	article.CategoryID = req.CategoryID
	article.Status = req.Status

	if req.Status == "published" && wasDraft {
		now := time.Now()
		article.PublishedAt = &now
	}

	if req.Status == "draft" {
		article.PublishedAt = nil
	}

	if err := db.Save(&article).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新文章失败"})
		return
	}

	var maxVer int
	db.Model(&models.ArticleVersion{}).Where("article_id = ?", article.ID).Select("COALESCE(MAX(version_number), 0)").Scan(&maxVer)
	newVersion := models.ArticleVersion{
		ArticleID:       article.ID,
		VersionNumber:   maxVer + 1,
		TitleSnapshot:   article.Title,
		ContentSnapshot: article.Content,
		CategoryIDSnap:  article.CategoryID,
		StatusSnapshot:  article.Status,
		CreatorID:       userID,
	}
	db.Create(&newVersion)

	cleanupOldVersions(article.ID)

	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}

func DeleteArticle(c *gin.Context) {
	id := c.Param("id")

	var article models.Article
	if err := db.First(&article, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
		return
	}

	userID := middleware.GetUserID(c)
	role := middleware.GetRole(c)

	if article.AuthorID != userID && role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权删除此文章"})
		return
	}

	if err := db.Delete(&article).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除文章失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

func GetMyArticles(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	keyword := c.Query("keyword")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	userID := middleware.GetUserID(c)

	query := db.Model(&models.Article{}).Where("author_id = ?", userID)

	if keyword != "" {
		query = query.Where("title LIKE ?", "%"+keyword+"%")
	}

	var total int64
	query.Count(&total)

	var articles []models.Article
	offset := (page - 1) * pageSize
	if err := query.Preload("Category").
		Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&articles).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取文章列表失败"})
		return
	}

	var list []models.ArticleListResponse
	for _, a := range articles {
		item := models.ArticleListResponse{
			ID:           a.ID,
			Title:        a.Title,
			CategoryID:   a.CategoryID,
			CategoryName: a.Category.Name,
			AuthorID:     a.AuthorID,
			AuthorName:   "",
			Status:       a.Status,
			PublishedAt:  a.PublishedAt,
			CreatedAt:    a.CreatedAt,
			UpdatedAt:    a.UpdatedAt,
		}
		list = append(list, item)
	}

	c.JSON(http.StatusOK, gin.H{
		"list":        list,
		"total":       total,
		"page":        page,
		"page_size":   pageSize,
		"total_pages": int(math.Ceil(float64(total) / float64(pageSize))),
	})
}

func SearchArticles(c *gin.Context) {
	keyword := c.Query("keyword")
	if keyword == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请输入搜索关键词"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	query := db.Model(&models.Article{}).Where("status = ? AND title LIKE ?", "published", "%"+keyword+"%")

	var total int64
	query.Count(&total)

	var articles []models.Article
	offset := (page - 1) * pageSize
	if err := query.Preload("Category").Preload("Author").
		Order("published_at DESC, created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&articles).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "搜索失败"})
		return
	}

	var list []models.ArticleListResponse
	for _, a := range articles {
		item := models.ArticleListResponse{
			ID:           a.ID,
			Title:        a.Title,
			CategoryID:   a.CategoryID,
			CategoryName: a.Category.Name,
			AuthorID:     a.AuthorID,
			AuthorName:   a.Author.Nickname,
			Status:       a.Status,
			PublishedAt:  a.PublishedAt,
			CreatedAt:    a.CreatedAt,
			UpdatedAt:    a.UpdatedAt,
		}
		list = append(list, item)
	}

	c.JSON(http.StatusOK, gin.H{
		"list":        list,
		"total":       total,
		"page":        page,
		"page_size":   pageSize,
		"total_pages": int(math.Ceil(float64(total) / float64(pageSize))),
	})
}
