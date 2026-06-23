package model

import (
	"time"

	IAdmodel "ideyanale-be/pkg/modules/insti-admin/model"
)

type Institution struct {
	InstitutionID   uint   `gorm:"primaryKey" json:"institution_id"`
	InstitutionName string `gorm:"unique;not null" json:"institution_name"`
	Description     string `json:"description"`
	Status          string `gorm:"default:'active';not null" json:"status"`

	JobPosition []IAdmodel.JobPosition `gorm:"foreignKey:InstitutionID" json:"jobposition"`

	CreatedAt time.Time `json:"created_at"`
}

func (Institution) TableName() string {
	return "institutions"
}