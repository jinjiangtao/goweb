package models

import (
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Name  string  `gorm:"not null" json:"name"`
	Code  string  `gorm:"uniqueIndex;not null" json:"code"`
	Price float64 `gorm:"type:decimal(10,2)" json:"price"`
	Spec  string  `json:"spec"`
}
