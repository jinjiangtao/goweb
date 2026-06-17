package models

import (
	"time"
)

type ArticleVersion struct {
	ID              uint      `json:"id" gorm:"primaryKey"`
	ArticleID       uint      `json:"article_id" gorm:"not null;index"`
	VersionNumber   int       `json:"version_number" gorm:"not null"`
	TitleSnapshot   string    `json:"title_snapshot" gorm:"size:255;not null"`
	ContentSnapshot string    `json:"content_snapshot" gorm:"type:text;not null"`
	CategoryIDSnap  uint      `json:"category_id_snap" gorm:"not null"`
	StatusSnapshot  string    `json:"status_snapshot" gorm:"size:20;not null"`
	CreatorID       uint      `json:"creator_id" gorm:"not null"`
	IsRollback      bool      `json:"is_rollback" gorm:"default:false"`
	RollbackFromVer *int      `json:"rollback_from_ver"`
	CreatedAt       time.Time `json:"created_at"`
}

func (ArticleVersion) TableName() string {
	return "article_versions"
}
