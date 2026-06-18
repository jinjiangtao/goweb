package models

import (
	"time"
)

type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Username  string    `json:"username" gorm:"uniqueIndex;size:50;not null"`
	Password  string    `json:"-" gorm:"size:255;not null"`
	Email     string    `json:"email" gorm:"size:100"`
	Avatar    string    `json:"avatar" gorm:"size:255"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Snippet struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	Title         string    `json:"title" gorm:"size:200;not null"`
	Description   string    `json:"description" gorm:"size:500"`
	Content       string    `json:"content" gorm:"type:text;not null"`
	Language      string    `json:"language" gorm:"size:50;not null;index"`
	Visibility    string    `json:"visibility" gorm:"size:20;not null;default:public;index"`
	UserID        uint      `json:"user_id" gorm:"index;not null"`
	User          User      `json:"user" gorm:"foreignKey:UserID"`
	Tags          string    `json:"tags" gorm:"size:500"`
	LikesCount    int       `json:"likes_count" gorm:"default:0"`
	FavsCount     int       `json:"favs_count" gorm:"default:0"`
	ViewsCount    int       `json:"views_count" gorm:"default:0"`
	CommentsCount int       `json:"comments_count" gorm:"default:0"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type Like struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"index;not null"`
	SnippetID uint      `json:"snippet_id" gorm:"index;not null"`
	CreatedAt time.Time `json:"created_at"`
}

type Favorite struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"index;not null"`
	SnippetID uint      `json:"snippet_id" gorm:"index;not null"`
	CreatedAt time.Time `json:"created_at"`
}

type Comment struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Content   string    `json:"content" gorm:"type:text;not null"`
	UserID    uint      `json:"user_id" gorm:"index;not null"`
	User      User      `json:"user" gorm:"foreignKey:UserID"`
	SnippetID uint      `json:"snippet_id" gorm:"index;not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6,max=50"`
	Email    string `json:"email" binding:"omitempty,email"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type SnippetRequest struct {
	Title       string `json:"title" binding:"required,max=200"`
	Description string `json:"description" binding:"max=500"`
	Content     string `json:"content" binding:"required"`
	Language    string `json:"language" binding:"required,max=50"`
	Visibility  string `json:"visibility" binding:"oneof=public private"`
	Tags        string `json:"tags"`
}

type CommentRequest struct {
	Content string `json:"content" binding:"required"`
}

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type PaginatedResponse struct {
	Code     int         `json:"code"`
	Message  string      `json:"message"`
	Data     interface{} `json:"data"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
}
