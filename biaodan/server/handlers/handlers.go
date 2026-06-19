package handlers

import (
	"biaodan/database"
	"biaodan/models"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type TemplateSchema struct {
	Fields []map[string]interface{} `json:"fields"`
}

type TemplateRequest struct {
	Name        string                 `json:"name" binding:"required"`
	Description string                 `json:"description"`
	Schema      map[string]interface{} `json:"schema" binding:"required"`
}

func ListTemplates(c *gin.Context) {
	var templates []models.FormTemplate
	keyword := c.Query("keyword")

	query := database.DB.Order("created_at DESC")
	if keyword != "" {
		query = query.Where("name LIKE ?", "%"+keyword+"%")
	}

	if err := query.Find(&templates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	result := make([]gin.H, 0)
	for _, t := range templates {
		var schema map[string]interface{}
		_ = json.Unmarshal([]byte(t.Schema), &schema)

		var count int64
		database.DB.Model(&models.FormSubmission{}).Where("template_id = ?", t.ID).Count(&count)

		result = append(result, gin.H{
			"id":          t.ID,
			"name":        t.Name,
			"description": t.Description,
			"schema":      schema,
			"status":      t.Status,
			"submissions": count,
			"createdAt":   t.CreatedAt,
			"updatedAt":   t.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{"data": result})
}

func GetTemplate(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var template models.FormTemplate

	if err := database.DB.First(&template, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Template not found"})
		return
	}

	var schema map[string]interface{}
	_ = json.Unmarshal([]byte(template.Schema), &schema)

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"id":          template.ID,
			"name":        template.Name,
			"description": template.Description,
			"schema":      schema,
			"status":      template.Status,
			"createdAt":   template.CreatedAt,
			"updatedAt":   template.UpdatedAt,
		},
	})
}

func CreateTemplate(c *gin.Context) {
	var req TemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	schemaBytes, _ := json.Marshal(req.Schema)

	template := models.FormTemplate{
		Name:        req.Name,
		Description: req.Description,
		Schema:      string(schemaBytes),
		Status:      1,
	}

	if err := database.DB.Create(&template).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var schema map[string]interface{}
	_ = json.Unmarshal([]byte(template.Schema), &schema)

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"id":          template.ID,
			"name":        template.Name,
			"description": template.Description,
			"schema":      schema,
			"status":      template.Status,
			"createdAt":   template.CreatedAt,
			"updatedAt":   template.UpdatedAt,
		},
	})
}

func UpdateTemplate(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var template models.FormTemplate

	if err := database.DB.First(&template, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Template not found"})
		return
	}

	var req TemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	schemaBytes, _ := json.Marshal(req.Schema)

	template.Name = req.Name
	template.Description = req.Description
	template.Schema = string(schemaBytes)

	if err := database.DB.Save(&template).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var schema map[string]interface{}
	_ = json.Unmarshal([]byte(template.Schema), &schema)

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"id":          template.ID,
			"name":        template.Name,
			"description": template.Description,
			"schema":      schema,
			"status":      template.Status,
			"createdAt":   template.CreatedAt,
			"updatedAt":   template.UpdatedAt,
		},
	})
}

func DeleteTemplate(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	if err := database.DB.Delete(&models.FormTemplate{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	database.DB.Where("template_id = ?", id).Delete(&models.FormSubmission{})

	c.JSON(http.StatusOK, gin.H{"message": "Deleted successfully"})
}

type SubmissionRequest struct {
	TemplateID uint                   `json:"templateId" binding:"required"`
	Data       map[string]interface{} `json:"data" binding:"required"`
}

func validateSubmission(template *models.FormTemplate, data map[string]interface{}) []string {
	var errors []string

	var schema TemplateSchema
	_ = json.Unmarshal([]byte(template.Schema), &schema)

	for _, field := range schema.Fields {
		label, _ := field["label"].(string)
		fieldId, _ := field["id"].(string)
		required, _ := field["required"].(bool)
		fieldType, _ := field["type"].(string)

		value, exists := data[fieldId]

		if required {
			if !exists || value == nil {
				errors = append(errors, label+"不能为空")
				continue
			}

			switch v := value.(type) {
			case string:
				if v == "" {
					errors = append(errors, label+"不能为空")
				}
			case []interface{}:
				if len(v) == 0 {
					errors = append(errors, label+"不能为空")
				}
			}
		}

		if fieldType == "email" && value != nil {
			if v, ok := value.(string); ok && v != "" {
				if !isValidEmail(v) {
					errors = append(errors, label+"格式不正确")
				}
			}
		}

		if fieldType == "phone" && value != nil {
			if v, ok := value.(string); ok && v != "" {
				if !isValidPhone(v) {
					errors = append(errors, label+"格式不正确")
				}
			}
		}
	}

	return errors
}

func isValidEmail(email string) bool {
	return len(email) > 3 && contains(email, "@") && contains(email, ".")
}

func isValidPhone(phone string) bool {
	if len(phone) != 11 {
		return false
	}
	for _, c := range phone {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func CreateSubmission(c *gin.Context) {
	var req SubmissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var template models.FormTemplate
	if err := database.DB.First(&template, req.TemplateID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Template not found"})
		return
	}

	if validationErrors := validateSubmission(&template, req.Data); len(validationErrors) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed", "details": validationErrors})
		return
	}

	dataBytes, _ := json.Marshal(req.Data)

	submission := models.FormSubmission{
		TemplateID: req.TemplateID,
		Data:       string(dataBytes),
		IP:         c.ClientIP(),
		UserAgent:  c.GetHeader("User-Agent"),
	}

	if err := database.DB.Create(&submission).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"id":        submission.ID,
			"createdAt": submission.CreatedAt,
		},
		"message": "提交成功",
	})
}

