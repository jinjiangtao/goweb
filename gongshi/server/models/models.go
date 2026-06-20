package models

import (
	"time"
)

type User struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Username string `gorm:"size:50;uniqueIndex" json:"username"`
	Name     string `gorm:"size:50" json:"name"`
	Role     string `gorm:"size:20;default:employee" json:"role"`
	Password string `gorm:"size:100" json:"-"`
	CreatedAt time.Time `json:"created_at"`
}

type Project struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"size:100" json:"name"`
	Code        string    `gorm:"size:50;uniqueIndex" json:"code"`
	Description string    `gorm:"size:500" json:"description"`
	Status      string    `gorm:"size:20;default:active" json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}

type WorkRecord struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserID      uint      `gorm:"index" json:"user_id"`
	ProjectID   uint      `gorm:"index" json:"project_id"`
	WorkDate    string    `gorm:"size:10;index" json:"work_date"`
	Hours       float64   `gorm:"type:decimal(4,1)" json:"hours"`
	Content     string    `gorm:"size:500" json:"content"`
	Status      string    `gorm:"size:20;default:pending" json:"status"`
	ReviewerID  *uint     `gorm:"index" json:"reviewer_id,omitempty"`
	RejectReason string  `gorm:"size:200" json:"reject_reason,omitempty"`
	ReviewedAt  *time.Time `json:"reviewed_at,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	User      User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Project   Project   `gorm:"foreignKey:ProjectID" json:"project,omitempty"`
	Reviewer  *User     `gorm:"foreignKey:ReviewerID" json:"reviewer,omitempty"`
}
