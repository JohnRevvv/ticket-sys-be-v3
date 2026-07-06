package model

type (
	ServerModel struct {
		ServerID      uint   `gorm:"primaryKey" json:"server_id"`
		InstitutionID uint   `gorm:"not null" json:"institution_id"`
		ServerName    string `gorm:"column:server_name;not null" json:"server_name"`
		ServerIP      string `gorm:"column:server_ip;not null" json:"server_ip"`
		Status        string `gorm:"default:'active';not null" json:"status"`

		Project []ProjectModel `gorm:"foreignKey:ServerID"`
	}

	ProjectModel struct {
		ProjectID	 uint   `gorm:"primaryKey" json:"project_id"`
		ServerID     uint   `gorm:"not null" json:"server_id"`
		ProjectName  string `gorm:"column:project_name;not null" json:"project_name"`
		Environment  string `gorm:"column:environment;not null" json:"environment"`
		Description  string `gorm:"column:description" json:"description"`
		Status       string `gorm:"default:'active';not null" json:"status"`
	}
)

func (ServerModel) TableName() string {
	return "servers"
}

func (ProjectModel) TableName() string {
	return "projects"
}