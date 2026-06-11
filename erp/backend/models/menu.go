package models

import (
	"gorm.io/gorm"
)

type Menu struct {
	gorm.Model
	ParentID  *uint   `json:"parentId"`
	Name      string  `gorm:"not null" json:"name"`
	Path      string  `json:"path"`
	Component string  `json:"component"`
	Icon      string  `json:"icon"`
	Sort      int     `gorm:"default:0" json:"sort"`
	Type      string  `gorm:"default:'menu'" json:"type"`
	Permission string `json:"permission"`
	Status    int     `gorm:"default:1" json:"status"`
	Children  []Menu  `gorm:"foreignKey:ParentID" json:"children,omitempty"`
}
