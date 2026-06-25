package script

import (
	"fmt"
	"ideyanale-be/pkg/config"
	"ideyanale-be/pkg/modules/users/model"
	"time"

	"gorm.io/gorm"
)

func getDB() (*gorm.DB, error) {
	if len(config.DBConnList) == 0 || config.DBConnList[0] == nil {
		return nil, fmt.Errorf("database not initialized")
	}
	return config.DBConnList[0], nil
}

func GetActiveUserByStaffID(encStaffID string) (*model.UserDetails, error) {
	var user model.UserDetails

	db := config.DBConnList[0]

	err := db.
		Table("users").
		Where("staff_id = ? AND status = ?", encStaffID, "active").
		First(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func SaveLoginOTP(data *model.LoginOTP) error {
	db := config.DBConnList[0]

	return db.Table("login_otps").Create(data).Error
}

func DeleteOTPByStaffID(encStaffID string) error {
	db := config.DBConnList[0]

	return db.Table("login_otps").
		Where("staff_id = ?", encStaffID).
		Delete(&model.LoginOTP{}).Error
}

func GetOTPByStaffID(encStaffID string) (*model.LoginOTP, error) {
	var otp model.LoginOTP

	db := config.DBConnList[0]

	err := db.Table("login_otps").
		Where("staff_id = ?", encStaffID).
		First(&otp).Error

	if err != nil {
		return nil, err
	}

	return &otp, nil
}

func SetUserLoginStatus(ID int, isLoggedIn bool) error {
	return config.DBConnList[0].Exec(`
		UPDATE users
		SET 
			last_login = ?,
			is_logged_in = ?
		WHERE id = ?
	`,
		time.Now(),
		isLoggedIn,
		ID,
	).Error
}
