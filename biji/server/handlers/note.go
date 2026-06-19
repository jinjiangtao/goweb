package handlers

import (
	"biji/database"
	"biji/models"
	"biji/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type CreateNoteRequest struct {
	Title       string   `json:"title"`
	Content     string   `json:"content"`
	TagIDs      []string `json:"tagIds"`
	IsEncrypted bool     `json:"isEncrypted"`
	Password    string   `json:"password"`
}

type UpdateNoteRequest struct {
	Title       string   `json:"title"`
	Content     string   `json:"content"`
	TagIDs      []string `json:"tagIds"`
	IsEncrypted *bool    `json:"isEncrypted"`
	Password    *string  `json:"password"`
}

type VerifyPasswordRequest struct {
	Password string `json:"password"`
}

func CreateNote(c *gin.Context) {
	var req CreateNoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	note := models.Note{
		ID:          utils.GenerateID(),
		Title:       req.Title,
		Content:     req.Content,
		PlainText:   utils.StripHTMLTags(req.Content),
		IsEncrypted: req.IsEncrypted,
	}

	if req.IsEncrypted && req.Password != "" {
		salt := utils.GenerateSalt()
		note.Salt = salt
		note.Password = utils.HashPassword(req.Password, salt)
		encryptedContent, err := utils.EncryptContent(req.Content, req.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Encryption failed"})
			return
		}
		note.Content = encryptedContent
		note.PlainText = ""
	}

	if len(req.TagIDs) > 0 {
		var tags []models.Tag
		database.DB.Where("id IN ?", req.TagIDs).Find(&tags)
		note.Tags = tags
	}

	if err := database.DB.Create(&note).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, note)
}

func GetNotes(c *gin.Context) {
	status := c.DefaultQuery("status", "normal")
	keyword := c.Query("keyword")
	tagIDs := c.QueryArray("tagIds")
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")

	db := database.DB.Preload("Tags")

	switch status {
	case "archived":
		db = db.Where("is_archived = ? AND is_deleted = ?", true, false)
	case "deleted":
		db = db.Where("is_deleted = ?", true)
	default:
		db = db.Where("is_archived = ? AND is_deleted = ?", false, false)
	}

	if keyword != "" {
		db = db.Where("(title LIKE ? OR plain_text LIKE ?)", "%"+keyword+"%", "%"+keyword+"%")
	}

	if len(tagIDs) > 0 {
		noteIDs := []string{}
		database.DB.Table("note_tags").Where("tag_id IN ?", tagIDs).Pluck("note_id", &noteIDs)
		if len(noteIDs) > 0 {
			db = db.Where("id IN ?", noteIDs)
		} else {
			db = db.Where("1 = 0")
		}
	}

	if startDate != "" {
		if t, err := time.Parse("2006-01-02", startDate); err == nil {
			db = db.Where("created_at >= ?", t)
		}
	}
	if endDate != "" {
		if t, err := time.Parse("2006-01-02", endDate); err == nil {
			db = db.Where("created_at <= ?", t.Add(24*time.Hour))
		}
	}

	db = db.Order("is_pinned DESC, updated_at DESC")

	var notes []models.Note
	if err := db.Find(&notes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, notes)
}

func GetNote(c *gin.Context) {
	id := c.Param("id")
	var note models.Note
	if err := database.DB.Preload("Tags").First(&note, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Note not found"})
		return
	}

	if note.IsEncrypted {
		note.Content = ""
		note.PlainText = ""
	}

	c.JSON(http.StatusOK, note)
}

func VerifyNotePassword(c *gin.Context) {
	id := c.Param("id")
	var req VerifyPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var note models.Note
	if err := database.DB.First(&note, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Note not found"})
		return
	}

	if !note.IsEncrypted {
		c.JSON(http.StatusOK, gin.H{"content": note.Content, "plainText": note.PlainText})
		return
	}

	if !utils.VerifyPassword(req.Password, note.Salt, note.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
		return
	}

	decryptedContent, err := utils.DecryptContent(note.Content, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Decryption failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"content":   decryptedContent,
		"plainText": utils.StripHTMLTags(decryptedContent),
	})
}

