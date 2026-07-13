package script

import (
	"ideyanale-be/pkg/config"
	"time"
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
