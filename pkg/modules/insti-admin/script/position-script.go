package script

import (
	"ideyanale-be/pkg/config"
	"ideyanale-be/pkg/modules/insti-admin/model"
)

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