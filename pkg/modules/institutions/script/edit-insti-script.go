package script

import (
	"fmt"
	"ideyanale-be/pkg/config"
	"ideyanale-be/pkg/modules/institutions/model"
)

func ChangeInstitutionStatus(institutionID uint, status string) error {
	result := config.DBConnList[0].Exec(`
		UPDATE institutions
		SET status = ?
		WHERE institution_id = ?
	`, status, institutionID)

	fmt.Println("Rows:", result.RowsAffected)
	fmt.Println("Error:", result.Error)

	return result.Error
}

func UpdateInstitution(institutionID uint, institutionCode, institutionName, description, institutionColor string) error {

	return config.DBConnList[0].Exec(`
		UPDATE institutions
		SET
			institution_code = ?,
			institution_name = ?,
			description = ?,
			institution_color = ?,
			status = ?
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
