package script

import (
	"ideyanale-be/pkg/config"
	"ideyanale-be/pkg/modules/super-admin/model"
)

func CreateSuperAdmin(encusername, role, encemail, password string) error {

	return config.DBConnList[0].Exec(`
		INSERT INTO admin (
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

func GetAllSuperAdmins() ([]model.SuperAdminDetails, error) {

	var users []model.SuperAdminDetails

	err := config.DBConnList[0].Raw(`
		SELECT id, username, 
		role, email, 
		password
		FROM admin
	`).Scan(&users).Error

	return users, err
}

func GetSuperAdminByID(id int) (*model.SuperAdminDetails, error) {

	var user model.SuperAdminDetails

	err := config.DBConnList[0].Raw(`
		SELECT id, username, role, email, password
		FROM admin
		WHERE id = ?
	`, id).Scan(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}
