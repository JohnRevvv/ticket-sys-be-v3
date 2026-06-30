package script

import (
	"ideyanale-be/pkg/config"
	encrypDecryptV1 "ideyanale-be/pkg/middleware/encryption/v1"
	IAdmodel "ideyanale-be/pkg/modules/insti-admin/model"

	"gorm.io/gorm"
)

// ==========================
// AUTO-GENERATE SCRIPTS
// ==========================
var DefaultTicketTypes = []string{
	"Service Request",
	// "General Inquiry",
	// "Technical Support",
	// "Billing",
	// "Complaint",
	// "Feature Request",
}

func AddDefaultTicketTypes(institutionID uint) error {
	return config.DBConnList[0].Transaction(func(tx *gorm.DB) error {
		for _, name := range DefaultTicketTypes {
			encName, err := encrypDecryptV1.EncryptV1(name, config.SecretKey)
			if err != nil {
				return err
			}

			if err := tx.Exec(
				`INSERT INTO ticket_types (
					ticket_type_name, 
					institution_id
				)
				 VALUES (?, ?)
				 `,
				encName,
				institutionID,
			).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

type DefaultCategory struct {
	Name string
}

var DefaultCategories = []DefaultCategory{
	{
		Name: "Server",
	},
	{
		Name: "Network",
	},
	{
		Name: "Access",
	},
}

func AddDefaultCategories(ticketTypeID uint) error {
	return config.DBConnList[0].Transaction(func(tx *gorm.DB) error {
		for _, category := range DefaultCategories {
			encName, err := encrypDecryptV1.EncryptV1(category.Name, config.SecretKey,)
			if err != nil {
				return err
			}

			if err := tx.Exec(`
				INSERT INTO categories (
					category_name,
					ticket_type_id
				)
				VALUES (?, ?)
			`,
				encName,
				ticketTypeID,
			).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

type DefaultSubCategory struct {
	SubCategoryName string
	SubjectName     string
	Description     string
}

var DefaultSubCategories = []DefaultSubCategory{
	{
		SubCategoryName: "Server-to-Server Connection",
		SubjectName:     "I Am Requesting for",
		Description:     "Request for server to server connection",
	},
}

func AddDefaultSubCategories(categoryID uint) error {
	return config.DBConnList[0].Transaction(func(tx *gorm.DB) error {
		for _, sub := range DefaultSubCategories {

			encSubCategory, err := encrypDecryptV1.EncryptV1(
				sub.SubCategoryName,
				config.SecretKey,
			)
			if err != nil {
				return err
			}

			encSubject, err := encrypDecryptV1.EncryptV1(
				sub.SubjectName,
				config.SecretKey,
			)
			if err != nil {
				return err
			}

			encDescription, err := encrypDecryptV1.EncryptV1(
				sub.Description,
				config.SecretKey,
			)
			if err != nil {
				return err
			}

			if err := tx.Exec(`
				INSERT INTO sub_categories (
					sub_category_name,
					subject_name,
					description,
					category_id
				)
				VALUES (?, ?, ?, ?)
			`,
				encSubCategory,
				encSubject,
				encDescription,
				categoryID,
			).Error; err != nil {
				return err
			}
		}

		return nil
	})
}


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
			sub_category_name,
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

func GetSubCategoryByID(subCategoryID int) (*IAdmodel.SubCategory, error) {
	var subcategory IAdmodel.SubCategory

	err := config.DBConnList[0].Raw(`
		SELECT
			sub_category_id,
			category_id,
			subject_name,
			sub_category_name,
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

func GetSubCategoriesByCategoryID(categoryID int) ([]IAdmodel.SubCategory, error) {
	var allSubCategories []IAdmodel.SubCategory

	err := config.DBConnList[0].Raw(`
		SELECT
			sub_category_id,
			category_id,
			subject_name,
			sub_category_name,
			description,
			status
		FROM sub_categories
		WHERE category_id = ?
	`, categoryID).Scan(&allSubCategories).Error

	return allSubCategories, err
}

//==========================
// EDIT SCRIPTS
//==========================

func EditTicketType(ticketTypeID int, ticketTypeName, status string) error {
	return config.DBConnList[0].Exec(`
		UPDATE ticket_types
		SET
			ticket_type_name = ?,
			status = ?
		WHERE ticket_type_id = ?
	`,
		ticketTypeName,
		status,
		ticketTypeID,
	).Error
}

func EditCategory(categoryID int, categoryName, status string) error {
	return config.DBConnList[0].Exec(`
		UPDATE categories
		SET
			category_name = ?,
			status = ?
		WHERE category_id = ?
	`,
		categoryName,
		status,
		categoryID,
	).Error
}

func EditSubCategory(subCategoryID int, subCategoryName, subjectName, description, status string) error {
	return config.DBConnList[0].Exec(`
		UPDATE sub_categories
		SET
			sub_category_name = ?,
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
