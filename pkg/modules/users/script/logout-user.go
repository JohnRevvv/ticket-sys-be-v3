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

func IsLoggedIn(id int) (bool, error) {
	var isLoggedIn bool
	err := config.DBConnList[0].Raw(`
		SELECT is_logged_in
		FROM users
		WHERE id = ?
	`,
		id,
	).Scan(&isLoggedIn).Error

	return isLoggedIn, err
}