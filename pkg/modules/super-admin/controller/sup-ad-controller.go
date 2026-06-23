package controller

import (
	"errors"
	"ideyanale-be/pkg/config"
	global "ideyanale-be/pkg/global/json_response"
	encrypDecryptV1 "ideyanale-be/pkg/middleware/encryption/v1"
	hashingV1 "ideyanale-be/pkg/middleware/hashing/v1"
	"ideyanale-be/pkg/middleware/jwt"
	"strings"

	SAdmodel "ideyanale-be/pkg/modules/super-admin/model"
	"ideyanale-be/pkg/modules/super-admin/script"

	"github.com/gofiber/fiber/v3"
)

func CreateSuperAdmin(c fiber.Ctx) error {

	type Req struct {
		UserName  string `json:"username"`
		Role      string `json:"role"`
		Email     string `json:"email"`
		Password  string `json:"password"`
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

	var found *SAdmodel.SuperAdminDetails

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
		return global.JSONResponseWithErrorV1(
			c,
			"401",
			"Invalid username/email or password",
			errors.New("user not found"),
			401,
		)
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

	// generate JWT token
	token, err := jwt.GenerateSuperAdminToken(
		found.ID,
		found.UserName,
	)

	if err != nil {
		return global.JSONResponseWithErrorV1(
			c,
			"500",
			"Failed to generate token",
			err,
			500,
		)
	}

	// decrypt email for response
	email, err := encrypDecryptV1.DecryptV1(found.Email, config.SecretKey)
	if err == nil {
		found.Email = email
	}

	found.Password = ""

	// response wrapper
	type LoginResponse struct {
		User  *SAdmodel.SuperAdminDetails `json:"user"`
		Token string                   `json:"token"`
	}

	return global.JSONResponseWithDataV1(
		c,
		"200",
		"Login successful",
		LoginResponse{
			User:  found,
			Token: token,
		},
		200,
	)
}


