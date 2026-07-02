package model

type (
	Server struct {
		ID            uint   `gorm:"primaryKey" json:"id"`
		ServerName    string `json:"server_name"`
		ServerIP	  string `json:"server_ip"`
		InstitutionID uint   `json:"institution_id"`

		Projects []Project `gorm:"foreignKey:ServerID"`
	}

	Project struct {
		ID            uint   `gorm:"primaryKey" json:"id"`
		ServerID      uint   `json:"server_id"`
		ProjectName   string `json:"project_name"`
		InstitutionID uint   `json:"institution_id"`
	}
)

func (Project) TableName() string {
	return "projects"
}

func (Server) TableName() string {
	return "servers"
}