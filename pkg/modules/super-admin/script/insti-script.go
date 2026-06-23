package script

import (
	"fmt"
	"ideyanale-be/pkg/config"
	"ideyanale-be/pkg/modules/super-admin/model"
	"time"

	"gorm.io/gorm"
)

func getDB() (*gorm.DB, error) {
	if len(config.DBConnList) == 0 || config.DBConnList[0] == nil {
		return nil, fmt.Errorf("database not initialized")
	}
	return config.DBConnList[0], nil
}

func AddInstitution(encinstitutionName, encinstitutionCode string) error {
	return config.DBConnList[0].Exec(`
		INSERT INTO institutions (
			institution_name,
			description,
			created_at
		)
		VALUES (?, ?, ?)
	`,
		encinstitutionName,
		encinstitutionCode,
		time.Now(),
	).Error
}

func GetInstitutions() ([]model.Institution, error) {
	var institutions []model.Institution

	err := config.DBConnList[0].Raw(`
		SELECT
			institution_id,
			institution_name,
			description,
			created_at
		FROM institutions
	`).Scan(&institutions).Error

	return institutions, err
}