func UpdateNote(c *gin.Context) {
	id := c.Param("id")
	var note models.Note
	if err := database.DB.First(&note, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Note not found"})
		return
	}

	var req UpdateNoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Title != "" {
		note.Title = req.Title
	}

	if req.Content != "" {
		if note.IsEncrypted {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Encrypted note must be decrypted first"})
			return
		}
		note.Content = req.Content
		note.PlainText = utils.StripHTMLTags(req.Content)
	}

	if req.IsEncrypted != nil && req.Password != nil {
		if *req.IsEncrypted && *req.Password != "" {
			salt := utils.GenerateSalt()
			note.Salt = salt
			note.Password = utils.HashPassword(*req.Password, salt)
			encryptedContent, err := utils.EncryptContent(note.Content, *req.Password)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Encryption failed"})
				return
			}
			note.Content = encryptedContent
			note.PlainText = ""
			note.IsEncrypted = true
		} else if !*req.IsEncrypted {
			note.IsEncrypted = false
			note.Password = ""
			note.Salt = ""
		}
	}

	if req.TagIDs != nil {
		var tags []models.Tag
		database.DB.Where("id IN ?", req.TagIDs).Find(&tags)
		database.DB.Model(&note).Association("Tags").Replace(tags)
	}

	if err := database.DB.Save(&note).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	database.DB.Preload("Tags").First(&note, id)
	c.JSON(http.StatusOK, note)
}

func DeleteNote(c *gin.Context) {
	id := c.Param("id")
	permanent := c.DefaultQuery("permanent", "false")

	var note models.Note
	if err := database.DB.First(&note, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Note not found"})
		return
	}

	if permanent == "true" {
		database.DB.Model(&note).Association("Tags").Clear()
		database.DB.Delete(&note)
	} else {
		note.IsDeleted = true
		note.DeletedAt = time.Now()
		database.DB.Save(&note)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Note deleted"})
}

func RestoreNote(c *gin.Context) {
	id := c.Param("id")
	var note models.Note
	if err := database.DB.First(&note, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Note not found"})
		return
	}

	note.IsDeleted = false
	note.DeletedAt = time.Time{}
	database.DB.Save(&note)

	c.JSON(http.StatusOK, note)
}

func TogglePin(c *gin.Context) {
	id := c.Param("id")
	var note models.Note
	if err := database.DB.First(&note, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Note not found"})
		return
	}

	note.IsPinned = !note.IsPinned
	database.DB.Save(&note)

	c.JSON(http.StatusOK, note)
}

func ToggleArchive(c *gin.Context) {
	id := c.Param("id")
	var note models.Note
	if err := database.DB.First(&note, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Note not found"})
		return
	}

	note.IsArchived = !note.IsArchived
	database.DB.Save(&note)

	c.JSON(http.StatusOK, note)
}

func EmptyTrash(c *gin.Context) {
	var notes []models.Note
	database.DB.Where("is_deleted = ?", true).Find(&notes)
	for _, note := range notes {
		database.DB.Model(&note).Association("Tags").Clear()
		database.DB.Delete(&note)
	}
	c.JSON(http.StatusOK, gin.H{"message": "Trash emptied"})
}

func GetStats(c *gin.Context) {
	var total int64
	var archived int64
	var deleted int64
	var pinned int64

	database.DB.Model(&models.Note{}).Where("is_deleted = ? AND is_archived = ?", false, false).Count(&total)
	database.DB.Model(&models.Note{}).Where("is_archived = ? AND is_deleted = ?", true, false).Count(&archived)
	database.DB.Model(&models.Note{}).Where("is_deleted = ?", true).Count(&deleted)
	database.DB.Model(&models.Note{}).Where("is_pinned = ? AND is_deleted = ? AND is_archived = ?", true, false, false).Count(&pinned)

	var tagsCount int64
	database.DB.Model(&models.Tag{}).Count(&tagsCount)

	c.JSON(http.StatusOK, gin.H{
		"total":    total,
		"archived": archived,
		"deleted":  deleted,
		"pinned":   pinned,
		"tags":     tagsCount,
	})
}
