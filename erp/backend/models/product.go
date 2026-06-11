package models

import (
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Name        string  `gorm:"not null" json:"name"`
	Code        string  `gorm:"uniqueIndex;not null" json:"code"`
	Category    string  `json:"category"`
	Description string  `json:"description"`
	Price       float64 `gorm:"type:decimal(10,2)" json:"price"`
	Stock       int     `gorm:"default:0" json:"stock"`
	Status      int     `gorm:"default:1" json:"status"`
}
