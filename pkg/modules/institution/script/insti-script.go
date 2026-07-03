package script

import (
	"errors"
	"ideyanale-be/pkg/config"
	"ideyanale-be/pkg/modules/institution/model"
	"time"

	"gorm.io/gorm"
)

func AddInstitution(encinstitutionCode, encinstitutionName, encdescription string) (uint, error) {
	var institutionID uint

	err := config.DBConnList[0].Raw(`
        INSERT INTO institutions (
            institution_code,
            institution_name,
            description,
            created_at
        )
        VALUES (?, ?, ?, ?)
        RETURNING institution_id
    `,
		encinstitutionCode,
		encinstitutionName,
		encdescription,
		time.Now(),
	).Scan(&institutionID).Error

	return institutionID, err
}

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

func UpdateInstitution(institutionID uint, institutionCode, institutionName, description, institutionColor string) error {

	return config.DBConnList[0].Exec(`
		UPDATE institutions
		SET
			institution_code = ?,
			institution_name = ?,
			description = ?,
			institution_color = ?
		WHERE institution_id = ?
	`,
		institutionCode,
		institutionName,
		description,
		institutionColor,
		institutionID,
	).Error
}

func UpsertInstitutionLogo(institutionID uint, fileName, fileKey string, uploadedBy uint) error {

	var count int64

	config.DBConnList[0].
		Model(&model.InstitutionLogo{}).
		Where("institution_id = ?", institutionID).
		Count(&count)

	if count == 0 {

		return config.DBConnList[0].Exec(`
			INSERT INTO institution_logos
			(
				institution_id,
				file_name,
				file_key,
				uploaded_by
			)
			VALUES (?, ?, ?, ?)
		`,
			institutionID,
			fileName,
			fileKey,
			uploadedBy,
		).Error
	}

	return config.DBConnList[0].Exec(`
		UPDATE institution_logos
		SET
			file_name = ?,
			file_key = ?,
			uploaded_by = ?
		WHERE institution_id = ?
	`,
		fileName,
		fileKey,
		uploadedBy,
		institutionID,
	).Error
}
