package handlers

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"shuiti/database"
	"shuiti/models"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type GenerateExamRequest struct {
	Mode       string   `json:"mode"`
	Count      int      `json:"count"`
	Types      []string `json:"types"`
	Categories []string `json:"categories"`
	Difficulty int      `json:"difficulty"`
}

func GenerateExam(c *gin.Context) {
	var req GenerateExamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Count <= 0 {
		req.Count = 10
	}

	query := database.DB.Model(&models.Question{})

	if len(req.Types) > 0 {
		query = query.Where("type IN ?", req.Types)
	}
	if len(req.Categories) > 0 {
		query = query.Where("category IN ?", req.Categories)
	}
	if req.Difficulty > 0 {
		query = query.Where("difficulty = ?", req.Difficulty)
	}

	var allQuestions []models.Question
	err := query.Find(&allQuestions).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(allQuestions) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"questions": []models.Question{},
			"count":     0,
		})
		return
	}

	var selectedQuestions []models.Question

	switch req.Mode {
	case "random":
		rand.Seed(time.Now().UnixNano())
		shuffleQuestions(allQuestions)
		if req.Count > len(allQuestions) {
			req.Count = len(allQuestions)
		}
		selectedQuestions = allQuestions[:req.Count]
	case "smart":
		typeQuestions := make(map[string][]models.Question)
		for _, q := range allQuestions {
			typeQuestions[string(q.Type)] = append(typeQuestions[string(q.Type)], q)
		}
		var smartResult []models.Question
		perType := req.Count / 3
		if perType == 0 {
			perType = 1
		}
		for _, qs := range typeQuestions {
			shuffleQuestions(qs)
			limit := perType
			if limit > len(qs) {
				limit = len(qs)
			}
			smartResult = append(smartResult, qs[:limit]...)
		}
		shuffleQuestions(smartResult)
		if req.Count > len(smartResult) {
			req.Count = len(smartResult)
		}
		if len(smartResult) > req.Count {
			smartResult = smartResult[:req.Count]
		}
		selectedQuestions = smartResult
	default:
		sort.Slice(allQuestions, func(i, j int) bool {
			return allQuestions[i].ID < allQuestions[j].ID
		})
		if req.Count > len(allQuestions) {
			req.Count = len(allQuestions)
		}
		selectedQuestions = allQuestions[:req.Count]
	}

	c.JSON(http.StatusOK, gin.H{
		"questions": selectedQuestions,
		"count":     len(selectedQuestions),
		"mode":      req.Mode,
	})
}

func shuffleQuestions(qs []models.Question) {
	rand.Shuffle(len(qs), func(i, j int) {
		qs[i], qs[j] = qs[j], qs[i]
	})
}

