package script

import (
	"ideyanale-be/pkg/config"
	"ideyanale-be/pkg/modules/users/model"

	"time"
)

func CreateUser(user *model.UserDetails) error {
	return config.DBConnList[0].Exec(`
		INSERT INTO users (
			username,
			staff_id,
			first_name,
			last_name,
			email,
			phone_no,
			institution_id,
			role_id,
			position_id,
			status,
			created_at,
			updated_at
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`,
		user.Username,
		user.StaffID,
		user.FirstName,
		user.LastName,
		user.Email,
		user.PhoneNo,
		user.InstitutionID,
		user.RoleID,
		user.PositionID,
		user.Status,
		time.Now(),
		time.Now(),
	).Error
}

func UserExists(staffID, email, phoneNo string) (bool, error) {
	var exists bool

	err := config.DBConnList[0].Raw(`
		SELECT EXISTS(
			SELECT 1
			FROM users
			WHERE staff_id = ?
			   OR email = ?
			   OR phone_no = ?
		)
	`,
		staffID,
		email,
		phoneNo,
	).Scan(&exists).Error

	if err != nil {
		return false, err
	}

	return exists, nil
}
