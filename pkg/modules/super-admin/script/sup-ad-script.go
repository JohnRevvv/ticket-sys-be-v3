package script

import (
	"ideyanale-be/pkg/config"
	"ideyanale-be/pkg/modules/super-admin/model"
)

func CreateSuperAdmin(encusername, role, encemail, password string) error {

	return config.DBConnList[0].Exec(`
		INSERT INTO superaccount (
			username,
			role,
			email,
			password
		)
		VALUES (?, 'super-admin', ?, ?)
	`,
		encusername,
		encemail,
		password,
	).Error
}

func GetAllSuperAdmins() ([]model.SuperAccount, error) {

	var users []model.SuperAccount

	err := config.DBConnList[0].Raw(`
		SELECT id, username, 
		role, email, 
		password
		FROM superaccount
	`).Scan(&users).Error

	return users, err
}

func GetSuperAdminByID(id int) (*model.SuperAccount, error) {

	var user model.SuperAccount

	err := config.DBConnList[0].Raw(`
		SELECT id, username, role, email, password
		FROM superaccount
		WHERE id = ?
	`, id).Scan(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func LogoutSuperAdmin(id int) error {
	return config.DBConnList[0].Exec(`
		UPDATE superaccount
		SET is_logged_in = false
		WHERE id = ?
	`, id).Error
}

func IsSuperAdminLoggedIn(id int) (bool, error) {
	var loggedIn bool

	err := config.DBConnList[0].Raw(`
		SELECT is_logged_in
		FROM superaccount
		WHERE id = ?
	`, id).Scan(&loggedIn).Error

	return loggedIn, err
}

func SetSuperAdminLoginStatus(id int, loggedIn bool) error {
	return config.DBConnList[0].Exec(`
		UPDATE superaccount
		SET is_logged_in = ?
		WHERE id = ?
	`, loggedIn, id).Error
}