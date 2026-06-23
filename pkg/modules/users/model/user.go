package model

import "time"

type (
	UserDetails struct {
		ID              int       `gorm:"primaryKey;autoIncrement" json:"id,omitempty"`
		Username        string    `json:"username,omitempty"`
		StaffID         string    `json:"staff_id,omitempty"`
		FirstName       string    `json:"first_name,omitempty"`
		LastName        string    `json:"last_name,omitempty"`
		Email           string    `json:"email,omitempty"`
		PhoneNo         string    `json:"phone_no,omitempty"`
		InstitutionID   int       `json:"institution_id,omitempty"`
		InstitutionName string    `json:"institution_name,omitempty"`
		JobPosition     string    `json:"job_position,omitempty"`
		Role            string    `json:"role"`
		Status          string    `json:"status"`
		LastLogin       string    `json:"last_login,omitempty"`
		IsLoggedIn      bool      `json:"is_logged_in,omitempty"`
		CreatedAt       time.Time `json:"created_at"`
	}

	LoginOTP struct {
		ID        int    `json:"id,omitempty"`
		StaffID   string `json:"staff_id,omitempty"`
		OTPHash   string `json:"otp_hash,omitempty"`
		ExpiresAt int64  `json:"expires_at,omitempty"`
	}
)

func (UserDetails) TableName() string {
	return "users"
}

func (LoginOTP) TableName() string {
	return "login_otps"
}
