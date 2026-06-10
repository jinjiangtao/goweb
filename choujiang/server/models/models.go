package models

import (
	"sync"
	"time"

	"gorm.io/gorm"
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
	ID             uint      `gorm:"primaryKey" json:"id"`
	Name           string    `json:"name"`
	Phone          string    `json:"phone"`
	PrizeID        uint      `json:"prizeId"`
	PrizeName      string    `json:"prizeName"`
	IsWin          bool      `json:"isWin"`
	Status         string    `json:"status"`
	ReceiverName   string    `json:"receiverName"`
	ReceiverPhone  string    `json:"receiverPhone"`
	Province       string    `json:"province"`
	City           string    `json:"city"`
	District       string    `json:"district"`
	DetailAddress  string    `json:"detailAddress"`
	ShippingStatus string    `json:"shippingStatus"`
	TrackingNumber string    `json:"trackingNumber"`
	CreatedAt      time.Time `json:"createdAt"`
}

type Address struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	ParentID uint   `gorm:"default:0" json:"parentId"`
	Name     string `json:"name"`
	Level    int    `json:"level"` // 1=省, 2=市, 3=区
	Code     string `json:"code"`
}
