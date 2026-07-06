package models

import (
	// UserModel "ideyanale-be/pkg/modules/users/model"
	"time"
)

type Ticket struct {
	ID              uint       `gorm:"primaryKey" json:"id"`
	TicketID        string     `gorm:"uniqueIndex;not null" json:"ticket_id"`
	ProjectID       uint       `json:"project_id"`
	InstitutionID   uint       `json:"institution_id"`
	TicketTypeID    uint       `json:"ticket_type_id"`
	CategoryID      uint       `json:"category_id"`
	SubCategoryID   uint       `json:"subcategory_id"`
	Subject         string     `json:"subject"`
	Description     string     `json:"description"`
	DueDate         *time.Time `json:"due_date"`
	InstitutionPool uint       `json:"institution_pool"`
	SubmitterID     uint       `json:"submitter_id"`
	// Submitter       UserModel.UserDetails `gorm:"foreignKey:SubmitterID;references:ID"`
	ResolverID uint `json:"resolver_id"`
	EndorserID uint `json:"endorser_id"`
	// Endorser        UserModel.UserDetails `gorm:"foreignKey:EndorserID;references:ID"`
	ApproverID uint `json:"approver_id"`

	TicketAttachment []TicketAttachment `gorm:"foreignKey:TicketID;references:TicketID"`
	Status           string             `json:"status" gorm:"default:'for endorsement'"`
	CreatedAt        time.Time          `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt        time.Time          `json:"updated_at" gorm:"autoUpdateTime"`

	CancelledBy        string     `json:"cancelled_by"`
	CancelledAt        *time.Time `json:"cancelled_at"`
	CancellationReason string     `json:"cancellation_reason"`

	StartedAt      *time.Time `json:"started_at"`
	ResolvedAt     *time.Time `json:"resolved_at"`
	ResolutionTime string     `json:"resolution_time" gorm:"column:resolution_time;default:''"`
	OnHold         bool       `json:"onhold" gorm:"column:on_hold;default:false"`
	HoldAt         *time.Time `json:"hold_at"`
	ClosedBy       string     `json:"closed_by"`
	ClosedAt       *time.Time `json:"closed_at"`

	EndorsedAt *time.Time `json:"endorsed_at"`
	ApprovedAt *time.Time `json:"approved_at"`

	CloseToken     string `json:"-" gorm:"column:close_token"`
	CloseTokenUsed bool   `json:"-" gorm:"column:close_token_used"`
}

func (Ticket) TableName() string {
	return "tickets"
}

// ── TicketRemark ──────────────────────────────────────────────────────────────

type TicketRemark struct {
	RemarkID  string    `gorm:"primaryKey" json:"remark_id"`
	TicketID  string    `json:"ticket_id"`
	UserID    string    `json:"user_id"`
	Username  string    `json:"username"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
}

func (TicketRemark) TableName() string {
	return "ticketremark"
}
