package model

import (
	"time"

)

type Institution struct {
	InstitutionID   uint   `gorm:"primaryKey" json:"institution_id"`
	InstitutionName string `gorm:"unique;not null" json:"institution_name"`
	Description     string `json:"description"`
	Status          string `gorm:"default:'active';not null" json:"status"`


	CreatedAt time.Time `json:"created_at"`
}

func (Institution) TableName() string {
	return "institutions"
}
