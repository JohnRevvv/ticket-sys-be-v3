package script

import (
	"errors"
	"ideyanale-be/pkg/config"
	"ideyanale-be/pkg/modules/institutions/model"

	"gorm.io/gorm"
)

func GetInstitutionByID(institutionID uint) (*model.Institution, error) {
	var institution model.Institution

	err := config.DBConnList[0].Raw(`
		SELECT
			institution_id,
			institution_name
		FROM institutions
		WHERE institution_id = ?
	`,
		institutionID,
	).Scan(&institution).Error

	if err != nil {
		return nil, err
	}

	if institution.InstitutionID == 0 {
		return nil, nil
	}

	return &institution, nil
}

func GetInstitutions() ([]model.Institution, error) {
	var institutions []model.Institution

	err := config.DBConnList[0].Raw(`
		SELECT
			institution_id,
			institution_code,
			institution_name,
			description,
			created_at
		FROM institutions
	`).Scan(&institutions).Error

	return institutions, err
}

func InstitutionExists(institutionCode, institutionName string) (bool, error) {
	var exists bool

	err := config.DBConnList[0].Raw(`
		SELECT EXISTS(
			SELECT 1
			FROM institutions
			WHERE institution_name = ?
			   OR institution_code = ?
		)
	`,
		institutionName,
		institutionCode,
	).Scan(&exists).Error

	if err != nil {
		return false, err
	}

	return exists, nil
}

func GetInstitutionLogo(institutionID uint) (*model.InstitutionLogo, error) {

	var logo model.InstitutionLogo

	err := config.DBConnList[0].
		Where("institution_id = ?", institutionID).
		First(&logo).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return &logo, err
}