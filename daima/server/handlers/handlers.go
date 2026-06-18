package handlers

import (
	"codesnippet/database"
	"codesnippet/models"
	"codesnippet/utils"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误: " + err.Error()})
		return
	}

	var existingUser models.User
	if err := database.DB.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "用户名已存在"})
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "密码加密失败"})
		return
	}

	user := models.User{
		Username: req.Username,
		Password: hashedPassword,
		Email:    req.Email,
		Avatar:   utils.GenerateAvatar(req.Username),
	}

	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "用户创建失败"})
		return
	}

	token, err := utils.GenerateToken(user.ID, user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "生成令牌失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "注册成功",
		"data": gin.H{
			"token": token,
			"user":  user,
		},
	})
}

func Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误: " + err.Error()})
		return
	}

	var user models.User
	if err := database.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "用户名或密码错误"})
		return
	}

	if !utils.CheckPasswordHash(req.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "用户名或密码错误"})
		return
	}

	token, err := utils.GenerateToken(user.ID, user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "生成令牌失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "登录成功",
		"data": gin.H{
			"token": token,
			"user":  user,
		},
	})
}

func GetCurrentUser(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "用户不存在"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "success", "data": user})
}

func CreateSnippet(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req models.SnippetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误: " + err.Error()})
		return
	}

	snippet := models.Snippet{
		Title:       req.Title,
		Description: req.Description,
		Content:     req.Content,
		Language:    req.Language,
		Visibility:  req.Visibility,
		Tags:        req.Tags,
		UserID:      userID.(uint),
	}

	if err := database.DB.Create(&snippet).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建片段失败"})
		return
	}

	database.DB.Preload("User").First(&snippet, snippet.ID)

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "创建成功", "data": snippet})
}

func GetSnippet(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的ID"})
		return
	}

	var snippet models.Snippet
	if err := database.DB.Preload("User").First(&snippet, uint(id)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "片段不存在"})
		return
	}

	userID, exists := c.Get("user_id")
	if snippet.Visibility == "private" {
		if !exists || snippet.UserID != userID.(uint) {
			c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "无权访问此片段"})
			return
		}
	}

	database.DB.Model(&snippet).UpdateColumn("views_count", snippet.ViewsCount+1)
	snippet.ViewsCount++

	response := gin.H{"snippet": snippet}

	if exists {
		var like models.Like
		response["liked"] = database.DB.Where("user_id = ? AND snippet_id = ?", userID, snippet.ID).First(&like).Error == nil

		var fav models.Favorite
		response["favorited"] = database.DB.Where("user_id = ? AND snippet_id = ?", userID, snippet.ID).First(&fav).Error == nil
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "success", "data": response})
}

func UpdateSnippet(c *gin.Context) {
	userID, _ := c.Get("user_id")

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的ID"})
		return
	}

	var snippet models.Snippet
	if err := database.DB.First(&snippet, uint(id)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "片段不存在"})
		return
	}

	if snippet.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "无权修改此片段"})
		return
	}

	var req models.SnippetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误: " + err.Error()})
		return
	}

	snippet.Title = req.Title
	snippet.Description = req.Description
	snippet.Content = req.Content
	snippet.Language = req.Language
	snippet.Visibility = req.Visibility
	snippet.Tags = req.Tags

	if err := database.DB.Save(&snippet).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败"})
		return
	}

	database.DB.Preload("User").First(&snippet, snippet.ID)

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "更新成功", "data": snippet})
}

func DeleteSnippet(c *gin.Context) {
	userID, _ := c.Get("user_id")

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的ID"})
		return
	}

	var snippet models.Snippet
	if err := database.DB.First(&snippet, uint(id)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "片段不存在"})
		return
	}

	if snippet.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "无权删除此片段"})
		return
	}

	database.DB.Where("snippet_id = ?", snippet.ID).Delete(&models.Like{})
	database.DB.Where("snippet_id = ?", snippet.ID).Delete(&models.Favorite{})
	database.DB.Where("snippet_id = ?", snippet.ID).Delete(&models.Comment{})
	database.DB.Delete(&snippet)

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "删除成功"})
}

