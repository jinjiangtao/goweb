
package models

import (
	"gorm.io/gorm"
	"sync"
	"time"
)

var DB *gorm.DB
var LotteryMutex sync.Mutex

func InitDB(db *gorm.DB) {
	DB = db
}

type Admin struct {
	ID        uint   `gorm:"primaryKey"`
	Username  string `gorm:"unique"`
	Password  string
	CreatedAt time.Time
}

type Prize struct {
	ID          uint   `gorm:"primaryKey"`
	Name        string
	Probability float64
	Stock       int
	StockUsed   int
	Description string
	ImageURL    string
	Enabled     bool
	CreatedAt   time.Time
}

type Record struct {
	ID         uint      `gorm:"primaryKey"`
	Name       string
	Phone      string
	PrizeID    uint
	PrizeName  string
	IsWin      bool
	Status     string
	CreatedAt  time.Time
}
