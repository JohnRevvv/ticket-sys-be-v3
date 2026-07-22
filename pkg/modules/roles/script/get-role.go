package script

import (
	"fmt"
	"ideyanale-be/pkg/config"
	Rolemodel "ideyanale-be/pkg/modules/roles/model"
)

func ExistingRole(institutionID uint, roleName string) (bool, error) {
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

func GetRoleByID(roleID uint) (Rolemodel.Roles, error) {
	var role Rolemodel.Roles

	err := config.DBConnList[0].Raw(`
		SELECT
			role_id,
			institution_id,
			role_name,
			can_create,
			can_endorse,
			can_approve,
			can_resolve,
			can_audit,
			created_at,
			updated_at
		FROM roles
		WHERE role_id = ?
		AND deleted_at IS NULL
	`, roleID).Scan(&role).Error

	if err != nil {
		return role, err
	}

	if role.RoleID == 0 {
		return role, fmt.Errorf("role not found")
	}

	return role, nil
}

func GetRoleByUserID(userID uint) (*Rolemodel.Roles, error) {

	var role Rolemodel.Roles

	err := config.DBConnList[0].
		Table("roles r").
		Select("r.*").
		Joins("JOIN users u ON u.role_id = r.role_id").
		Where("u.id = ?", userID).
		First(&role).Error

	if err != nil {
		return nil, err
	}

	return &role, nil
}

func RoleNameExistsForInstitution(roleName string, institutionID uint, excludeRoleID uint) (bool, error) {
	var exists bool

	err := config.DBConnList[0].Raw(`
		SELECT EXISTS (
			SELECT 1
			FROM roles
			WHERE institution_id = ?
			AND LOWER(TRIM(role_name)) = LOWER(TRIM(?))
			AND role_id != ?
			AND deleted_at IS NULL
		)
	`, institutionID, roleName, excludeRoleID).Scan(&exists).Error

	return exists, err
}

func GetRolesByInstitutionID(institutionID uint) ([]Rolemodel.Roles, error) {
	var roles []Rolemodel.Roles

	err := config.DBConnList[0].Raw(`
		SELECT
			role_id,
			institution_id,
			role_name,
			can_create,
			can_endorse,
			can_approve,
			can_resolve,
			can_audit,
			created_at,
			updated_at
		FROM roles
		WHERE institution_id = ?
		AND deleted_at IS NULL
		ORDER BY role_name ASC
	`, institutionID).Scan(&roles).Error

	return roles, err
}
