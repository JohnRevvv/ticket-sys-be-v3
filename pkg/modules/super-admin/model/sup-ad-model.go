package model

type SuperAccount struct {
	ID         int    `gorm:"column:id" json:"id"`
	UserName   string `gorm:"column:username" json:"username"`
	Role       string `gorm:"column:role" json:"role"`
	Email      string `gorm:"column:email" json:"email"`
	Password   string `gorm:"column:password" json:"password"`
	IsLoggedIn bool   `gorm:"default:false; column:is_logged_in" json:"is_logged_in,omitempty"`
}

func (SuperAccount) TableName() string {
	return "superaccount"
}
