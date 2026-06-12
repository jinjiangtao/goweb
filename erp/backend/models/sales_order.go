package models

import (
	"time"
	"gorm.io/gorm"
)

type SalesOrderStatus string

const (
	SalesOrderStatusDraft       SalesOrderStatus = "draft"
	SalesOrderStatusApproved    SalesOrderStatus = "approved"
	SalesOrderStatusPartialShip SalesOrderStatus = "partial_ship"
	SalesOrderStatusCompleted   SalesOrderStatus = "completed"
	SalesOrderStatusCancelled   SalesOrderStatus = "cancelled"
)

type SalesOrder struct {
	gorm.Model
	OrderNo       string               `gorm:"uniqueIndex;not null" json:"orderNo"`
	CustomerID    uint                 `gorm:"not null" json:"customerId"`
	Customer      Customer             `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
	Items         []SalesOrderItem     `gorm:"foreignKey:SalesOrderID" json:"items,omitempty"`
	TotalAmount   float64              `gorm:"type:decimal(10,2);not null" json:"totalAmount"`
	Status        SalesOrderStatus     `gorm:"not null;default:'draft'" json:"status"`
	DeliveryDate  *time.Time           `json:"deliveryDate"`
	Remark        string               `json:"remark"`
	CreatedByID   uint                 `gorm:"not null" json:"createdById"`
	CreatedBy     User                 `gorm:"foreignKey:CreatedByID" json:"createdBy,omitempty"`
}

type SalesOrderItem struct {
	gorm.Model
	SalesOrderID uint    `gorm:"not null" json:"salesOrderId"`
	ProductID    uint    `gorm:"not null" json:"productId"`
	Product      Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	Quantity     float64 `gorm:"type:decimal(10,2);not null" json:"quantity"`
	UnitPrice    float64 `gorm:"type:decimal(10,2);not null" json:"unitPrice"`
	Amount       float64 `gorm:"type:decimal(10,2);not null" json:"amount"`
}
