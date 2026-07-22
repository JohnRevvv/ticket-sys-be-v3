package controller

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v3"

	"ideyanale-be/pkg/config"
	global "ideyanale-be/pkg/global/json_response"
	encrypDecryptV1 "ideyanale-be/pkg/middleware/encryption/v1"
	"ideyanale-be/pkg/middleware/jwt"
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

	type RoleResp struct {
		RoleID           uint   `json:"role_id"`
		RoleName         string `json:"role_name"`
		CanCreateTicket  bool   `json:"can_create"`
		CanEndorseTicket bool   `json:"can_endorse"`
		CanApproveTicket bool   `json:"can_approve"`
		CanResolveTicket bool   `json:"can_resolve"`
		CanAudit         bool   `json:"can_audit"`
	}

	type JobPositionResp struct {
		PositionID   uint   `json:"position_id"`
		PositionName string `json:"position_name"`
	}

	type UserDetailsResp struct {
		ID              int             `json:"id,omitempty"`
		Username        string          `json:"username,omitempty"`
		StaffID         string          `json:"staff_id,omitempty"`
		FirstName       string          `json:"first_name,omitempty"`
		LastName        string          `json:"last_name,omitempty"`
		Email           string          `json:"email,omitempty"`
		PhoneNo         string          `json:"phone_no,omitempty"`
		InstitutionID   uint            `json:"institution_id,omitempty"`
		InstitutionName string          `json:"institution_name,omitempty"`
		JobPosition     JobPositionResp `json:"job_positions"`
		Role            RoleResp        `json:"role"`
		Status          string          `json:"status"`
		LastLogin       string          `json:"last_login,omitempty"`
		IsLoggedIn      bool            `json:"is_logged_in,omitempty"`
		CreatedAt       time.Time       `json:"created_at"`
	}

	response := make([]UserDetailsResp, 0, len(users))

	for _, user := range users {

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

		// 				fmt.Println(user.RoleID)
		// fmt.Println(user.Role.RoleID)

		response = append(response, UserDetailsResp{
			ID:            user.ID,
			Username:      decryptedUsername,
			StaffID:       decryptedStaffID,
			FirstName:     decryptedFirstName,
			LastName:      decryptedLastName,
			Email:         decryptedEmail,
			PhoneNo:       decryptedPhoneNo,
			InstitutionID: user.InstitutionID,
			
			JobPosition: JobPositionResp{
				PositionID:   user.JobPosition.PositionID,
				PositionName: user.JobPosition.PositionName,
			},

			Role: RoleResp{
				RoleID:           user.Role.RoleID,
				RoleName:         user.Role.RoleName,
				CanCreateTicket:  user.Role.CanCreateTicket,
				CanEndorseTicket: user.Role.CanEndorseTicket,
				CanApproveTicket: user.Role.CanApproveTicket,
				CanResolveTicket: user.Role.CanResolveTicket,
				CanAudit:         user.Role.CanAudit,
			},
			Status:     user.Status,
			LastLogin:  user.LastLogin,
			IsLoggedIn: user.IsLoggedIn,
			CreatedAt:  user.CreatedAt,
		})
	}

	return global.JSONResponseWithDataV1(c, "200", "Users fetched successfully", response, 200)
}

func GetAllUsers(c fiber.Ctx) error {

	users, err := script.GetAllUsers()
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Failed to fetch users", err, 500)
	}

	type RoleResp struct {
		RoleID           uint   `json:"role_id"`
		RoleName         string `json:"role_name"`
		CanCreateTicket  bool   `json:"can_create"`
		CanEndorseTicket bool   `json:"can_endorse"`
		CanApproveTicket bool   `json:"can_approve"`
		CanResolveTicket bool   `json:"can_resolve"`
		CanAudit         bool   `json:"can_audit"`
	}

	type UserDetailsResp struct {
		ID              int       `json:"id,omitempty"`
		Username        string    `json:"username,omitempty"`
		StaffID         string    `json:"staff_id,omitempty"`
		FirstName       string    `json:"first_name,omitempty"`
		LastName        string    `json:"last_name,omitempty"`
		Email           string    `json:"email,omitempty"`
		PhoneNo         string    `json:"phone_no,omitempty"`
		InstitutionID   uint      `json:"institution_id,omitempty"`
		InstitutionName string    `json:"institution_name,omitempty"`
		JobPositionID   uint      `json:"job_position_id,omitempty"`
		Role            RoleResp  `json:"role"`
		Status          string    `json:"status"`
		LastLogin       string    `json:"last_login,omitempty"`
		IsLoggedIn      bool      `json:"is_logged_in,omitempty"`
		CreatedAt       time.Time `json:"created_at"`
	}

	response := make([]UserDetailsResp, 0, len(users))

	for _, user := range users {

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

		response = append(response, UserDetailsResp{
			ID:            user.ID,
			Username:      decryptedUsername,
			StaffID:       decryptedStaffID,
			FirstName:     decryptedFirstName,
			LastName:      decryptedLastName,
			Email:         decryptedEmail,
			PhoneNo:       decryptedPhoneNo,
			InstitutionID: user.InstitutionID,
			JobPositionID: user.PositionID,
			Role: RoleResp{
				RoleID:           user.Role.RoleID,
				RoleName:         user.Role.RoleName,
				CanCreateTicket:  user.Role.CanCreateTicket,
				CanEndorseTicket: user.Role.CanEndorseTicket,
				CanApproveTicket: user.Role.CanApproveTicket,
				CanResolveTicket: user.Role.CanResolveTicket,
				CanAudit:         user.Role.CanAudit,
			},
			Status:     user.Status,
			LastLogin:  user.LastLogin,
			IsLoggedIn: user.IsLoggedIn,
			CreatedAt:  user.CreatedAt,
		})
	}

	return global.JSONResponseWithDataV1(c, "200", "Users fetched successfully", response, 200)
}

