package model

import (
	"time"

	IAdmodel "ideyanale-be/pkg/modules/insti-admin/model"
)

type Institution struct {
	InstitutionID   uint   `gorm:"primaryKey" json:"institution_id"`
	InstitutionCode	string `gorm:"column:institution_code" json:"institution_code"`
	InstitutionName string `gorm:"column:institution_name;not null" json:"institution_name"`
	Description     string `gorm:"column:description" json:"description"`
	Status          string `gorm:"default:'active';not null" json:"status"`

	JobPosition []IAdmodel.JobPosition `gorm:"foreignKey:InstitutionID"`

	CreatedAt time.Time `json:"created_at"`
}

func (Institution) TableName() string {
	return "institutions"
}