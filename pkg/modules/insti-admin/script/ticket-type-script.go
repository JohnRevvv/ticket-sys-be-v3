package script

import (
	"ideyanale-be/pkg/config"
	IAdmodel "ideyanale-be/pkg/modules/insti-admin/model"
)

// ==========================
// POST SCRIPTS
// ==========================
func AddTicketType(tickettypename string, institutionID int) error {
	return config.DBConnList[0].Exec(`
		INSERT INTO ticket_types (
			ticket_type_name,
			institution_id
		)
		VALUES (?, ?)
	`,
		tickettypename,
		institutionID,
	).Error
}

func AddCategory(categoryname string, tickettypeID int) error {
	return config.DBConnList[0].Exec(`
		INSERT INTO categories (
			category_name,
			ticket_type_id
		)
		VALUES (?, ?)
	`,
		categoryname,
		tickettypeID,
	).Error
}

func AddSubCategory(subCategoryName, subjectName, description string, categoryID int) error {
	return config.DBConnList[0].Exec(`
		INSERT INTO sub_categories (
			subcategory_name,
			subject_name,
			description,
			category_id
		)
		VALUES (?, ?, ?, ?)
	`,
		subCategoryName,
		subjectName,
		description,
		categoryID,
	).Error
}

//==========================
// GET SCRIPTS
//==========================

func GetTicketTypesByInstitutionID(institutionID int) ([]IAdmodel.TicketType, error) {
	var alltickettypes []IAdmodel.TicketType

	err := config.DBConnList[0].Raw(`
		SELECT
			ticket_type_id,
			ticket_type_name,
			institution_id
		FROM ticket_types
		WHERE institution_id = ?
	`, institutionID).Scan(&alltickettypes).Error

	return alltickettypes, err
}

func GetTicketTypeByID(ticketTypeID int) (*IAdmodel.TicketType, error) {
	var tickettype IAdmodel.TicketType

	err := config.DBConnList[0].Raw(`
		SELECT
			ticket_type_id,
			ticket_type_name,
			institution_id,
			status
		FROM ticket_types
		WHERE ticket_type_id = ?
	`, ticketTypeID).Scan(&tickettype).Error

	if err != nil {
		return nil, err
	}

	if tickettype.TicketTypeID == 0 {
		return nil, nil
	}

	return &tickettype, nil
}

func GetCategoriesByTicketTypeID(ticketTypeID int) ([]IAdmodel.Category, error) {
	var allcategories []IAdmodel.Category

	err := config.DBConnList[0].Raw(`
		SELECT
			category_id,
			ticket_type_id,
			category_name,
			status
		FROM categories
		WHERE ticket_type_id = ?
	`, ticketTypeID).Scan(&allcategories).Error

	return allcategories, err
}

func GetCategoryByID(categoryID int) (*IAdmodel.Category, error) {
	var category IAdmodel.Category

	err := config.DBConnList[0].Raw(`
		SELECT
			category_id,
			ticket_type_id,
			category_name,
			status
		FROM categories
		WHERE category_id = ?
	`, categoryID).Scan(&category).Error

	if err != nil {
		return nil, err
	}

	if category.CategoryID == 0 {
		return nil, nil
	}

	return &category, nil
}

func GetSubCategoriesByCategoryID(categoryID int) ([]IAdmodel.SubCategory, error) {
	var allsubcategories []IAdmodel.SubCategory

	err := config.DBConnList[0].Raw(`
		SELECT
			sub_category_id,
			category_id,
			subject_name,
			subcategory_name,
			description,
			status
		FROM sub_categories
		WHERE category_id = ?
	`, categoryID).Scan(&allsubcategories).Error

	return allsubcategories, err
}

func GetSubCategoryByID(subCategoryID int) (*IAdmodel.SubCategory, error) {
	var subcategory IAdmodel.SubCategory

	err := config.DBConnList[0].Raw(`
		SELECT
			sub_category_id,
			category_id,
			subject_name,
			subcategory_name,
			description,
			status
		FROM sub_categories
		WHERE sub_category_id = ?
	`, subCategoryID).Scan(&subcategory).Error

	if err != nil {
		return nil, err
	}

	if subcategory.SubCategoryID == 0 {
		return nil, nil
	}

	return &subcategory, nil
}

//==========================
// EDIT SCRIPTS
//==========================

func EditSubCategory(subCategoryID int, subCategoryName, subjectName, description, status string) error {
	return config.DBConnList[0].Exec(`
		UPDATE sub_categories
		SET
			subcategory_name = ?,
			subject_name = ?,
			description = ?,
			status = ?
		WHERE sub_category_id = ?
	`,
		subCategoryName,
		subjectName,
		description,
		status,
		subCategoryID,
	).Error
}