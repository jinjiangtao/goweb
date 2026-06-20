package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

type JSON map[string]interface{}

func (j JSON) Value() (driver.Value, error) {
	return json.Marshal(j)
}

func (j *JSON) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to unmarshal JSON value")
	}
	return json.Unmarshal(bytes, j)
}

type Template struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"size:100;not null" json:"name"`
	Description string    `gorm:"size:500" json:"description"`
	Category    string    `gorm:"size:50" json:"category"`
	Thumbnail   string    `gorm:"size:500" json:"thumbnail"`
	Content     JSON      `gorm:"type:text" json:"content"`
	StyleConfig JSON      `gorm:"type:text" json:"style_config"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Resume struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	UserIdentity   string    `gorm:"size:100;not null;index" json:"user_identity"`
	TemplateID     uint      `gorm:"not null" json:"template_id"`
	Title          string    `gorm:"size:200;not null" json:"title"`
	CurrentVersion uint      `gorm:"not null;default:1" json:"current_version"`
	Content        JSON      `gorm:"type:text" json:"content"`
	StyleConfig    JSON      `gorm:"type:text" json:"style_config"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type Version struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	ResumeID      uint      `gorm:"not null;index" json:"resume_id"`
	VersionNumber uint      `gorm:"not null;index" json:"version_number"`
	SnapshotName  string    `gorm:"size:200" json:"snapshot_name"`
	Content       JSON      `gorm:"type:text" json:"content"`
	StyleConfig   JSON      `gorm:"type:text" json:"style_config"`
	CreatedAt     time.Time `json:"created_at"`
}
