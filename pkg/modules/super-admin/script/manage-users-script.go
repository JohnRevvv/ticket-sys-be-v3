package script

import (
	"ideyanale-be/pkg/config"
	"ideyanale-be/pkg/modules/users/model"
	)

func ChangeRoleToAdmin(user *model.UserDetails) error {
	return config.DBConnList[0].Exec(`
		UPDATE users
		SET role = 'insti-admin'
		WHERE id = ?
	`,
		user.ID,
	).Error
}
