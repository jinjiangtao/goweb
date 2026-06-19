package handlers

import (
	"encoding/json"
	"net/http"
	"shuiti/database"
	"shuiti/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func GetWrongQuestions(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	query := database.DB.Model(&models.WrongQuestion{}).Where("user_id = ?", userID)

	var total int64
	query.Count(&total)

	var wrongQuestions []models.WrongQuestion
	offset := (page - 1) * pageSize
	query.Preload("Question").Order("last_wrong_at DESC").Offset(offset).Limit(pageSize).Find(&wrongQuestions)

	c.JSON(http.StatusOK, gin.H{
		"total":     total,
		"list":      wrongQuestions,
		"page":      page,
		"page_size": pageSize,
	})
}

func DeleteWrongQuestion(c *gin.Context) {
	userID := c.Query("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	result := database.DB.Where("id = ? AND user_id = ?", uint(id), userID).Delete(&models.WrongQuestion{})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Wrong question not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

func ClearWrongQuestions(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	result := database.DB.Where("user_id = ?", userID).Delete(&models.WrongQuestion{})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"deleted": result.RowsAffected,
	})
}

func GetStatistics(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	var examCount int64
	database.DB.Model(&models.ExamRecord{}).Where("user_id = ?", userID).Count(&examCount)

	var totalQuestions int64
	var totalCorrect int64
	var avgScore float64
	rows := database.DB.Model(&models.ExamRecord{}).Where("user_id = ?", userID).
		Select("COALESCE(SUM(total_count), 0), COALESCE(SUM(correct_count), 0), COALESCE(AVG(score), 0)").
		Row()
	rows.Scan(&totalQuestions, &totalCorrect, &avgScore)

	var wrongCount int64
	database.DB.Model(&models.WrongQuestion{}).Where("user_id = ?", userID).Count(&wrongCount)

	var records []models.ExamRecord
	database.DB.Where("user_id = ?", userID).Order("id DESC").Limit(30).Find(&records)

	type DailyStat struct {
		Date      string  `json:"date"`
		ExamCount int     `json:"exam_count"`
		Questions int     `json:"questions"`
		Correct   int     `json:"correct"`
		Accuracy  float64 `json:"accuracy"`
	}

	dailyMap := make(map[string]*DailyStat)
	for _, r := range records {
		date := r.CreatedAt.Format("2006-01-02")
		if _, ok := dailyMap[date]; !ok {
			dailyMap[date] = &DailyStat{Date: date}
		}
		dailyMap[date].ExamCount++
		dailyMap[date].Questions += r.TotalCount
		dailyMap[date].Correct += r.CorrectCount
	}

	var dailyStats []DailyStat
	for _, v := range dailyMap {
		if v.Questions > 0 {
			v.Accuracy = float64(v.Correct) / float64(v.Questions) * 100
		}
		dailyStats = append(dailyStats, *v)
	}

	for i := 0; i < len(dailyStats)-1; i++ {
		for j := i + 1; j < len(dailyStats); j++ {
			if dailyStats[i].Date > dailyStats[j].Date {
				dailyStats[i], dailyStats[j] = dailyStats[j], dailyStats[i]
			}
		}
	}

	var typeStats []struct {
		Type     string  `json:"type"`
		Total    int64   `json:"total"`
		Correct  int64   `json:"correct"`
		Accuracy float64 `json:"accuracy"`
	}

	var allRecords []models.ExamRecord
	database.DB.Where("user_id = ?", userID).Find(&allRecords)

	typeDetailMap := make(map[models.QuestionType]struct{ Total, Correct int64 })

	for _, r := range allRecords {
		var qids []uint
		var answers map[uint]string
		importJSON(r.QuestionIDs, &qids)
		importJSON(r.Answers, &answers)

		var questions []models.Question
		database.DB.Where("id IN ?", qids).Find(&questions)
		qMap := make(map[uint]models.Question)
		for _, q := range questions {
			qMap[q.ID] = q
		}

		for _, qid := range qids {
			q, ok := qMap[qid]
			if !ok {
				continue
			}
			userAns := answers[qid]
			stat := typeDetailMap[q.Type]
			stat.Total++
			if checkAnswer(q.Type, q.Answer, userAns) {
				stat.Correct++
			}
			typeDetailMap[q.Type] = stat
		}
	}

	for t, s := range typeDetailMap {
		acc := float64(0)
		if s.Total > 0 {
			acc = float64(s.Correct) / float64(s.Total) * 100
		}
		typeStats = append(typeStats, struct {
			Type     string  `json:"type"`
			Total    int64   `json:"total"`
			Correct  int64   `json:"correct"`
			Accuracy float64 `json:"accuracy"`
		}{
			Type:     string(t),
			Total:    s.Total,
			Correct:  s.Correct,
			Accuracy: acc,
		})
	}

	overallAccuracy := float64(0)
	if totalQuestions > 0 {
		overallAccuracy = float64(totalCorrect) / float64(totalQuestions) * 100
	}

	c.JSON(http.StatusOK, gin.H{
		"overview": gin.H{
			"exam_count":       examCount,
			"total_questions":  totalQuestions,
			"total_correct":    totalCorrect,
			"overall_accuracy": overallAccuracy,
			"avg_score":        avgScore,
			"wrong_count":      wrongCount,
		},
		"daily_stats":  dailyStats,
		"type_stats":   typeStats,
		"recent_exams": records,
	})
}

func importJSON(s string, v interface{}) {
	json.Unmarshal([]byte(s), v)
}

func GetWrongQuestionsForPractice(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	count, _ := strconv.Atoi(c.DefaultQuery("count", "0"))

	var wrongQuestions []models.WrongQuestion
	query := database.DB.Where("user_id = ?", userID).Preload("Question").Order("wrong_count DESC, last_wrong_at DESC")

	if count > 0 {
		query.Limit(count)
	}
	query.Find(&wrongQuestions)

	var questions []models.Question
	for _, wq := range wrongQuestions {
		questions = append(questions, wq.Question)
	}

	c.JSON(http.StatusOK, gin.H{
		"questions": questions,
		"count":     len(questions),
		"mode":      "wrong",
	})
}

func RemoveWrongAfterCorrect(c *gin.Context) {
	var req struct {
		UserID     string `json:"user_id" binding:"required"`
		QuestionID uint   `json:"question_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	database.DB.Where("user_id = ? AND question_id = ?", req.UserID, req.QuestionID).Delete(&models.WrongQuestion{})

	c.JSON(http.StatusOK, gin.H{"success": true})
}

func GetDashboard(c *gin.Context) {
	var totalQuestions int64
	database.DB.Model(&models.Question{}).Count(&totalQuestions)

	var typeCount []struct {
		Type  string `json:"type"`
		Count int64  `json:"count"`
	}
	database.DB.Model(&models.Question{}).Select("type, count(*) as count").Group("type").Scan(&typeCount)

	var totalExams int64
	database.DB.Model(&models.ExamRecord{}).Count(&totalExams)

	var recentExams []models.ExamRecord
	database.DB.Order("id DESC").Limit(10).Find(&recentExams)

	var todayStart time.Time
	todayStart = time.Now().Truncate(24 * time.Hour)
	var todayExams int64
	database.DB.Model(&models.ExamRecord{}).Where("created_at >= ?", todayStart).Count(&todayExams)

	c.JSON(http.StatusOK, gin.H{
		"total_questions": totalQuestions,
		"total_exams":     totalExams,
		"today_exams":     todayExams,
		"type_count":      typeCount,
		"recent_exams":    recentExams,
	})
}
