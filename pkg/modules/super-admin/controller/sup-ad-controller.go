package controller

import (
	"errors"
	"ideyanale-be/pkg/config"
	global "ideyanale-be/pkg/global/json_response"
	encrypDecryptV1 "ideyanale-be/pkg/middleware/encryption/v1"
	hashingV1 "ideyanale-be/pkg/middleware/hashing/v1"
	"ideyanale-be/pkg/middleware/jwt"
	"strconv"
	"strings"

	SAdmodel "ideyanale-be/pkg/modules/super-admin/model"
	"ideyanale-be/pkg/modules/super-admin/script"

	"github.com/gofiber/fiber/v3"
)

func CreateSuperAdmin(c fiber.Ctx) error {

	type Req struct {
		UserName string `json:"username"`
		Role     string `json:"role"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var req Req

	if err := c.Bind().Body(&req); err != nil {
		return global.JSONResponseWithErrorV1(c, "400", "Invalid request", err, 400)
	}

	username := strings.TrimSpace(strings.ToLower(req.UserName))
	email := strings.TrimSpace(strings.ToLower(req.Email))

	// encrypt username/email (storage only)
	encUserName, err := encrypDecryptV1.EncryptV1(username, config.SecretKey)
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Encrypt username failed", err, 500)
	}

	encEmail, err := encrypDecryptV1.EncryptV1(email, config.SecretKey)
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Encrypt email failed", err, 500)
	}

	// bcrypt password
	hashedPassword, err := hashingV1.GenerateHash(req.Password)
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Password hash failed", err, 500)
	}

	// save
	err = script.CreateSuperAdmin(
		encUserName,
		req.Role,
		encEmail,
		hashedPassword,
	)

	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Create user failed", err, 500)
	}

	return global.JSONResponseV1(c, "200", "User created successfully", 200)
}

func LoginSuperAdmin(c fiber.Ctx) error {

	type LoginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var req LoginRequest

	if err := c.Bind().Body(&req); err != nil {
		return global.JSONResponseWithErrorV1(c, "400", "Invalid request", err, 400)
	}

	input := strings.TrimSpace(strings.ToLower(req.Username))

	users, err := script.GetAllSuperAdmins()
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Failed to fetch users", err, 500)
	}

	var found *SAdmodel.SuperAccount

	for i := range users {

		u := &users[i]

		decUsername, err1 := encrypDecryptV1.DecryptV1(u.UserName, config.SecretKey)
		decEmail, err2 := encrypDecryptV1.DecryptV1(u.Email, config.SecretKey)

		if err1 == nil && strings.ToLower(decUsername) == input {
			found = u
			break
		}

		if err2 == nil && strings.ToLower(decEmail) == input {
			found = u
			break
		}
	}

	if found == nil {
		return global.JSONResponseWithErrorV1(c, "401", "Invalid username/email or password", errors.New("user not found"), 401)
	}

	// bcrypt password check
	if !hashingV1.ValidateHash(req.Password, found.Password) {
		return global.JSONResponseWithErrorV1(c, "401", "Invalid username/email or password", errors.New("invalid password"), 401)
	}

	// generate JWT token
	token, err := jwt.GenerateSuperAdminToken(
		found.ID,
		found.UserName,
	)

	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Failed to generate token", err, 500)
	}

	// bcrypt password check
	if !hashingV1.ValidateHash(req.Password, found.Password) {
		return global.JSONResponseWithErrorV1(
			c,
			"401",
			"Invalid username/email or password",
			errors.New("invalid password"),
			401,
		)
	}

	// update login status
	if err := script.SetSuperAdminLoginStatus(int(found.ID), true); err != nil {
		return global.JSONResponseWithErrorV1(
			c,
			"500",
			"Failed to update login status",
			err,
			500,
		)
	}

	found.Password = ""

	// response wrapper
	type LoginResponse struct {
		User  *SAdmodel.SuperAccount `json:"user"`
		Token string                 `json:"token"`
	}

	return global.JSONResponseWithDataV1(c, "200", "Login successful",
		LoginResponse{
			User:  found,
			Token: token,
		}, 200)
}

func LogoutSuperAdmin(c fiber.Ctx) error {

	// get userID from context (set by JWT middleware)
	ID, ok := c.Locals("id").(int)
	if !ok {
		return global.JSONResponseWithErrorV1(c, "401", "Unauthorized", nil, 401)
	}

	// update DB logout status
	if err := script.LogoutSuperAdmin(ID); err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Failed to logout user", err, 500)
	}

	return global.JSONResponseWithDataV1(c, "200", "Logout successful", nil, 200)
}

func GetSuperAdminByID(c fiber.Ctx) error {

	if err := jwt.RequireRoles(c, "Super-Admin"); err != nil {
		return global.JSONResponseWithErrorV1(c, "403", "Forbidden", err, 403)
	}

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id <= 0 {
		return global.JSONResponseWithErrorV1(c, "400", "Invalid super admin ID", nil, 400)
	}

	superAdmin, err := script.GetSuperAdminByID(id)
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Failed to retrieve super admin", err, 500)
	}

	if superAdmin == nil || superAdmin.ID == 0 {
		return global.JSONResponseWithErrorV1(c, "404", "Super admin not found", nil, 404)
	}

	// Decrypt username
	superAdmin.UserName, err = encrypDecryptV1.DecryptV1(superAdmin.UserName, config.SecretKey)
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Failed to decrypt username", err, 500)
	}

	// Decrypt email
	superAdmin.Email, err = encrypDecryptV1.DecryptV1(superAdmin.Email, config.SecretKey)
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Failed to decrypt email", err, 500)
	}

	return global.JSONResponseWithDataV1(
		c,
		"200",
		"Super admin retrieved successfully",
		superAdmin,
		200,
	)
}
