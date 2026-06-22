package handlers

import (
	"net/http"
	"weixiu/config"
	"weixiu/models"

	"github.com/gin-gonic/gin"
)

func GetOperationLogs(c *gin.Context) {
	var logs []models.OperationLog
	query := config.DB.Order("created_at DESC").Limit(200)

	action := c.Query("action")
	if action != "" {
		query = query.Where("action = ?", action)
	}

	query.Find(&logs)
	c.JSON(http.StatusOK, logs)
}
