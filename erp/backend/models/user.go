package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"uniqueIndex;not null" json:"username"`
	Password string `gorm:"not null" json:"-"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Avatar   string `json:"avatar"`
	Status   int    `gorm:"default:1" json:"status"`
	RoleID   uint   `json:"roleId"`
	Role     Role   `gorm:"foreignKey:RoleID" json:"role,omitempty"`
}
