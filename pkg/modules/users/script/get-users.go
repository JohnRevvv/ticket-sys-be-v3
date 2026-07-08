package script

import (
	"fmt"
	"ideyanale-be/pkg/config"
	"ideyanale-be/pkg/modules/users/model"
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

func RoleBelongsToInstitution(roleID, institutionID uint) (bool, error) {
	var exists bool

	err := config.DBConnList[0].Raw(`
		SELECT EXISTS (
			SELECT 1
			FROM roles
			WHERE role_id = ?
			AND institution_id = ?
			AND deleted_at IS NULL
		)
	`, roleID, institutionID).Scan(&exists).Error

	return exists, err
}

func GetUserByID(userID int) (model.UserDetails, error) {
	var user model.UserDetails

	err := config.DBConnList[0].
		Preload("Role").
		Where("id = ? AND deleted_at IS NULL", userID).
		First(&user).Error

	return user, err
}

func CountUsers() (int64, error) {
	var count int64

	err := config.DBConnList[0].Raw(`
		SELECT COUNT(*)
		FROM users
	`).Scan(&count).Error

	return count, err
}

func CountUsersByInstitutionID(institutionID int) (int64, error) {
	var count int64

	err := config.DBConnList[0].Raw(`
		SELECT COUNT(*)
		FROM users
		WHERE institution_id = ?
	`, institutionID).Scan(&count).Error

	return count, err
}

func GetDefaultUserRoleID(institutionID uint) (uint, error) {
	var roleID uint

	err := config.DBConnList[0].Raw(`
		SELECT role_id
		FROM roles
		WHERE institution_id = ?
		AND role_name = 'User'
		AND deleted_at IS NULL
		LIMIT 1
	`, institutionID).Scan(&roleID).Error

	if err != nil {
		return 0, err
	}

	if roleID == 0 {
		return 0, fmt.Errorf("default 'User' role not found for institution %d — did AddDefaultRoles run?", institutionID)
	}

	return roleID, nil
}