package model

type SuperAdminDetails struct {
	ID       int    `gorm:"column:id" json:"id"`
	UserName string `gorm:"column:username" json:"username"`
	Role     string `gorm:"column:role" json:"role"`
	Email    string `gorm:"column:email" json:"email"`
	Password string `gorm:"column:password" json:"password"`
}

func (SuperAdminDetails) TableName() string {
	return "admin"
}