func GetUserByID(c fiber.Ctx) error {

	userID, err := strconv.Atoi(c.Params("id"))
	if err != nil || userID <= 0 {
		return global.JSONResponseWithErrorV1(c, "400", "Invalid user id", nil, 400)
	}

	user, err := script.GetUserByID(userID)
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Failed to fetch user", err, 500)
	}

	if user.ID == 0 {
		return global.JSONResponseWithErrorV1(c, "404", "User not found", nil, 404)
	}

	type RoleResp struct {
		RoleID           uint   `json:"role_id"`
		RoleName         string `json:"role_name"`
		CanCreateTicket  bool   `json:"can_create"`
		CanEndorseTicket bool   `json:"can_endorse"`
		CanApproveTicket bool   `json:"can_approve"`
		CanResolveTicket bool   `json:"can_resolve"`
		CanAudit         bool   `json:"can_audit"`
	}

	type UserDetailsResp struct {
		ID            int       `json:"id,omitempty"`
		Username      string    `json:"username,omitempty"`
		StaffID       string    `json:"staff_id,omitempty"`
		FirstName     string    `json:"first_name,omitempty"`
		LastName      string    `json:"last_name,omitempty"`
		Email         string    `json:"email,omitempty"`
		PhoneNo       string    `json:"phone_no,omitempty"`
		InstitutionID uint      `json:"institution_id,omitempty"`
		JobPositionID uint      `json:"job_position_id,omitempty"`
		Role          RoleResp  `json:"role"`
		Status        string    `json:"status"`
		LastLogin     string    `json:"last_login,omitempty"`
		IsLoggedIn    bool      `json:"is_logged_in,omitempty"`
		CreatedAt     time.Time `json:"created_at"`
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

	resp := UserDetailsResp{
		ID:            user.ID,
		Username:      decryptedUsername,
		StaffID:       decryptedStaffID,
		FirstName:     decryptedFirstName,
		LastName:      decryptedLastName,
		Email:         decryptedEmail,
		PhoneNo:       decryptedPhoneNo,
		InstitutionID: user.InstitutionID,
		JobPositionID: user.PositionID,
		Role: RoleResp{
			RoleID:           user.Role.RoleID,
			RoleName:         user.Role.RoleName,
			CanCreateTicket:  user.Role.CanCreateTicket,
			CanEndorseTicket: user.Role.CanEndorseTicket,
			CanApproveTicket: user.Role.CanApproveTicket,
			CanResolveTicket: user.Role.CanResolveTicket,
			CanAudit:         user.Role.CanAudit,
		},
		Status:     user.Status,
		LastLogin:  user.LastLogin,
		IsLoggedIn: user.IsLoggedIn,
		CreatedAt:  user.CreatedAt,
	}

	return global.JSONResponseWithDataV1(c, "200", "User fetched successfully", resp, 200)
}

func CountUsers(c fiber.Ctx) error {
	if err := jwt.RequireRoles(c, "Super-Admin", "Insti-Admin"); err != nil {
		return global.JSONResponseWithErrorV1(c, "403", "Forbidden", err, fiber.StatusForbidden)
	}

	inst := c.Locals("institution_id")
	if inst == nil {
		return global.JSONResponseWithErrorV1(c, "401", "Unauthorized institution", nil, fiber.StatusUnauthorized)
	}

	institutionID, ok := inst.(int)
	if !ok {
		return global.JSONResponseWithErrorV1(c, "500", "Invalid institution id type", nil, fiber.StatusInternalServerError)
	}

	count, err := script.CountUsersByInstitutionID(institutionID)
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Failed to count users", err, fiber.StatusInternalServerError)
	}

	return global.JSONResponseWithDataV1(
		c,
		"200",
		"User count retrieved successfully",
		fiber.Map{
			"count": count,
		},
		fiber.StatusOK,
	)
}
