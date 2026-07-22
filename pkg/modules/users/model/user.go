package model

import (
	Positionmodel "ideyanale-be/pkg/modules/positions/models"
	Rolemodel "ideyanale-be/pkg/modules/roles/model"
	"time"

	"gorm.io/gorm"
)

// UserDetails represents a staff/user record tied to an institution and role.
type UserDetails struct {
	ID            int            `gorm:"primaryKey;autoIncrement" json:"id,omitempty"`
	Username      string         `json:"username,omitempty"`
	StaffID       string         `json:"staff_id,omitempty"`
	FirstName     string         `json:"first_name,omitempty"`
	LastName      string         `json:"last_name,omitempty"`
	Email         string         `json:"email,omitempty"`
	PhoneNo       string         `json:"phone_no,omitempty"`
	InstitutionID uint           `json:"institution_id,omitempty"`
	PositionID    uint           `json:"position_id,omitempty"`
	RoleID        uint           `gorm:"column:role_id;not null" json:"role_id"`
	Status        string         `json:"status"`
	LastLogin     string         `json:"last_login,omitempty"`
	IsLoggedIn    bool           `gorm:"default:false" json:"is_logged_in,omitempty"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`

	// Role        Rolemodel.Roles           `gorm:"-" json:"role"`
	// JobPosition Positionmodel.JobPosition `gorm:"-" json:"job_position"`
	Role        Rolemodel.Roles           `gorm:"foreignKey:RoleID;references:RoleID;constraint:false"`
	JobPosition Positionmodel.JobPosition `gorm:"foreignKey:PositionID;references:PositionID;constraint:false"`
}

func (UserDetails) TableName() string {
	return "users"
}

// LoginOTP stores a one-time-password used during staff login.
type LoginOTP struct {
	ID        int    `json:"id,omitempty"`
	StaffID   string `json:"staff_id,omitempty"`
	OTPHash   string `json:"otp_hash,omitempty"`
	ExpiresAt int64  `json:"expires_at,omitempty"`
}

func (LoginOTP) TableName() string {
	return "login_otps"
}
