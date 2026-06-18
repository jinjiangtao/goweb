package models

import (
	"time"
)

type Category struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"size:50;not null"`
	Color     string    `json:"color" gorm:"size:20;not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Event struct {
	ID            uint       `json:"id" gorm:"primaryKey"`
	Title         string     `json:"title" gorm:"size:200;not null"`
	Description   string     `json:"description" gorm:"type:text"`
	StartTime     time.Time  `json:"start_time" gorm:"not null"`
	EndTime       time.Time  `json:"end_time" gorm:"not null"`
	CategoryID    *uint      `json:"category_id"`
	Category      *Category  `json:"category" gorm:"foreignKey:CategoryID"`
	IsRecurring   bool       `json:"is_recurring" gorm:"default:false"`
	RecurringType string     `json:"recurring_type" gorm:"size:20"`
	RecurringEnd  *time.Time `json:"recurring_end"`
	HasReminder   bool       `json:"has_reminder" gorm:"default:false"`
	ReminderTime  *time.Time `json:"reminder_time"`
	Location      string     `json:"location" gorm:"size:500"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}
