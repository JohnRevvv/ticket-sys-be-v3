package script

import (
	"ideyanale-be/pkg/modules/users/model"
	"ideyanale-be/pkg/config"
)

func GetUsersByInstitutionID(ID string) ([]model.UserDetails, error) {
	var users []model.UserDetails

	err := config.DBConnList[0].Raw(`
		SELECT 
			id,
			username,
			staff_id,
			first_name,
			last_name,
			email,
			phone_no,
			institution_id,
			institution_name,
			job_position,
			role,
			status,
			last_login,
			is_logged_in,
			created_at
		FROM users
		WHERE institution_id = ?
		ORDER BY id DESC
	`, ID).Scan(&users).Error

	return users, err
}

// func GetUserByID(ID string) (*model.UserDetails, error) {
	
// 	var user model.UserDetails

// 	err := config.DBConnList[0].Table("users").Where("id = ?", ID).First(&user).Error
// 	if err != nil {
// 		return nil, err
// 	}
	
// 	return &user, nil
// }

func GetUserByID(userID int) (model.UserDetails, error) {
	var user model.UserDetails

	err := config.DBConnList[0].Raw(`
		SELECT *
		FROM users
		WHERE id = ?
	`, userID).Scan(&user).Error

	return user, err
}