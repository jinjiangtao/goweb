package models

import (
	"gorm.io/gorm"
)

type Role struct {
	gorm.Model
	Name        string    `gorm:"uniqueIndex;not null" json:"name"`
	Code        string    `gorm:"uniqueIndex;not null" json:"code"`
	Description string    `json:"description"`
	Status      int       `gorm:"default:1" json:"status"`
	Menus       []Menu    `gorm:"many2many:role_menus" json:"menus,omitempty"`
}
