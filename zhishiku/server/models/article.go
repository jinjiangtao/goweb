package models

import (
	"time"

	"gorm.io/gorm"
)

type Article struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Title       string         `json:"title" gorm:"size:255;not null"`
	Content     string         `json:"content" gorm:"type:text;not null"`
	CategoryID  uint           `json:"category_id" gorm:"not null"`
	Category    Category       `json:"category" gorm:"foreignKey:CategoryID"`
	AuthorID    uint           `json:"author_id" gorm:"not null"`
	Author      User           `json:"author" gorm:"foreignKey:AuthorID"`
	Status      string         `json:"status" gorm:"size:20;default:draft"`
	PublishedAt *time.Time     `json:"published_at"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

type ArticleListResponse struct {
	ID           uint       `json:"id"`
	Title        string     `json:"title"`
	CategoryID   uint       `json:"category_id"`
	CategoryName string     `json:"category_name"`
	AuthorID     uint       `json:"author_id"`
	AuthorName   string     `json:"author_name"`
	Status       string     `json:"status"`
	PublishedAt  *time.Time `json:"published_at"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

func (Article) TableName() string {
	return "articles"
}
