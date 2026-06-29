package script

import (
	"fmt"
	"ideyanale-be/pkg/config"
	"ideyanale-be/pkg/modules/insti-admin/model"
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

func ChangeUserRole(userID int, role string) error {
	return config.DBConnList[0].Exec(`
		UPDATE users
		SET role = ?,
			updated_at = NOW()
		WHERE id = ?
	`, role, userID).Error
}

func GetRoleByName(roleName string) (model.Roles, error) {
	var role model.Roles

	err := config.DBConnList[0].Raw(`
		SELECT *
		FROM roles
		WHERE role_name = ?
		LIMIT 1
	`, roleName).Scan(&role).Error

	return role, err
}

func GetRoleByID(roleID uint) (model.Roles, error) {
	var role model.Roles

	err := config.DBConnList[0].Raw(`
		SELECT *
		FROM roles
		WHERE role_id = ?
	`, roleID).Scan(&role).Error

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
