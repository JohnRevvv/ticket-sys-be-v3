package model

import (
	"time"

	Positionmodel "ideyanale-be/pkg/modules/positions/models"
	Rolemodel "ideyanale-be/pkg/modules/roles/model"
)

type (
	Institution struct {
		InstitutionID    uint   `gorm:"primaryKey" json:"institution_id"`
		InstitutionCode  string `gorm:"column:institution_code" json:"institution_code"`
		InstitutionName  string `gorm:"column:institution_name;not null" json:"institution_name"`
		Description      string `gorm:"column:description" json:"description"`
		InstitutionColor string `gorm:"column:institution_color" json:"institution_color"`
		Status           string `gorm:"default:'active';not null" json:"status"`

		JobPosition     []Positionmodel.JobPosition `gorm:"-" json:"job_positions"`
		Role            []Rolemodel.Roles           `gorm:"-" json:"roles"`
		InstitutionLogo InstitutionLogo             `gorm:"-" json:"institution_logo"`

		CreatedAt time.Time `json:"created_at"`
	}

	InstitutionLogo struct {
		ID            uint   `gorm:"primaryKey" json:"id"`
		InstitutionID uint   `json:"institution_id"`
		FileName      string `json:"file_name"`
		FileKey       string `json:"file_key"`
		UploadedBy    uint   `json:"uploaded_by"`
	}
)

func (Institution) TableName() string {
	return "institutions"
}

func (InstitutionLogo) TableName() string {
	return "institution_logos"
}
