package models

import (
	"time"
)

type QuestionType string

const (
	SingleChoice   QuestionType = "single"
	MultipleChoice QuestionType = "multiple"
	TrueFalse      QuestionType = "truefalse"
)

type Question struct {
	ID         uint         `gorm:"primaryKey" json:"id"`
	Type       QuestionType `gorm:"size:20;not null" json:"type"`
	Content    string       `gorm:"type:text;not null" json:"content"`
	Options    string       `gorm:"type:text" json:"options"`
	Answer     string       `gorm:"size:200;not null" json:"answer"`
	Analysis   string       `gorm:"type:text" json:"analysis"`
	Difficulty int          `gorm:"default:1" json:"difficulty"`
	Category   string       `gorm:"size:100" json:"category"`
	CreatedAt  time.Time    `json:"created_at"`
	UpdatedAt  time.Time    `json:"updated_at"`
}

type ExamRecord struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	UserID       string    `gorm:"size:100;not null;index" json:"user_id"`
	ExamName     string    `gorm:"size:200" json:"exam_name"`
	QuestionIDs  string    `gorm:"type:text;not null" json:"question_ids"`
	Answers      string    `gorm:"type:text;not null" json:"answers"`
	TotalCount   int       `gorm:"not null" json:"total_count"`
	CorrectCount int       `gorm:"not null" json:"correct_count"`
	Score        float64   `gorm:"not null" json:"score"`
	Mode         string    `gorm:"size:50" json:"mode"`
	CreatedAt    time.Time `json:"created_at"`
}

type WrongQuestion struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserID      string    `gorm:"size:100;not null;index" json:"user_id"`
	QuestionID  uint      `gorm:"not null;index" json:"question_id"`
	UserAnswer  string    `gorm:"size:200;not null" json:"user_answer"`
	WrongCount  int       `gorm:"default:1" json:"wrong_count"`
	LastWrongAt time.Time `json:"last_wrong_at"`
	Question    Question  `gorm:"foreignKey:QuestionID" json:"question,omitempty"`
}

type AnswerSubmission struct {
	QuestionID uint   `json:"question_id" binding:"required"`
	Answer     string `json:"answer"`
}

type ExamSubmitRequest struct {
	UserID   string             `json:"user_id" binding:"required"`
	ExamName string             `json:"exam_name"`
	Mode     string             `json:"mode"`
	Answers  []AnswerSubmission `json:"answers" binding:"required"`
}

type ExamResult struct {
	TotalCount   int                    `json:"total_count"`
	CorrectCount int                    `json:"correct_count"`
	Score        float64                `json:"score"`
	Details      []QuestionResultDetail `json:"details"`
}

type QuestionResultDetail struct {
	QuestionID uint     `json:"question_id"`
	Question   Question `json:"question"`
	UserAnswer string   `json:"user_answer"`
	IsCorrect  bool     `json:"is_correct"`
}
