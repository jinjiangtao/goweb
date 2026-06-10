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
	ID        uint      `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"unique" json:"username"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"createdAt"`
}

type Prize struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `json:"name"`
	Probability float64   `json:"probability"`
	Stock       int       `json:"stock"`
	StockUsed   int       `json:"stockUsed"`
	Description string    `json:"description"`
	ImageURL    string    `json:"imageUrl"`
	Enabled     bool      `json:"enabled"`
	CreatedAt   time.Time `json:"createdAt"`
}

type Record struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `json:"name"`
	Phone     string    `json:"phone"`
	PrizeID   uint      `json:"prizeId"`
	PrizeName string    `json:"prizeName"`
	IsWin     bool      `json:"isWin"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
}