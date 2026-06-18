package models

import (
	"time"

	"gorm.io/gorm"
)

type ApprovalRequest struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	RequestNo       string         `gorm:"size:50;uniqueIndex;not null" json:"request_no"`
	Title           string         `gorm:"size:200;not null" json:"title"`
	Type            string         `gorm:"size:50" json:"type"`
	TemplateID      uint           `json:"template_id"`
	Template        ProcessTemplate `gorm:"foreignKey:TemplateID" json:"template,omitempty"`
	FormData        string         `gorm:"type:text" json:"form_data"`
	CurrentNodeID   string         `gorm:"size:50" json:"current_node_id"`
	Status          string         `gorm:"size:20;default:'pending'" json:"status"`
	ApplicantID     uint           `json:"applicant_id"`
	Applicant       User           `gorm:"foreignKey:ApplicantID" json:"applicant,omitempty"`
	ProcessInstance string         `gorm:"type:text" json:"process_instance"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

type ApprovalRecord struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	RequestID       uint           `json:"request_id"`
	Request         ApprovalRequest `gorm:"foreignKey:RequestID" json:"request,omitempty"`
	NodeID          string         `gorm:"size:50" json:"node_id"`
	NodeName        string         `gorm:"size:100" json:"node_name"`
	ApproverID      uint           `json:"approver_id"`
	Approver        User           `gorm:"foreignKey:ApproverID" json:"approver,omitempty"`
	Action          string         `gorm:"size:20" json:"action"`
	Comment         string         `gorm:"size:500" json:"comment"`
	CCUsers         string         `gorm:"type:text" json:"cc_users"`
	SignType        string         `gorm:"size:20;default:'normal'" json:"sign_type"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
}

type ProcessInstance struct {
	Nodes         []NodeInstance `json:"nodes"`
	Edges         []Edge         `json:"edges"`
	CurrentNodeID string         `json:"current_node_id"`
	History       []HistoryRecord `json:"history"`
}

type NodeInstance struct {
	NodeID   string   `json:"node_id"`
	Status   string   `json:"status"`
	Approved []uint   `json:"approved"`
	Rejected []uint   `json:"rejected"`
	CCUsers  []uint   `json:"cc_users"`
}

type HistoryRecord struct {
	NodeID     string    `json:"node_id"`
	NodeName   string    `json:"node_name"`
	Action     string    `json:"action"`
	ApproverID uint      `json:"approver_id"`
	Approver   string    `json:"approver"`
	Comment    string    `json:"comment"`
	Timestamp  time.Time `json:"timestamp"`
}

const (
	ApprovalStatusPending  = "pending"
	ApprovalStatusApproved = "approved"
	ApprovalStatusRejected = "rejected"
	ApprovalStatusRevoked  = "revoked"

	NodeStatusPending    = "pending"
	NodeStatusProcessing = "processing"
	NodeStatusApproved   = "approved"
	NodeStatusRejected   = "rejected"
	NodeStatusSkipped    = "skipped"

	ApprovalActionApprove   = "approve"
	ApprovalActionReject    = "reject"
	ApprovalActionRevoke    = "revoke"
	ApprovalActionAddSign   = "add_sign"
	ApprovalActionTransfer  = "transfer"
	ApprovalActionAutoPass  = "auto_pass"
	ApprovalActionStart     = "start"
	ApprovalActionEnd       = "end"
)
