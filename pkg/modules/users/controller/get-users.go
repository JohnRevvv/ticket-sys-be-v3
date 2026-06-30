package controller

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v3"

	"ideyanale-be/pkg/config"
	global "ideyanale-be/pkg/global/json_response"
	encrypDecryptV1 "ideyanale-be/pkg/middleware/encryption/v1"
	script "ideyanale-be/pkg/modules/users/script"
)

func GetUsersByInstitutionID(c fiber.Ctx) error {

	ID := c.Params("institution_id")
	if ID == "" {
		return global.JSONResponseWithErrorV1(c, "400", "Institution ID is required", nil, 400)
	}

	users, err := script.GetUsersByInstitutionID(ID)
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Failed to fetch users", err, 500)
	}

	return global.JSONResponseWithDataV1(c, "200", "Users fetched successfully", users, 200)
}

func GetUserByID(c fiber.Ctx) error {

	userID, err := strconv.Atoi(c.Params("id"))
	if err != nil || userID <= 0 {
		return global.JSONResponseWithErrorV1(
			c,
			"400",
			"Invalid user id",
			nil,
			400,
		)
	}

	user, err := script.GetUserByID(userID)
	if err != nil {
		return global.JSONResponseWithErrorV1(
			c,
			"500",
			"Failed to fetch user",
			err,
			500,
		)
	}

	if user.ID == 0 {
		return global.JSONResponseWithErrorV1(
			c,
			"404",
			"User not found",
			nil,
			404,
		)
	}

	type UserDetailsResp struct {
		ID              int       `json:"id,omitempty"`
		Username        string    `json:"username,omitempty"`
		StaffID         string    `json:"staff_id,omitempty"`
		FirstName       string    `json:"first_name,omitempty"`
		LastName        string    `json:"last_name,omitempty"`
		Email           string    `json:"email,omitempty"`
		PhoneNo         string    `json:"phone_no,omitempty"`
		InstitutionID   int       `json:"institution_id,omitempty"`
		InstitutionName string    `json:"institution_name,omitempty"`
		JobPosition     string    `json:"job_position,omitempty"`
		Role            string    `json:"role"`
		Status          string    `json:"status"`
		LastLogin       string    `json:"last_login,omitempty"`
		IsLoggedIn      bool      `json:"is_logged_in,omitempty"`
		CreatedAt       time.Time `json:"created_at"`
	}

	decryptedUsername, err := encrypDecryptV1.DecryptV1(user.Username, config.SecretKey)
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Decrypt username failed", err, 500)
	}

	decryptedStaffID, err := encrypDecryptV1.DecryptV1(user.StaffID, config.SecretKey)
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Decrypt staff ID failed", err, 500)
	}

	decryptedFirstName, err := encrypDecryptV1.DecryptV1(user.FirstName, config.SecretKey)
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Decrypt first name failed", err, 500)
	}

	decryptedLastName, err := encrypDecryptV1.DecryptV1(user.LastName, config.SecretKey)
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Decrypt last name failed", err, 500)
	}

	decryptedEmail, err := encrypDecryptV1.DecryptV1(user.Email, config.SecretKey)
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Decrypt email failed", err, 500)
	}

	decryptedPhoneNo, err := encrypDecryptV1.DecryptV1(user.PhoneNo, config.SecretKey)
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Decrypt phone number failed", err, 500)
	}

	decryptedInstitutionName, err := encrypDecryptV1.DecryptV1(user.InstitutionName, config.SecretKey)
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Decrypt institution name failed", err, 500)
	}

	resp := UserDetailsResp{
		ID:              user.ID,
		Username:        decryptedUsername,
		StaffID:         decryptedStaffID,
		FirstName:       decryptedFirstName,
		LastName:        decryptedLastName,
		Email:           decryptedEmail,
		PhoneNo:         decryptedPhoneNo,
		InstitutionID:   user.InstitutionID,
		InstitutionName: decryptedInstitutionName,
		JobPosition:     user.JobPosition,
		Role:            user.Role,
		Status:          user.Status,
		LastLogin:       user.LastLogin,
		IsLoggedIn:      user.IsLoggedIn,
		CreatedAt:       user.CreatedAt,
	}

	return global.JSONResponseWithDataV1(
		c,
		"200",
		"User fetched successfully",
		resp,
		200,
	)
}

// helper (keeps controller clean)
func now() time.Time {
	return time.Now()
}
