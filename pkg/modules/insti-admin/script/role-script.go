package script

import ( 
	"ideyanale-be/pkg/config"
	InsAdmodel "ideyanale-be/pkg/modules/insti-admin/model"
)


func AddRole(roleName string, institutionID int, canCreate, canEndorse, canApprove,  canResolve, canAudit bool,) error {

	return config.DBConnList[0].Exec(`
		INSERT INTO roles (
			institution_id,
			role_name,
			can_create,
			can_endorse,
			can_approve,
			can_resolve,
			can_audit,
			created_at,
			updated_at
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
	`,
		institutionID,
		roleName,
		canCreate,
		canEndorse,
		canApprove,
		canResolve,
		canAudit,
	).Error
}

func ExistingRole(institutionID int, roleName string) (bool, error) {
	var exists bool

	err := config.DBConnList[0].Raw(`
		SELECT EXISTS (
			SELECT 1
			FROM roles
			WHERE institution_id = ?
			AND LOWER(TRIM(role_name)) = LOWER(TRIM(?))
		)
	`, institutionID, roleName).Scan(&exists).Error

	return exists, err
}

func GetRolesByInstitutionID(institutionID int) ([]InsAdmodel.Roles, error) {
	var roles []InsAdmodel.Roles

	err := config.DBConnList[0].Raw(`
		SELECT *
		FROM roles
		WHERE institution_id = ?
	`, institutionID).Scan(&roles).Error

	return roles, err
}