package model

import (
	"time"

	"gorm.io/gorm"
)

type Roles struct {
	RoleID           uint           `gorm:"column:role_id;primaryKey" json:"role_id"`
	InstitutionID    uint           `gorm:"column:institution_id;not null" json:"institution_id"`
	RoleName         string         `gorm:"column:role_name;not null" json:"role_name"`
	CanCreateTicket  bool           `gorm:"column:can_create;default:false" json:"can_create"`
	CanEndorseTicket bool           `gorm:"column:can_endorse;default:false" json:"can_endorse"`
	CanApproveTicket bool           `gorm:"column:can_approve;default:false" json:"can_approve"`
	CanResolveTicket bool           `gorm:"column:can_resolve;default:false" json:"can_resolve"`
	CanAudit         bool           `gorm:"column:can_audit;default:false" json:"can_audit"`
	CreatedAt        time.Time      `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`
}