func ListSnippets(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	language := c.Query("language")
	sort := c.DefaultQuery("sort", "latest")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	query := database.DB.Model(&models.Snippet{}).Where("visibility = ?", "public")

	if language != "" && language != "all" {
		query = query.Where("language = ?", language)
	}

	var total int64
	query.Count(&total)

	var snippets []models.Snippet
	offset := (page - 1) * pageSize

	dbQuery := query.Preload("User").Offset(offset).Limit(pageSize)
	switch sort {
	case "likes":
		dbQuery = dbQuery.Order("likes_count DESC, created_at DESC")
	case "views":
		dbQuery = dbQuery.Order("views_count DESC, created_at DESC")
	case "favs":
		dbQuery = dbQuery.Order("favs_count DESC, created_at DESC")
	default:
		dbQuery = dbQuery.Order("created_at DESC")
	}

	dbQuery.Find(&snippets)

	c.JSON(http.StatusOK, gin.H{
		"code":      200,
		"message":   "success",
		"data":      snippets,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func SearchSnippets(c *gin.Context) {
	keyword := c.Query("q")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	language := c.Query("language")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	query := database.DB.Model(&models.Snippet{}).Where("visibility = ?", "public")

	if keyword != "" {
		like := "%" + keyword + "%"
		query = query.Where("(title LIKE ? OR description LIKE ? OR content LIKE ? OR tags LIKE ?)", like, like, like, like)
	}

	if language != "" && language != "all" {
		query = query.Where("language = ?", language)
	}

	var total int64
	query.Count(&total)

	var snippets []models.Snippet
	offset := (page - 1) * pageSize
	query.Preload("User").Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&snippets)

	c.JSON(http.StatusOK, gin.H{
		"code":      200,
		"message":   "success",
		"data":      snippets,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func GetMySnippets(c *gin.Context) {
	userID, _ := c.Get("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	visibility := c.DefaultQuery("visibility", "all")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	query := database.DB.Model(&models.Snippet{}).Where("user_id = ?", userID)
	if visibility != "all" {
		query = query.Where("visibility = ?", visibility)
	}

	var total int64
	query.Count(&total)

	var snippets []models.Snippet
	offset := (page - 1) * pageSize
	query.Preload("User").Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&snippets)

	c.JSON(http.StatusOK, gin.H{
		"code":      200,
		"message":   "success",
		"data":      snippets,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func GetMyFavorites(c *gin.Context) {
	userID, _ := c.Get("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	var favs []models.Favorite
	query := database.DB.Model(&models.Favorite{}).Where("user_id = ?", userID)

	var total int64
	query.Count(&total)

	offset := (page - 1) * pageSize
	query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&favs)

	snippetIDs := make([]uint, 0, len(favs))
	for _, f := range favs {
		snippetIDs = append(snippetIDs, f.SnippetID)
	}

	var snippets []models.Snippet
	if len(snippetIDs) > 0 {
		database.DB.Preload("User").Where("id IN ?", snippetIDs).Order("created_at DESC").Find(&snippets)
	}

	c.JSON(http.StatusOK, gin.H{
		"code":      200,
		"message":   "success",
		"data":      snippets,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func ToggleLike(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的ID"})
		return
	}

	var snippet models.Snippet
	if err := database.DB.First(&snippet, uint(id)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "片段不存在"})
		return
	}

	var like models.Like
	result := database.DB.Where("user_id = ? AND snippet_id = ?", userID, snippet.ID).First(&like)

	liked := false
	if result.Error == nil {
		database.DB.Delete(&like)
		snippet.LikesCount--
		liked = false
	} else {
		like = models.Like{UserID: userID.(uint), SnippetID: snippet.ID}
		database.DB.Create(&like)
		snippet.LikesCount++
		liked = true
	}

	database.DB.Model(&snippet).UpdateColumn("likes_count", snippet.LikesCount)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data": gin.H{
			"liked":       liked,
			"likes_count": snippet.LikesCount,
		},
	})
}

func ToggleFavorite(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的ID"})
		return
	}

	var snippet models.Snippet
	if err := database.DB.First(&snippet, uint(id)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "片段不存在"})
		return
	}

	var fav models.Favorite
	result := database.DB.Where("user_id = ? AND snippet_id = ?", userID, snippet.ID).First(&fav)

	favorited := false
	if result.Error == nil {
		database.DB.Delete(&fav)
		snippet.FavsCount--
		favorited = false
	} else {
		fav = models.Favorite{UserID: userID.(uint), SnippetID: snippet.ID}
		database.DB.Create(&fav)
		snippet.FavsCount++
		favorited = true
	}

	database.DB.Model(&snippet).UpdateColumn("favs_count", snippet.FavsCount)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data": gin.H{
			"favorited":  favorited,
			"favs_count": snippet.FavsCount,
		},
	})
}

