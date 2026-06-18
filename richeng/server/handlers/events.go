package handlers

import (
	"net/http"
	"richeng-server/database"
	"richeng-server/models"
	"time"

	"github.com/gin-gonic/gin"
)

type EventQueryParams struct {
	StartDate  string `form:"start_date"`
	EndDate    string `form:"end_date"`
	CategoryID uint   `form:"category_id"`
	Keyword    string `form:"keyword"`
}

func expandRecurringEvents(events []models.Event, startDate, endDate time.Time) []models.Event {
	expanded := make([]models.Event, 0)

	for _, event := range events {
		if !event.IsRecurring || event.RecurringType == "" {
			if event.StartTime.Before(endDate) && event.EndTime.After(startDate) {
				expanded = append(expanded, event)
			}
			continue
		}

		duration := event.EndTime.Sub(event.StartTime)
		currentStart := event.StartTime
		recurringEnd := endDate
		if event.RecurringEnd != nil && event.RecurringEnd.Before(endDate) {
			recurringEnd = *event.RecurringEnd
		}

		for currentStart.Before(recurringEnd) {
			if !currentStart.Before(startDate) {
				newEvent := event
				newEvent.StartTime = currentStart
				newEvent.EndTime = currentStart.Add(duration)
				expanded = append(expanded, newEvent)
			}

			switch event.RecurringType {
			case "daily":
				currentStart = currentStart.AddDate(0, 0, 1)
			case "weekly":
				currentStart = currentStart.AddDate(0, 0, 7)
			case "monthly":
				currentStart = currentStart.AddDate(0, 1, 0)
			default:
				goto stopExpand
			}
		}
	stopExpand:
	}

	return expanded
}

func parseISO8601(s string) (time.Time, error) {
	t, err := time.Parse(time.RFC3339Nano, s)
	if err != nil {
		t, err = time.Parse(time.RFC3339, s)
	}
	return t, err
}

func GetEvents(c *gin.Context) {
	var params EventQueryParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var startDate, endDate time.Time
	var startOk, endOk bool

	if params.StartDate != "" {
		if t, err := parseISO8601(params.StartDate); err == nil {
			startDate = t
			startOk = true
		}
	}
	if params.EndDate != "" {
		if t, err := parseISO8601(params.EndDate); err == nil {
			endDate = t
			endOk = true
		}
	}

	if !startOk {
		startDate = time.Now().AddDate(-1, 0, 0)
	}
	if !endOk {
		endDate = time.Now().AddDate(1, 0, 0)
	}

	events := make([]models.Event, 0)
	query := database.DB.Preload("Category")

	if startOk && endOk {
		query = query.Where(
			"(end_time >= ? AND start_time <= ?) OR is_recurring = ?",
			startDate, endDate, true,
		)
	} else if startOk {
		query = query.Where("end_time >= ? OR is_recurring = ?", startDate, true)
	} else if endOk {
		query = query.Where("start_time <= ? OR is_recurring = ?", endDate, true)
	}

	if params.CategoryID > 0 {
		query = query.Where("category_id = ?", params.CategoryID)
	}
	if params.Keyword != "" {
		query = query.Where("title LIKE ? OR description LIKE ?", "%"+params.Keyword+"%", "%"+params.Keyword+"%")
	}

	if err := query.Order("start_time").Find(&events).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	expandedEvents := expandRecurringEvents(events, startDate, endDate)

	c.JSON(http.StatusOK, gin.H{"data": expandedEvents})
}

func GetEvent(c *gin.Context) {
	var event models.Event
	id := c.Param("id")

	if err := database.DB.Preload("Category").First(&event, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": event})
}

func CreateEvent(c *gin.Context) {
	var event models.Event
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if event.CategoryID != nil && *event.CategoryID == 0 {
		event.CategoryID = nil
	}

	if err := database.DB.Create(&event).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	database.DB.Preload("Category").First(&event, event.ID)
	c.JSON(http.StatusCreated, gin.H{"data": event})
}

func UpdateEvent(c *gin.Context) {
	var event models.Event
	id := c.Param("id")

	if err := database.DB.First(&event, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	var input map[string]interface{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if categoryID, ok := input["category_id"].(float64); ok && categoryID == 0 {
		input["category_id"] = nil
	}

	if err := database.DB.Model(&event).Updates(input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	database.DB.Preload("Category").First(&event, id)
	c.JSON(http.StatusOK, gin.H{"data": event})
}

func DeleteEvent(c *gin.Context) {
	var event models.Event
	id := c.Param("id")

	if err := database.DB.First(&event, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	if err := database.DB.Delete(&event).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Event deleted successfully"})
}
