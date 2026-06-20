package models

import (
	"time"

	"gorm.io/gorm"
)

type LogisticsOrder struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	OrderNo     string         `gorm:"size:50;uniqueIndex;not null" json:"order_no"`
	Sender      string         `gorm:"size:100" json:"sender"`
	Receiver    string         `gorm:"size:100" json:"receiver"`
	Origin      string         `gorm:"size:200" json:"origin"`
	Destination string         `gorm:"size:200" json:"destination"`
	GoodsName   string         `gorm:"size:200" json:"goods_name"`
	Weight      float64        `json:"weight"`
	Status      string         `gorm:"size:50;default:'pending'" json:"status"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Nodes       []LogisticsNode `gorm:"foreignKey:OrderID" json:"nodes,omitempty"`
}

type LogisticsNode struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	OrderID        uint           `gorm:"index;not null" json:"order_id"`
	NodeName       string         `gorm:"size:100;not null" json:"node_name"`
	Status         string         `gorm:"size:50;default:'in_transit'" json:"status"`
	Location       string         `gorm:"size:200" json:"location"`
	Description    string         `gorm:"size:500" json:"description"`
	Operator       string         `gorm:"size:100" json:"operator"`
	IsAbnormal     bool           `gorm:"default:false" json:"is_abnormal"`
	AbnormalReason string         `gorm:"size:500" json:"abnormal_reason"`
	OccurredAt     time.Time      `json:"occurred_at"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

const (
	OrderStatusPending   = "pending"
	OrderStatusPicked    = "picked"
	OrderStatusInTransit = "in_transit"
	OrderStatusArrived   = "arrived"
	OrderStatusDelivered = "delivered"
	OrderStatusException = "exception"
)

const (
	NodeStatusPending   = "pending"
	NodeStatusPicked    = "picked"
	NodeStatusInTransit = "in_transit"
	NodeStatusArrived   = "arrived"
	NodeStatusDelivered = "delivered"
	NodeStatusException = "exception"
)
