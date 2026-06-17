package handlers

import (
	"math"
	"net/http"
	"strconv"

	"zhishiku-server/diff"
	"zhishiku-server/middleware"
	"zhishiku-server/models"

	"github.com/gin-gonic/gin"
)

func GetVersions(c *gin.Context) {
	articleID := c.Param("id")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var article models.Article
	if err := db.First(&article, articleID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
		return
	}

	var total int64
	db.Model(&models.ArticleVersion{}).Where("article_id = ?", articleID).Count(&total)

	var versions []models.ArticleVersion
	offset := (page - 1) * pageSize
	if err := db.Where("article_id = ?", articleID).
		Order("version_number DESC").
		Offset(offset).Limit(pageSize).
		Find(&versions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取版本列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"list":        versions,
		"total":       total,
		"page":        page,
		"page_size":   pageSize,
		"total_pages": int(math.Ceil(float64(total) / float64(pageSize))),
	})
}

func GetVersionDetail(c *gin.Context) {
	versionID := c.Param("versionId")

	var version models.ArticleVersion
	if err := db.First(&version, versionID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "版本不存在"})
		return
	}

	var category models.Category
	categoryName := ""
	if err := db.First(&category, version.CategoryIDSnap).Error; err == nil {
		categoryName = category.Name
	}

	c.JSON(http.StatusOK, gin.H{
		"id":                version.ID,
		"article_id":        version.ArticleID,
		"version_number":    version.VersionNumber,
		"title_snapshot":    version.TitleSnapshot,
		"content_snapshot":  version.ContentSnapshot,
		"category_id_snap":  version.CategoryIDSnap,
		"category_name":     categoryName,
		"status_snapshot":   version.StatusSnapshot,
		"is_rollback":       version.IsRollback,
		"rollback_from_ver": version.RollbackFromVer,
		"created_at":        version.CreatedAt,
	})
}

func CompareVersions(c *gin.Context) {
	articleID := c.Param("id")

	oldVerID := c.Query("old_version_id")
	newVerID := c.Query("new_version_id")
	if oldVerID == "" || newVerID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请提供对比版本ID"})
		return
	}

	var oldVer, newVer models.ArticleVersion
	if err := db.First(&oldVer, oldVerID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "旧版本不存在"})
		return
	}
	if err := db.First(&newVer, newVerID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "新版本不存在"})
		return
	}

	if oldVer.ArticleID != newVer.ArticleID || strconv.FormatUint(uint64(oldVer.ArticleID), 10) != articleID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "版本不属于同一篇文章"})
		return
	}

	result := diff.DiffResult{}

	titleChanged := oldVer.TitleSnapshot != newVer.TitleSnapshot
	result.TitleDiff = &diff.TextDiff{
		OldValue: oldVer.TitleSnapshot,
		NewValue: newVer.TitleSnapshot,
		Changed:  titleChanged,
	}

	var oldCat, newCat models.Category
	oldCatName, newCatName := "", ""
	if err := db.First(&oldCat, oldVer.CategoryIDSnap).Error; err == nil {
		oldCatName = oldCat.Name
	}
	if err := db.First(&newCat, newVer.CategoryIDSnap).Error; err == nil {
		newCatName = newCat.Name
	}

	catChanged := oldVer.CategoryIDSnap != newVer.CategoryIDSnap
	result.CatDiff = &diff.CategoryDiff{
		OldID:   oldVer.CategoryIDSnap,
		OldName: oldCatName,
		NewID:   newVer.CategoryIDSnap,
		NewName: newCatName,
		Changed: catChanged,
	}

	result.Lines = diff.ComputeDiff(oldVer.ContentSnapshot, newVer.ContentSnapshot)
	result.HasChange = titleChanged || catChanged || diff.HasContentChange(oldVer.ContentSnapshot, newVer.ContentSnapshot)

	c.JSON(http.StatusOK, result)
}

func RollbackVersion(c *gin.Context) {
	versionID := c.Param("versionId")

	var targetVer models.ArticleVersion
	if err := db.First(&targetVer, versionID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "该版本已不存在"})
		return
	}

	articleIDStr := c.Param("id")
	var article models.Article
	if err := db.First(&article, articleIDStr).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
		return
	}

	if article.ID != targetVer.ArticleID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "版本不属于该文章"})
		return
	}

	userID := middleware.GetUserID(c)
	role := middleware.GetRole(c)
	if article.AuthorID != userID && role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权操作此文章"})
		return
	}

	article.Title = targetVer.TitleSnapshot
	article.Content = targetVer.ContentSnapshot
	article.CategoryID = targetVer.CategoryIDSnap
	article.Status = targetVer.StatusSnapshot

	if err := db.Save(&article).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "回滚失败"})
		return
	}

	var maxVer int
	db.Model(&models.ArticleVersion{}).Where("article_id = ?", article.ID).Select("COALESCE(MAX(version_number), 0)").Scan(&maxVer)

	rollbackFrom := targetVer.VersionNumber
	newVersion := models.ArticleVersion{
		ArticleID:       article.ID,
		VersionNumber:   maxVer + 1,
		TitleSnapshot:   targetVer.TitleSnapshot,
		ContentSnapshot: targetVer.ContentSnapshot,
		CategoryIDSnap:  targetVer.CategoryIDSnap,
		StatusSnapshot:  targetVer.StatusSnapshot,
		CreatorID:       userID,
		IsRollback:      true,
		RollbackFromVer: &rollbackFrom,
	}
	db.Create(&newVersion)

	cleanupOldVersions(article.ID)

	c.JSON(http.StatusOK, gin.H{
		"message":        "回滚成功",
		"new_version_id": newVersion.ID,
	})
}

func DeleteVersion(c *gin.Context) {
	versionID := c.Param("versionId")

	var version models.ArticleVersion
	if err := db.First(&version, versionID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "版本不存在"})
		return
	}

	var article models.Article
	if err := db.First(&article, version.ArticleID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
		return
	}

	var currentMaxVer int
	db.Model(&models.ArticleVersion{}).Where("article_id = ?", article.ID).Select("COALESCE(MAX(version_number), 0)").Scan(&currentMaxVer)
	if version.VersionNumber == currentMaxVer {
		c.JSON(http.StatusBadRequest, gin.H{"error": "当前版本不可删除"})
		return
	}

	role := middleware.GetRole(c)
	if role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "需要管理员权限"})
		return
	}

	if err := db.Delete(&version).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除版本失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "版本已删除"})
}

func cleanupOldVersions(articleID uint) {
	var count int64
	db.Model(&models.ArticleVersion{}).Where("article_id = ?", articleID).Count(&count)

	if count > 50 {
		var oldest []models.ArticleVersion
		db.Where("article_id = ?", articleID).
			Order("version_number ASC").
			Limit(int(count - 50)).
			Find(&oldest)

		for _, v := range oldest {
			db.Delete(&v)
		}
	}
}
