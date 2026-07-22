package model

type JobPosition struct {
	PositionID    uint   `gorm:"primaryKey" json:"position_id"`
	InstitutionID uint   `gorm:"not null" json:"institution_id"`
	PositionName  string `gorm:"column:position_name;not null" json:"position_name"`
	Status        string `gorm:"default:'active';not null" json:"status"`
}

func (JobPosition) TableName() string {
	return "job_positions"
}
