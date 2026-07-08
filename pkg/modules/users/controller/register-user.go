package controller

import (

	"strconv"
	"strings"

	"ideyanale-be/pkg/config"
	global "ideyanale-be/pkg/global/json_response"
	encrypDecryptV1 "ideyanale-be/pkg/middleware/encryption/v1"

	"ideyanale-be/pkg/modules/users/model"
	"ideyanale-be/pkg/modules/users/script"

	"github.com/gofiber/fiber/v3"
)

func RegisterUser(c fiber.Ctx) error {

	type Req struct {
		StaffID         string `json:"staff_id"`
		Email           string `json:"email"`
		FirstName       string `json:"first_name"`
		LastName        string `json:"last_name"`
		PhoneNo         string `json:"phone_no"`
		InstitutionID   uint   `json:"institution_id"`
		JobPosition     string `json:"job_position"`
		Status          string `json:"status"`
		// RoleID removed — always defaults to "User" on registration
	}

	var req Req

	// bind request
	if err := c.Bind().Body(&req); err != nil {
		return global.JSONResponseWithErrorV1(c, "400", "Invalid request", err, 400)
	}

	// -------------------------
	// Normalize inputs
	// -------------------------
	staffID := strings.TrimSpace(req.StaffID)
	staffID = strings.ReplaceAll(staffID, "-", "")

	// validate length
	if len(staffID) != 11 {
		return global.JSONResponseWithErrorV1(c, "400", "Staff ID must be exactly 11 digits", nil, 400)
	}

	// validate numeric only
	if _, err := strconv.Atoi(staffID); err != nil {
		return global.JSONResponseWithErrorV1(c, "400", "Staff ID must contain numbers only", nil, 400)
	}

	// format: 202511-12345
	staffID = staffID[:6] + "-" + staffID[6:]

	firstName := strings.TrimSpace(req.FirstName)
	lastName := strings.TrimSpace(req.LastName)
	email := strings.ToLower(strings.TrimSpace(req.Email))
	phoneNo := strings.TrimSpace(req.PhoneNo)
	jobPosition := strings.TrimSpace(req.JobPosition)

	// validations
	emailParts := strings.Split(email, "@")
	if len(emailParts) != 2 {
		return global.JSONResponseWithErrorV1(c, "400", "Invalid email format", nil, 400)
	}

	username := emailParts[0]

	if staffID == "" {
		return global.JSONResponseWithErrorV1(c, "400", "Staff ID is required", nil, 400)
	}

	if firstName == "" {
		return global.JSONResponseWithErrorV1(c, "400", "Firstname is required", nil, 400)
	}

	if lastName == "" {
		return global.JSONResponseWithErrorV1(c, "400", "Lastname is required", nil, 400)
	}

	if email == "" {
		return global.JSONResponseWithErrorV1(c, "400", "Email is required", nil, 400)
	}

	if phoneNo == "" {
		return global.JSONResponseWithErrorV1(c, "400", "Phone number is required", nil, 400)
	}

	if req.InstitutionID == 0 {
		return global.JSONResponseWithErrorV1(c, "400", "Institution ID is required", nil, 400)
	}

	// default status if not supplied
	status := strings.TrimSpace(req.Status)
	if status == "" {
		status = "active"
	}

	// NEW: resolve the default "User" role for this institution — never trust client-supplied role_id
	defaultRoleID, err := script.GetDefaultUserRoleID(req.InstitutionID)
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Failed to resolve default role", err, 500)
	}

	// -------------------------
	// Encryption layer
	// -------------------------
	encrypt := func(val string) (string, error) {
		return encrypDecryptV1.EncryptV1(val, config.SecretKey)
	}

	encStaffID, err := encrypt(staffID)
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Encrypt staff ID failed", err, 500)
	}

	encUserName, err := encrypt(username)
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Encrypt staff ID failed", err, 500)
	}

	encFirstName, err := encrypt(firstName)
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Encrypt first name failed", err, 500)
	}

	encLastName, err := encrypt(lastName)
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Encrypt last name failed", err, 500)
	}

	encEmail, err := encrypt(email)
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Encrypt email failed", err, 500)
	}

	encPhoneNo, err := encrypt(phoneNo)
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Encrypt phone number failed", err, 500)
	}

	// check duplicates after encryption
	exists, err := script.UserExists(encStaffID, encEmail, encPhoneNo)
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Validation failed", err, 500)
	}

	if exists {
		return global.JSONResponseWithErrorV1(c, "409", "Staff ID, email, or phone number already exists", nil, 409)
	}

	// -------------------------
	// DB layer (script)
	// -------------------------
	err = script.CreateUser(&model.UserDetails{
		StaffID:         encStaffID,
		Username:        encUserName,
		FirstName:       encFirstName,
		LastName:        encLastName,
		Email:           encEmail,
		PhoneNo:         encPhoneNo,
		InstitutionID:   req.InstitutionID,
		JobPosition:     jobPosition,
		RoleID:          defaultRoleID, // always "User" on self-registration
		Status:          status,
	})

	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Create user failed", err, 500)
	}

	// -------------------------
	// success response
	// -------------------------
	return global.JSONResponseV1(c, "200", "User created successfully", 200)
}

