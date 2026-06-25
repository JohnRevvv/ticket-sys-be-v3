package script

import (
	"ideyanale-be/pkg/config"
)

func ChangeRoleToAdmin(userID int)error {
	return config.DBConnList[0].Exec(`
		UPDATE users
		SET role = 'insti-admin'
		WHERE id = ?
	`,
		userID,
	).Error
}

func ChangeUserStatus(userID int, status string) error {
	return config.DBConnList[0].Exec(`
		UPDATE users
		SET status = ?
		WHERE id = ?
	`,
		status,
		userID,
	).Error
}
