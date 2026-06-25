package script

import (
	"ideyanale-be/pkg/config"
	encrypDecryptV1 "ideyanale-be/pkg/middleware/encryption/v1"
	"ideyanale-be/pkg/modules/insti-admin/model"

	"gorm.io/gorm"
)

// ==========================
// AUTO-GENERATE SCRIPTS
// ==========================
var DefaultPositions = []string{
	"None",
}

func AddDefaultPositions(institutionID uint) error {
	return config.DBConnList[0].Transaction(func(tx *gorm.DB) error {

		for _, position := range DefaultPositions {

			encPosition, err := encrypDecryptV1.EncryptV1(position, config.SecretKey)
			if err != nil {
				return err
			}

			if err := tx.Exec(`
			INSERT INTO job_positions (
				position_name,
				institution_id
			)
			VALUES (?, ?)
		`,
				encPosition,
				institutionID,
			).Error; err != nil {
				return err
			}

		}
		return nil
		
	})
}

func AddPosition(position string, institutionID int) error {
	return config.DBConnList[0].Exec(`
		INSERT INTO job_positions (
			position_name,
			institution_id
		)
		VALUES (?, ?)
	`,
		position,
		institutionID,
	).Error
}

func GetPositionsByInstitutionID(institutionID int) ([]model.JobPosition, error) {
	var positions []model.JobPosition

	err := config.DBConnList[0].Raw(`
		SELECT
			position_id,
			position_name,
			institution_id
		FROM job_positions
		WHERE institution_id = ?
	`, institutionID).Scan(&positions).Error

	return positions, err
}

func UpdatePosition(positionID int, institutionID int, encPositionName string) (int64, error) {
	result := config.DBConnList[0].Exec(`
		UPDATE job_positions
		SET position_name = ?
		WHERE position_id = ?
		AND institution_id = ?
	`, encPositionName, positionID, institutionID)

	return result.RowsAffected, result.Error
}
