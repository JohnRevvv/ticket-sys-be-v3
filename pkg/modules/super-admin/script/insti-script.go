package script

import (
	"fmt"
	"ideyanale-be/pkg/config"
	SAdmodel "ideyanale-be/pkg/modules/super-admin/model"
	"time"

	"gorm.io/gorm"
)

func getDB() (*gorm.DB, error) {
	if len(config.DBConnList) == 0 || config.DBConnList[0] == nil {
		return nil, fmt.Errorf("database not initialized")
	}
	return config.DBConnList[0], nil
}

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

func GetInstitutions() ([]SAdmodel.Institution, error) {
	var institutions []SAdmodel.Institution

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
