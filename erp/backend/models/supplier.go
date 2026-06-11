package models

import (
	"gorm.io/gorm"
)

type Supplier struct {
	gorm.Model
	Name    string `gorm:"not null" json:"name"`
	Contact string `json:"contact"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
	Remark  string `json:"remark"`
}
