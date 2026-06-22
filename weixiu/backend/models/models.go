package models

import (
	"time"
)

type UserRole string

const (
	RoleAdmin     UserRole = "admin"
	RoleUser      UserRole = "user"
	RoleTechnician UserRole = "technician"
)

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"uniqueIndex;size:50;not null" json:"username"`
	Password  string    `gorm:"size:255;not null" json:"-"`
	RealName  string    `gorm:"size:50" json:"real_name"`
	Role      UserRole  `gorm:"size:20;not null;default:user" json:"role"`
	Phone     string    `gorm:"size:20" json:"phone"`
	Email     string    `gorm:"size:100" json:"email"`
	Area      string    `gorm:"size:100" json:"area"`
	Skills    string    `gorm:"size:255" json:"skills"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type DeviceType string

const (
	DeviceComputer  DeviceType = "computer"
	DevicePrinter   DeviceType = "printer"
	DeviceNetwork   DeviceType = "network"
	DeviceProjector DeviceType = "projector"
	DeviceAircon    DeviceType = "aircon"
	DeviceOther     DeviceType = "other"
)

type Device struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	DeviceCode   string     `gorm:"uniqueIndex;size:50;not null" json:"device_code"`
	DeviceName   string     `gorm:"size:100;not null" json:"device_name"`
	DeviceType   DeviceType `gorm:"size:30;not null" json:"device_type"`
	Location     string     `gorm:"size:200;not null" json:"location"`
	Area         string     `gorm:"size:100" json:"area"`
	Status       string     `gorm:"size:20;default:normal" json:"status"`
	BuyDate      *time.Time `json:"buy_date"`
	Description  string     `gorm:"size:500" json:"description"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

type OrderStatus string

const (
	StatusPending    OrderStatus = "pending"
	StatusAssigned   OrderStatus = "assigned"
	StatusProcessing OrderStatus = "processing"
	StatusCompleted  OrderStatus = "completed"
	StatusRejected   OrderStatus = "rejected"
	StatusCancelled  OrderStatus = "cancelled"
)

type OrderPriority string

const (
	PriorityLow    OrderPriority = "low"
	PriorityMedium OrderPriority = "medium"
	PriorityHigh   OrderPriority = "high"
	PriorityUrgent OrderPriority = "urgent"
)

type WorkOrder struct {
	ID             uint          `gorm:"primaryKey" json:"id"`
	OrderNo        string        `gorm:"uniqueIndex;size:30;not null" json:"order_no"`
	Title          string        `gorm:"size:200;not null" json:"title"`
	Description    string        `gorm:"size:1000;not null" json:"description"`
	FaultImages    string        `gorm:"size:1000" json:"fault_images"`
	DeviceType     DeviceType    `gorm:"size:30;not null" json:"device_type"`
	DeviceName     string        `gorm:"size:100" json:"device_name"`
	DeviceID       *uint         `json:"device_id"`
	Area           string        `gorm:"size:100;not null" json:"area"`
	Location       string        `gorm:"size:200" json:"location"`
	ContactName    string        `gorm:"size:50" json:"contact_name"`
	ContactPhone   string        `gorm:"size:20" json:"contact_phone"`
	Priority       OrderPriority `gorm:"size:20;default:medium" json:"priority"`
	Status         OrderStatus   `gorm:"size:20;default:pending" json:"status"`
	UserID         uint          `gorm:"not null" json:"user_id"`
	User           *User         `gorm:"foreignKey:UserID" json:"user,omitempty"`
	TechnicianID   *uint         `json:"technician_id"`
	Technician     *User         `gorm:"foreignKey:TechnicianID" json:"technician,omitempty"`
	ExpectTime     *time.Time    `json:"expect_time"`
	AssignTime     *time.Time    `json:"assign_time"`
	StartProcessTime *time.Time  `json:"start_process_time"`
	CompleteTime   *time.Time    `json:"complete_time"`
	RejectReason   string        `gorm:"size:500" json:"reject_reason"`
	IsOverdue      bool          `gorm:"default:false" json:"is_overdue"`
	CreatedAt      time.Time     `json:"created_at"`
	UpdatedAt      time.Time     `json:"updated_at"`
}

type RepairLog struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	OrderID      uint      `gorm:"not null;index" json:"order_id"`
	WorkOrder    *WorkOrder `gorm:"foreignKey:OrderID" json:"-"`
	TechnicianID uint      `gorm:"not null" json:"technician_id"`
	Technician   *User     `gorm:"foreignKey:TechnicianID" json:"technician,omitempty"`
	Content      string    `gorm:"size:2000;not null" json:"content"`
	Status       OrderStatus `gorm:"size:20" json:"status"`
	Images       string    `gorm:"size:1000" json:"images"`
	CreatedAt    time.Time `json:"created_at"`
}

type OrderEvaluation struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	OrderID    uint      `gorm:"uniqueIndex;not null" json:"order_id"`
	WorkOrder  *WorkOrder `gorm:"foreignKey:OrderID" json:"-"`
	UserID     uint      `gorm:"not null" json:"user_id"`
	Rating     int       `gorm:"not null" json:"rating"`
	Content    string    `gorm:"size:1000" json:"content"`
	CreatedAt  time.Time `json:"created_at"`
}

type OperationLog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `json:"user_id"`
	Username  string    `gorm:"size:50" json:"username"`
	Action    string    `gorm:"size:50;not null" json:"action"`
	Target    string    `gorm:"size:100" json:"target"`
	Detail    string    `gorm:"size:1000" json:"detail"`
	IP        string    `gorm:"size:50" json:"ip"`
	CreatedAt time.Time `json:"created_at"`
}
