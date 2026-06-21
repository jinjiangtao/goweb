package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Board struct {
	ID          string    `gorm:"primaryKey;type:varchar(36)" json:"id"`
	Name        string    `gorm:"type:varchar(100);not null" json:"name"`
	Background  string    `gorm:"type:varchar(50);default:'#ffffff'" json:"background"`
	Thumbnail   string    `gorm:"type:text" json:"thumbnail,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	Elements    []Element `gorm:"foreignKey:BoardID;constraint:OnDelete:CASCADE" json:"elements,omitempty"`
	Histories   []History `gorm:"foreignKey:BoardID;constraint:OnDelete:CASCADE" json:"histories,omitempty"`
}

func (b *Board) BeforeCreate(tx *gorm.DB) error {
	if b.ID == "" {
		b.ID = uuid.New().String()
	}
	return nil
}

type Element struct {
	ID        string    `gorm:"primaryKey;type:varchar(36)" json:"id"`
	BoardID   string    `gorm:"type:varchar(36);index;not null" json:"boardId"`
	Type      string    `gorm:"type:varchar(20);not null" json:"type"`
	UserID    string    `gorm:"type:varchar(36);not null" json:"userId"`
	Points    string    `gorm:"type:text;not null" json:"-"`
	PointsArr []Point   `gorm:"-" json:"points"`
	Style     string    `gorm:"type:text;not null" json:"-"`
	StyleObj  DrawStyle `gorm:"-" json:"style"`
	Text      string    `gorm:"type:text" json:"text,omitempty"`
	FontSize  int       `gorm:"default:14" json:"fontSize,omitempty"`
	FontFamily string   `gorm:"type:varchar(50)" json:"fontFamily,omitempty"`
	ImageData string    `gorm:"type:text" json:"imageData,omitempty"`
	X         float64   `gorm:"default:0" json:"x,omitempty"`
	Y         float64   `gorm:"default:0" json:"y,omitempty"`
	Width     float64   `gorm:"default:0" json:"width,omitempty"`
	Height    float64   `gorm:"default:0" json:"height,omitempty"`
	Timestamp int64     `gorm:"not null" json:"timestamp"`
	CreatedAt time.Time `json:"-"`
}

func (e *Element) BeforeCreate(tx *gorm.DB) error {
	if e.ID == "" {
		e.ID = uuid.New().String()
	}
	return nil
}

type Point struct {
	X        float64 `json:"x"`
	Y        float64 `json:"y"`
	Pressure float64 `json:"pressure,omitempty"`
}

type DrawStyle struct {
	Color       string  `json:"color"`
	Width       float64 `json:"width"`
	Opacity     float64 `json:"opacity"`
	Fill        string  `json:"fill,omitempty"`
	FillOpacity float64 `json:"fillOpacity,omitempty"`
}

type History struct {
	ID         string    `gorm:"primaryKey;type:varchar(36)" json:"id"`
	BoardID    string    `gorm:"type:varchar(36);index;not null" json:"boardId"`
	Action     string    `gorm:"type:varchar(50);not null" json:"action"`
	UserID     string    `gorm:"type:varchar(36);not null" json:"userId"`
	UserName   string    `gorm:"type:varchar(50);not null" json:"userName"`
	Snapshot   string    `gorm:"type:text" json:"snapshot,omitempty"`
	Timestamp  time.Time `json:"timestamp"`
}

func (h *History) BeforeCreate(tx *gorm.DB) error {
	if h.ID == "" {
		h.ID = uuid.New().String()
	}
	return nil
}

type OperationLog struct {
	ID        string    `gorm:"primaryKey;type:varchar(36)"`
	BoardID   string    `gorm:"type:varchar(36);index;not null"`
	UserID    string    `gorm:"type:varchar(36);not null"`
	Action    string    `gorm:"type:varchar(50);not null"`
	Data      string    `gorm:"type:text"`
	Timestamp time.Time `gorm:"index"`
}

func (o *OperationLog) BeforeCreate(tx *gorm.DB) error {
	if o.ID == "" {
		o.ID = uuid.New().String()
	}
	return nil
}
