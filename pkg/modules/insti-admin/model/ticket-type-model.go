package model

type (
	TicketType struct {
		TicketTypeID   uint   `gorm:"primaryKey" json:"ticket_type_id"`
		InstitutionID  uint   `gorm:"not null" json:"institution_id"`
		TicketTypeName string `gorm:"column:ticket_type_name;not null" json:"ticket_type_name"`
		Status         string `gorm:"default:'active';not null" json:"status"`

		Category []Category `gorm:"foreignKey:CategoryID"`
	}

	Category struct {
		CategoryID   uint   `gorm:"primaryKey" json:"category_id"`
		TicketTypeID uint   `gorm:"not null" json:"ticket_type_id"`
		CategoryName string `gorm:"column:category_name;not null" json:"category_name"`
		Status       string `gorm:"default:'active';not null" json:"status"`

		SubCategory []SubCategory `gorm:"foreignKey:CategoryID"`
	}

	SubCategory struct {
		SubCategoryID   uint   `gorm:"primaryKey" json:"sub_category_id"`
		CategoryID      uint   `gorm:"not null" json:"category_id"`
		SubjectName     string `gorm:"column:subject_name;not null" json:"subject_name"`
		SubCategoryName string `gorm:"column:subcategory_name;not null" json:"sub_category_name"`
		Description     string `gorm:"column:description" json:"description"`
		Status          string `gorm:"default:'active';not null" json:"status"`
	}
)

func (TicketType) TableName() string {
	return "ticket_types"
}

func (Category) TableName() string {
	return "categories"
}

func (SubCategory) TableName() string {
	return "sub_categories"
}