func ListSubmissions(c *gin.Context) {
	templateId, _ := strconv.Atoi(c.Query("templateId"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	keyword := c.Query("keyword")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var template models.FormTemplate
	if err := database.DB.First(&template, templateId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Template not found"})
		return
	}

	var schema TemplateSchema
	_ = json.Unmarshal([]byte(template.Schema), &schema)

	query := database.DB.Model(&models.FormSubmission{}).Where("template_id = ?", templateId)
	var total int64
	query.Count(&total)

	var submissions []models.FormSubmission
	if err := query.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&submissions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	records := make([]gin.H, 0)
	for _, s := range submissions {
		var data map[string]interface{}
		_ = json.Unmarshal([]byte(s.Data), &data)

		if keyword != "" {
			found := false
			for _, v := range data {
				if str, ok := v.(string); ok && contains(str, keyword) {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		records = append(records, gin.H{
			"id":        s.ID,
			"data":      data,
			"ip":        s.IP,
			"createdAt": s.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"records": records,
			"total":   total,
			"page":    page,
			"size":    pageSize,
			"fields":  schema.Fields,
		},
	})
}

func DeleteSubmission(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	if err := database.DB.Delete(&models.FormSubmission{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Deleted successfully"})
}

func ExportSubmissions(c *gin.Context) {
	templateId, _ := strconv.Atoi(c.Query("templateId"))

	var template models.FormTemplate
	if err := database.DB.First(&template, templateId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Template not found"})
		return
	}

	var schema TemplateSchema
	_ = json.Unmarshal([]byte(template.Schema), &schema)

	var submissions []models.FormSubmission
	database.DB.Where("template_id = ?", templateId).Order("created_at DESC").Find(&submissions)

	csvContent := ""
	for i, field := range schema.Fields {
		if i > 0 {
			csvContent += ","
		}
		label, _ := field["label"].(string)
		csvContent += "\"" + escapeCSV(label) + "\""
	}
	csvContent += ",\"提交时间\"\n"

	for _, s := range submissions {
		var data map[string]interface{}
		_ = json.Unmarshal([]byte(s.Data), &data)

		for i, field := range schema.Fields {
			if i > 0 {
				csvContent += ","
			}
			fieldId, _ := field["id"].(string)
			value := data[fieldId]
			strValue := formatValue(value)
			csvContent += "\"" + escapeCSV(strValue) + "\""
		}
		csvContent += ",\"" + s.CreatedAt.Format("2006-01-02 15:04:05") + "\"\n"
	}

	filename := template.Name + "_" + time.Now().Format("20060102150405") + ".csv"

	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", "attachment; filename=\""+filename+"\"")
	c.String(http.StatusOK, "\uFEFF"+csvContent)
}

func escapeCSV(s string) string {
	result := ""
	for _, c := range s {
		if c == '"' {
			result += "\"\""
		} else {
			result += string(c)
		}
	}
	return result
}

func formatValue(v interface{}) string {
	if v == nil {
		return ""
	}
	switch val := v.(type) {
	case string:
		return val
	case []interface{}:
		result := ""
		for i, item := range val {
			if i > 0 {
				result += "; "
			}
			result += formatValue(item)
		}
		return result
	default:
		b, _ := json.Marshal(v)
		return string(b)
	}
}

func GetTemplateStats(c *gin.Context) {
	templateId, _ := strconv.Atoi(c.Param("templateId"))

	var total int64
	database.DB.Model(&models.FormSubmission{}).Where("template_id = ?", templateId).Count(&total)

	var todayCount int64
	today := time.Now().Format("2006-01-02")
	database.DB.Model(&models.FormSubmission{}).
		Where("template_id = ? AND DATE(created_at) = ?", templateId, today).
		Count(&todayCount)

	var last7Days []gin.H
	for i := 6; i >= 0; i-- {
		date := time.Now().AddDate(0, 0, -i).Format("2006-01-02")
		var count int64
		database.DB.Model(&models.FormSubmission{}).
			Where("template_id = ? AND DATE(created_at) = ?", templateId, date).
			Count(&count)
		last7Days = append(last7Days, gin.H{
			"date":  date,
			"count": count,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"total":      total,
			"todayCount": todayCount,
			"last7Days":  last7Days,
		},
	})
}
