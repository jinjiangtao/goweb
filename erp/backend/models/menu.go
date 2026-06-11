package models

import (
	"gorm.io/gorm"
)

type Menu struct {
	gorm.Model
	ParentID *uint  `json:"parentId"`
	Name     string `gorm:"not null" json:"name"`
	Path     string `json:"path"`
	Icon     string `json:"icon"`
	Sort     int    `gorm:"default:0" json:"sort"`
	Hidden   bool   `gorm:"default:false" json:"hidden"`
	Children []Menu `gorm:"foreignKey:ParentID" json:"children,omitempty"`
}
