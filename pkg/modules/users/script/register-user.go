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
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, 'user', ?, 'pending', ?)
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

