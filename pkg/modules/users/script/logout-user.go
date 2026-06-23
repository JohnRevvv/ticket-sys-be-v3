package script

import (
	"ideyanale-be/pkg/config"
	
)

func LogoutUser(ID int) error {
	return config.DBConnList[0].Exec(`
		UPDATE users
		SET is_logged_in = ?
		WHERE id = ?
	`,
		false,
		ID,
	).Error
}