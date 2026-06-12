package models

import (
	"time"
	"gorm.io/gorm"
)

type PurchaseOrderStatus string

const (
	PurchaseOrderStatusDraft       PurchaseOrderStatus = "draft"
	PurchaseOrderStatusApproved    PurchaseOrderStatus = "approved"
	PurchaseOrderStatusPartialReceive PurchaseOrderStatus = "partial_receive"
	PurchaseOrderStatusCompleted   PurchaseOrderStatus = "completed"
	PurchaseOrderStatusCancelled   PurchaseOrderStatus = "cancelled"
)

type PurchaseOrder struct {
	gorm.Model
	OrderNo       string                `gorm:"uniqueIndex;not null" json:"orderNo"`
	SupplierID    uint                  `gorm:"not null" json:"supplierId"`
	Supplier      Supplier              `gorm:"foreignKey:SupplierID" json:"supplier,omitempty"`
	Items         []PurchaseOrderItem   `gorm:"foreignKey:PurchaseOrderID" json:"items,omitempty"`
	TotalAmount   float64               `gorm:"type:decimal(10,2);not null" json:"totalAmount"`
	Status        PurchaseOrderStatus   `gorm:"not null;default:'draft'" json:"status"`
	ExpectedDate  *time.Time            `json:"expectedDate"`
	Remark        string                `json:"remark"`
	CreatedByID   uint                  `gorm:"not null" json:"createdById"`
	CreatedBy     User                  `gorm:"foreignKey:CreatedByID" json:"createdBy,omitempty"`
}

type PurchaseOrderItem struct {
	gorm.Model
	PurchaseOrderID uint    `gorm:"not null" json:"purchaseOrderId"`
	ProductID       uint    `gorm:"not null" json:"productId"`
	Product         Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	Quantity        float64 `gorm:"type:decimal(10,2);not null" json:"quantity"`
	UnitPrice       float64 `gorm:"type:decimal(10,2);not null" json:"unitPrice"`
	Amount          float64 `gorm:"type:decimal(10,2);not null" json:"amount"`
}
