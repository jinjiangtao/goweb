package models

import (
	"time"

	"gorm.io/gorm"
)

var DB *gorm.DB

type Admin struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Username  string         `gorm:"uniqueIndex;size:50;not null" json:"username"`
	Password  string         `gorm:"size:255;not null" json:"-"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type House struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	Title         string         `gorm:"size:200;not null" json:"title"`
	Community     string         `gorm:"size:100;not null" json:"community"`
	District      string         `gorm:"size:50;not null" json:"district"`
	Address       string         `gorm:"size:300;not null" json:"address"`
	RoomType      string         `gorm:"size:50;not null" json:"room_type"`
	Area          float64        `gorm:"not null" json:"area"`
	Price         float64        `gorm:"not null" json:"price"`
	Floor         string         `gorm:"size:50" json:"floor"`
	Orientation   string         `gorm:"size:50" json:"orientation"`
	Decoration    string         `gorm:"size:50" json:"decoration"`
	Facilities    string         `gorm:"size:500" json:"facilities"`
	Description   string         `gorm:"type:text;not null" json:"description"`
	ContactName   string         `gorm:"size:50;not null" json:"contact_name"`
	ContactPhone  string         `gorm:"size:20;not null" json:"contact_phone"`
	Images        string         `gorm:"type:text" json:"images"`
	Status        int            `gorm:"default:1" json:"status"`
	IsRecommended bool           `gorm:"default:false" json:"is_recommended"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

func (House) TableName() string {
	return "houses"
}