func SubmitExam(c *gin.Context) {
	var req models.ExamSubmitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	questionIDs := make([]uint, 0, len(req.Answers))
	answerMap := make(map[uint]string)
	for _, a := range req.Answers {
		questionIDs = append(questionIDs, a.QuestionID)
		answerMap[a.QuestionID] = a.Answer
	}

	var questions []models.Question
	database.DB.Where("id IN ?", questionIDs).Find(&questions)

	qMap := make(map[uint]models.Question)
	for _, q := range questions {
		qMap[q.ID] = q
	}

	correctCount := 0
	var details []models.QuestionResultDetail
	var wrongQuestionsToSave []models.WrongQuestion

	for _, qid := range questionIDs {
		q, ok := qMap[qid]
		if !ok {
			continue
		}
		userAnswer := answerMap[qid]
		isCorrect := checkAnswer(q.Type, q.Answer, userAnswer)

		if isCorrect {
			correctCount++
		} else {
			var wq models.WrongQuestion
			database.DB.Where("user_id = ? AND question_id = ?", req.UserID, qid).First(&wq)
			if wq.ID == 0 {
				wrongQuestionsToSave = append(wrongQuestionsToSave, models.WrongQuestion{
					UserID:      req.UserID,
					QuestionID:  qid,
					UserAnswer:  userAnswer,
					WrongCount:  1,
					LastWrongAt: time.Now(),
				})
			} else {
				database.DB.Model(&wq).Updates(map[string]interface{}{
					"wrong_count":   wq.WrongCount + 1,
					"user_answer":   userAnswer,
					"last_wrong_at": time.Now(),
				})
			}
		}

		details = append(details, models.QuestionResultDetail{
			QuestionID: qid,
			Question:   q,
			UserAnswer: userAnswer,
			IsCorrect:  isCorrect,
		})
	}

	if len(wrongQuestionsToSave) > 0 {
		database.DB.Create(&wrongQuestionsToSave)
	}

	totalCount := len(questionIDs)
	score := 0.0
	if totalCount > 0 {
		score = float64(correctCount) / float64(totalCount) * 100
	}

	answersJSON, _ := json.Marshal(answerMap)
	qidsJSON, _ := json.Marshal(questionIDs)

	examName := req.ExamName
	if examName == "" {
		examName = "练习" + time.Now().Format("2006-01-02 15:04")
	}

	record := models.ExamRecord{
		UserID:       req.UserID,
		ExamName:     examName,
		QuestionIDs:  string(qidsJSON),
		Answers:      string(answersJSON),
		TotalCount:   totalCount,
		CorrectCount: correctCount,
		Score:        score,
		Mode:         req.Mode,
	}
	database.DB.Create(&record)

	result := models.ExamResult{
		TotalCount:   totalCount,
		CorrectCount: correctCount,
		Score:        score,
		Details:      details,
	}

	c.JSON(http.StatusOK, gin.H{
		"success":   true,
		"record_id": record.ID,
		"result":    result,
	})
}

func checkAnswer(qType models.QuestionType, correctAnswer, userAnswer string) bool {
	if userAnswer == "" {
		return false
	}

	switch qType {
	case models.MultipleChoice:
		correctSet := buildSet(strings.ToUpper(correctAnswer))
		userSet := buildSet(strings.ToUpper(userAnswer))
		if len(correctSet) != len(userSet) {
			return false
		}
		for k := range correctSet {
			if !userSet[k] {
				return false
			}
		}
		return true
	default:
		return strings.EqualFold(strings.TrimSpace(correctAnswer), strings.TrimSpace(userAnswer))
	}
}

func buildSet(s string) map[string]bool {
	set := make(map[string]bool)
	for _, c := range s {
		ch := string(c)
		if ch != "," && ch != " " && ch != "，" {
			set[ch] = true
		}
	}
	return set
}

func GetExamRecords(c *gin.Context) {
	userID := c.Query("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	query := database.DB.Model(&models.ExamRecord{})
	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}

	var total int64
	query.Count(&total)

	var records []models.ExamRecord
	offset := (page - 1) * pageSize
	query.Order("id DESC").Offset(offset).Limit(pageSize).Find(&records)

	c.JSON(http.StatusOK, gin.H{
		"total":     total,
		"list":      records,
		"page":      page,
		"page_size": pageSize,
	})
}

func GetExamRecord(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var record models.ExamRecord
	err = database.DB.First(&record, uint(id)).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
		return
	}

	var questionIDs []uint
	json.Unmarshal([]byte(record.QuestionIDs), &questionIDs)

	var answers map[uint]string
	json.Unmarshal([]byte(record.Answers), &answers)

	var questions []models.Question
	database.DB.Where("id IN ?", questionIDs).Find(&questions)
	qMap := make(map[uint]models.Question)
	for _, q := range questions {
		qMap[q.ID] = q
	}

	var details []models.QuestionResultDetail
	for _, qid := range questionIDs {
		q, ok := qMap[qid]
		if !ok {
			continue
		}
		userAns := answers[qid]
		details = append(details, models.QuestionResultDetail{
			QuestionID: qid,
			Question:   q,
			UserAnswer: userAns,
			IsCorrect:  checkAnswer(q.Type, q.Answer, userAns),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"record":  record,
		"details": details,
	})
}
