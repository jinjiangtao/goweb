package models

import (
	"gorm.io/gorm"
)

type Customer struct {
	gorm.Model
	Name    string `gorm:"not null" json:"name"`
	Contact string `json:"contact"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
	Remark  string `json:"remark"`
}
