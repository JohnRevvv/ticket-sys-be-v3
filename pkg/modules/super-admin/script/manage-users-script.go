package script

import (
	"fmt"
	"ideyanale-be/pkg/config"
	Rolemodel "ideyanale-be/pkg/modules/roles/model"
)

func ChangeRoleToAdmin(userID int)error {
	return config.DBConnList[0].Exec(`
		UPDATE users
		SET role = 'Insti-Admin'
		WHERE id = ?
	`,
		userID,
	).Error
}

func GetRoleByID(roleID uint) (Rolemodel.Roles, error) {
	var role Rolemodel.Roles

	err := config.DBConnList[0].
		Where("role_id = ? AND deleted_at IS NULL", roleID).
		First(&role).Error

	return role, err
}

func ChangeUserRole(userID int, roleID uint) error {
	return config.DBConnList[0].Exec(`
		UPDATE users
		SET role_id = ?,
			updated_at = NOW()
		WHERE id = ?
	`, roleID, userID).Error
}

func GetRoleByName(roleName string) (Rolemodel.Roles, error) {
	var role Rolemodel.Roles

	err := config.DBConnList[0].Raw(`
		SELECT *
		FROM roles
		WHERE role_name = ?
		LIMIT 1
	`, roleName).Scan(&role).Error

	return role, err
}

func ChangeUserStatus(userID int, status string) error {
	result := config.DBConnList[0].Exec(`
		UPDATE users
		SET status = ?
		WHERE id = ?
	`, status, userID)

	fmt.Println("Rows:", result.RowsAffected)
	fmt.Println("Error:", result.Error)

	return result.Error
}