func GetComments(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的ID"})
		return
	}

	var comments []models.Comment
	database.DB.Preload("User").Where("snippet_id = ?", uint(id)).Order("created_at DESC").Find(&comments)

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "success", "data": comments})
}

func CreateComment(c *gin.Context) {
	userID, _ := c.Get("user_id")

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的ID"})
		return
	}

	var snippet models.Snippet
	if err := database.DB.First(&snippet, uint(id)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "片段不存在"})
		return
	}

	var req models.CommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误: " + err.Error()})
		return
	}

	comment := models.Comment{
		Content:   req.Content,
		UserID:    userID.(uint),
		SnippetID: snippet.ID,
	}

	if err := database.DB.Create(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "评论失败"})
		return
	}

	snippet.CommentsCount++
	database.DB.Model(&snippet).UpdateColumn("comments_count", snippet.CommentsCount)

	database.DB.Preload("User").First(&comment, comment.ID)

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "评论成功", "data": comment})
}

func DeleteComment(c *gin.Context) {
	userID, _ := c.Get("user_id")
	commentID, err := strconv.ParseUint(c.Param("comment_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的ID"})
		return
	}

	var comment models.Comment
	if err := database.DB.First(&comment, uint(commentID)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "评论不存在"})
		return
	}

	if comment.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "无权删除此评论"})
		return
	}

	database.DB.Delete(&comment)

	var snippet models.Snippet
	database.DB.First(&snippet, comment.SnippetID)
	snippet.CommentsCount--
	if snippet.CommentsCount < 0 {
		snippet.CommentsCount = 0
	}
	database.DB.Model(&snippet).UpdateColumn("comments_count", snippet.CommentsCount)

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "删除成功"})
}

func GetLanguages(c *gin.Context) {
	languages := []string{
		"javascript", "typescript", "python", "go", "java", "c", "cpp",
		"csharp", "ruby", "php", "rust", "swift", "kotlin", "html", "css",
		"scss", "json", "yaml", "markdown", "sql", "shell", "bash", "dart",
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "success", "data": languages})
}

func ForkSnippet(c *gin.Context) {
	userID, _ := c.Get("user_id")

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的ID"})
		return
	}

	var snippet models.Snippet
	if err := database.DB.First(&snippet, uint(id)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "片段不存在"})
		return
	}

	userIDVal, exists := c.Get("user_id")
	if snippet.Visibility == "private" {
		if !exists || snippet.UserID != userIDVal.(uint) {
			c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "无权Fork此片段"})
			return
		}
	}

	newSnippet := models.Snippet{
		Title:       snippet.Title + " (Fork)",
		Description: snippet.Description,
		Content:     snippet.Content,
		Language:    snippet.Language,
		Visibility:  "public",
		Tags:        snippet.Tags,
		UserID:      userID.(uint),
	}

	if err := database.DB.Create(&newSnippet).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Fork失败"})
		return
	}

	database.DB.Preload("User").First(&newSnippet, newSnippet.ID)

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "Fork成功", "data": newSnippet})
}

func DiffSnippets(c *gin.Context) {
	type DiffRequest struct {
		Original string `json:"original" binding:"required"`
		Modified string `json:"modified" binding:"required"`
	}

	var req DiffRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误: " + err.Error()})
		return
	}

	origLines := strings.Split(req.Original, "\n")
	modLines := strings.Split(req.Modified, "\n")

	diff := computeDiff(origLines, modLines)

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "success", "data": diff})
}

func computeDiff(orig, mod []string) []map[string]interface{} {
	result := []map[string]interface{}{}
	m, n := len(orig), len(mod)
	dp := make([][]int, m+1)
	for i := range dp {
		dp[i] = make([]int, n+1)
	}
	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			if orig[i-1] == mod[j-1] {
				dp[i][j] = dp[i-1][j-1] + 1
			} else {
				dp[i][j] = max(dp[i-1][j], dp[i][j-1])
			}
		}
	}
	i, j := m, n
	for i > 0 || j > 0 {
		if i > 0 && j > 0 && orig[i-1] == mod[j-1] {
			result = append([]map[string]interface{}{{"type": "equal", "line": orig[i-1], "orig": i, "mod": j}}, result...)
			i--
			j--
		} else if j > 0 && (i == 0 || dp[i][j-1] >= dp[i-1][j]) {
			result = append([]map[string]interface{}{{"type": "add", "line": mod[j-1], "mod": j}}, result...)
			j--
		} else if i > 0 {
			result = append([]map[string]interface{}{{"type": "remove", "line": orig[i-1], "orig": i}}, result...)
			i--
		}
	}
	return result
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
