package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"size:50;uniqueIndex;not null" json:"username"`
	Password  string    `gorm:"size:255;not null" json:"-"`
	Nickname  string    `gorm:"size:100" json:"nickname"`
	Avatar    string    `gorm:"size:255" json:"avatar"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Album struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserID      uint      `gorm:"index;not null" json:"user_id"`
	Name        string    `gorm:"size:100;not null" json:"name"`
	Description string    `gorm:"size:500" json:"description"`
	Cover       string    `gorm:"size:255" json:"cover"`
	IsPrivate   bool      `gorm:"default:false" json:"is_private"`
	Password    string    `gorm:"size:100" json:"-"`
	SortOrder   int       `gorm:"default:0" json:"sort_order"`
	PhotoCount  int       `gorm:"default:0" json:"photo_count"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	User        User      `gorm:"foreignKey:UserID" json:"-"`
}

type Media struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	UserID       uint      `gorm:"index;not null" json:"user_id"`
	AlbumID      uint      `gorm:"index;not null" json:"album_id"`
	FileName     string    `gorm:"size:255;not null" json:"file_name"`
	OriginalName string    `gorm:"size:255;not null" json:"original_name"`
	FileSize     int64     `json:"file_size"`
	FileType     string    `gorm:"size:20" json:"file_type"`
	MimeType     string    `gorm:"size:100" json:"mime_type"`
	FilePath     string    `gorm:"size:500;not null" json:"file_path"`
	ThumbPath    string    `gorm:"size:500" json:"thumb_path"`
	Width        int       `json:"width"`
	Height       int       `json:"height"`
	Duration     int       `json:"duration"`
	Latitude     float64   `json:"latitude"`
	Longitude    float64   `json:"longitude"`
	Location     string    `gorm:"size:255" json:"location"`
	SortOrder    int       `gorm:"default:0" json:"sort_order"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Album        Album     `gorm:"foreignKey:AlbumID" json:"-"`
	User         User      `gorm:"foreignKey:UserID" json:"-"`
}

type ShareLink struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	UserID    uint       `gorm:"index;not null" json:"user_id"`
	AlbumID   uint       `gorm:"index;not null" json:"album_id"`
	Token     string     `gorm:"size:100;uniqueIndex;not null" json:"token"`
	Password  string     `gorm:"size:100" json:"-"`
	ViewCount int        `gorm:"default:0" json:"view_count"`
	MaxViews  int        `gorm:"default:0" json:"max_views"`
	ExpiresAt *time.Time `json:"expires_at"`
	IsActive  bool       `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	Album     Album      `gorm:"foreignKey:AlbumID" json:"album,omitempty"`
	User      User       `gorm:"foreignKey:UserID" json:"-"`
}

type VisitRecord struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ShareID   uint      `gorm:"index;not null" json:"share_id"`
	IPAddress string    `gorm:"size:50" json:"ip_address"`
	UserAgent string    `gorm:"size:500" json:"user_agent"`
	VisitedAt time.Time `json:"visited_at"`
	ShareLink ShareLink `gorm:"foreignKey:ShareID" json:"-"`
}

type Tag struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index;not null" json:"user_id"`
	Name      string    `gorm:"size:50;not null" json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type MediaTag struct {
	MediaID uint `gorm:"primaryKey" json:"media_id"`
	TagID   uint `gorm:"primaryKey" json:"tag_id"`
}

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&User{},
		&Album{},
		&Media{},
		&ShareLink{},
		&VisitRecord{},
		&Tag{},
		&MediaTag{},
	)
}
