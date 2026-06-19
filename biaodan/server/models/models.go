package models

import (
	"time"

	"gorm.io/gorm"
)

type FormTemplate struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	Name        string         `gorm:"size:255;not null" json:"name"`
	Description string         `gorm:"type:text" json:"description"`
	Schema      string         `gorm:"type:text;not null" json:"schema"`
	Status      int            `gorm:"default:1" json:"status"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

type FormSubmission struct {
	ID         uint           `gorm:"primarykey" json:"id"`
	TemplateID uint           `gorm:"index;not null" json:"templateId"`
	Data       string         `gorm:"type:text;not null" json:"data"`
	IP         string         `gorm:"size:100" json:"ip"`
	UserAgent  string         `gorm:"size:500" json:"userAgent"`
	CreatedAt  time.Time      `json:"createdAt"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}
