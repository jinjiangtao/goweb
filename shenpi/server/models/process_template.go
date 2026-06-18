package models

import (
	"time"

	"gorm.io/gorm"
)

type ProcessTemplate struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"size:100;not null" json:"name"`
	Description string         `gorm:"size:500" json:"description"`
	Type        string         `gorm:"size:50;default:'general'" json:"type"`
	FormConfig  string         `gorm:"type:text" json:"form_config"`
	Nodes       string         `gorm:"type:text" json:"nodes"`
	Edges       string         `gorm:"type:text" json:"edges"`
	Status      string         `gorm:"size:20;default:'draft'" json:"status"`
	CreatorID   uint           `json:"creator_id"`
	Creator     User           `gorm:"foreignKey:CreatorID" json:"creator,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

type Node struct {
	ID            string                 `json:"id"`
	Type          string                 `json:"type"`
	Label         string                 `json:"label"`
	X             int                    `json:"x"`
	Y             int                    `json:"y"`
	Properties    map[string]interface{} `json:"properties"`
	Approvers     []uint                 `json:"approvers"`
	ApproverType  string                 `json:"approver_type"`
	Condition     string                 `json:"condition"`
	ConditionExpr string                 `json:"condition_expr"`
}

type Edge struct {
	ID         string `json:"id"`
	Source     string `json:"source"`
	Target     string `json:"target"`
	Label      string `json:"label"`
	Condition  string `json:"condition"`
}

type FormField struct {
	Name        string                 `json:"name"`
	Label       string                 `json:"label"`
	Type        string                 `json:"type"`
	Required    bool                   `json:"required"`
	Placeholder string                 `json:"placeholder"`
	Options     []map[string]string    `json:"options"`
	Validation  map[string]interface{} `json:"validation"`
	Default     interface{}            `json:"default"`
}
