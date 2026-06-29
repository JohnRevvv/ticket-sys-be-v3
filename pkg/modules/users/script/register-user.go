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
			institution_name,
			role,
			job_position,
			status,
			created_at
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, 'User', ?, 'pending', ?)
	`,
		user.Username,
		user.StaffID,
		user.FirstName,
		user.LastName,
		user.Email,
		user.PhoneNo,
		user.InstitutionID,
		user.InstitutionName,
		user.JobPosition,
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

// func UserExists(staffID, email, phoneNo string) (bool, error) {
// 	db, err := getDB()
// 	if err != nil {
// 		return false, err
// 	}

// 	var count int64

// 	err = db.Table("users").
// 		Where("staff_id = ? OR email = ? OR phone_no = ?", staffID, email, phoneNo).
// 		Count(&count).Error

// 	if err != nil {
// 		return false, err
// 	}

// 	return count > 0, nil
// }

