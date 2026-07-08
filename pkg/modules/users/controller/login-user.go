package controller

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"ideyanale-be/pkg/config"
	global "ideyanale-be/pkg/global/json_response"
	middleware "ideyanale-be/pkg/middleware/autologout"
	encrypDecryptV1 "ideyanale-be/pkg/middleware/encryption/v1"
	"ideyanale-be/pkg/middleware/jwt"
	"ideyanale-be/pkg/modules/users/model"
	"ideyanale-be/pkg/modules/users/script"
	"ideyanale-be/pkg/services/email"

	"github.com/gofiber/fiber/v3"
	"golang.org/x/crypto/bcrypt"
)

type LoginReq struct {
	StaffID string `json:"staff_id"`
}

type VerifyOTPReq struct {
	StaffID string `json:"staff_id"`
	OTP     string `json:"otp"`
}

func generateOTP() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}

func LoginWithOTP(c fiber.Ctx) error {

	var req LoginReq

	if err := c.Bind().Body(&req); err != nil {
		return global.JSONResponseWithErrorV1(c, "400", "Invalid request", err, 400)
	}

	staffID := strings.TrimSpace(req.StaffID)
	staffID = strings.ReplaceAll(staffID, "-", "")

	if len(staffID) != 11 {
		return global.JSONResponseWithErrorV1(c, "400", "Invalid staff ID", nil, 400)
	}

	if _, err := strconv.Atoi(staffID); err != nil {
		return global.JSONResponseWithErrorV1(c, "400", "Staff ID must be numeric", nil, 400)
	}

	formattedStaffID := staffID[:6] + "-" + staffID[6:]

	// encrypt (same style as register)
	encStaffID, err := encrypDecryptV1.EncryptV1(formattedStaffID, config.SecretKey)
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Encryption failed", err, 500)
	}

	// -------------------------
	// CHECK ACTIVE USER FIRST
	// -------------------------
	user, err := script.GetActiveUserByStaffID(encStaffID)
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "404", "User not found or inactive", err, 404)
	}

	// -------------------------
	// CHECK IF ALREADY LOGGED IN
	// -------------------------
	if user.IsLoggedIn {
		return global.JSONResponseWithErrorV1(
			c,
			"409",
			"User is already logged in on another device",
			nil,
			409,
		)
	}

	// -------------------------
	// DELETE OLD OTP (OPTION 1 FIX)
	// -------------------------
	_ = script.DeleteOTPByStaffID(encStaffID)

	// -------------------------
	// GENERATE OTP
	// -------------------------
	otp := generateOTP()

	log.Printf("[LOGIN OTP] StaffID: %s | OTP: %s", formattedStaffID, otp)

	otpHash, err := bcrypt.GenerateFromPassword([]byte(otp), bcrypt.DefaultCost)
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "OTP hashing failed", err, 500)
	}

	// save OTP
	err = script.SaveLoginOTP(&model.LoginOTP{
		StaffID:   encStaffID,
		OTPHash:   string(otpHash),
		ExpiresAt: time.Now().Add(5 * time.Minute).Unix(),
	})

	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Failed to store OTP", err, 500)
	}

	// -------------------------
	// SEND EMAIL (service layer recommended)
	// -------------------------
	decEmail, err := encrypDecryptV1.DecryptV1(user.Email, config.SecretKey)
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Decrypt email failed", err, 500)
	}

	go func(emailAddr, otp string) {
		_ = email.SendOTP(emailAddr, otp)
	}(decEmail, otp)

	return global.JSONResponseV1(c, "200", "OTP sent successfully", 200)
}

func VerifyOTP(c fiber.Ctx) error {

	var req VerifyOTPReq

	if err := c.Bind().Body(&req); err != nil {
		return global.JSONResponseWithErrorV1(c, "400", "Invalid request", err, 400)
	}

	staffID := strings.TrimSpace(strings.ReplaceAll(req.StaffID, "-", ""))

	if len(staffID) != 11 {
		return global.JSONResponseWithErrorV1(c, "400", "Invalid staff ID", nil, 400)
	}

	formatted := staffID[:6] + "-" + staffID[6:]

	encStaffID, err := encrypDecryptV1.EncryptV1(formatted, config.SecretKey)
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Encryption failed", err, 500)
	}

	// get OTP record
	record, err := script.GetOTPByStaffID(encStaffID)
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "404", "OTP not found", err, 404)
	}

	// check expiry
	if time.Now().Unix() > record.ExpiresAt {
		return global.JSONResponseWithErrorV1(c, "400", "OTP expired", nil, 400)
	}

	// validate OTP
	if err := bcrypt.CompareHashAndPassword([]byte(record.OTPHash), []byte(req.OTP)); err != nil {
		return global.JSONResponseWithErrorV1(c, "400", "Invalid OTP", nil, 400)
	}

	// get user for JWT — NOTE: must preload Role, see script change below
	user, err := script.GetActiveUserByStaffID(encStaffID)
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "404", "User not found", err, 404)
	}

	// NEW: guard against a user whose role record is missing/deleted
	if user.Role.RoleID == 0 {
		return global.JSONResponseWithErrorV1(c, "403", "No valid role assigned to this user", nil, 403)
	}

	// generate JWT — now passes RoleID + permission flags, not the whole struct
	token, err := jwt.GenerateUserToken(
		user.ID,
		user.StaffID,
		user.InstitutionID,
		user.Role.RoleID,
		user.Role.RoleName,
		jwt.Permissions{
			CanCreateTicket:  user.Role.CanCreateTicket,
			CanEndorseTicket: user.Role.CanEndorseTicket,
			CanApproveTicket: user.Role.CanApproveTicket,
			CanResolveTicket: user.Role.CanResolveTicket,
			CanAudit:         user.Role.CanAudit,
		},
	)

	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Failed to generate token", err, 500)
	}

	// update login status
	err = script.SetUserLoginStatus(user.ID, true)
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Failed to update login status", err, 500)
	}

	// cleanup OTP
	_ = script.DeleteOTPByStaffID(encStaffID)

	// FIXED: TouchActivity expects (int, string) — pass RoleName, not RoleID
	middleware.TouchActivity(user.ID, user.Role.RoleName)

	// response
	return global.JSONResponseWithDataV1(
		c,
		"200",
		"Login successful",
		map[string]any{
			"token": token,
		},
		200,
	)
}

// func VerifyOTP(c fiber.Ctx) error {

// 	var req VerifyOTPReq

// 	if err := c.Bind().Body(&req); err != nil {
// 		return global.JSONResponseWithErrorV1(c, "400", "Invalid request", err, 400)
// 	}

// 	staffID := strings.TrimSpace(strings.ReplaceAll(req.StaffID, "-", ""))

// 	if len(staffID) != 11 {
// 		return global.JSONResponseWithErrorV1(c, "400", "Invalid staff ID", nil, 400)
// 	}

// 	formatted := staffID[:6] + "-" + staffID[6:]

// 	encStaffID, err := encrypDecryptV1.EncryptV1(formatted, config.SecretKey)
// 	if err != nil {
// 		return global.JSONResponseWithErrorV1(c, "500", "Encryption failed", err, 500)
// 	}

// 	// get OTP record
// 	record, err := script.GetOTPByStaffID(encStaffID)
// 	if err != nil {
// 		return global.JSONResponseWithErrorV1(c, "404", "OTP not found", err, 404)
// 	}

// 	// check expiry
// 	if time.Now().Unix() > record.ExpiresAt {
// 		return global.JSONResponseWithErrorV1(c, "400", "OTP expired", nil, 400)
// 	}

// 	// validate OTP
// 	if err := bcrypt.CompareHashAndPassword([]byte(record.OTPHash), []byte(req.OTP)); err != nil {
// 		return global.JSONResponseWithErrorV1(c, "400", "Invalid OTP", nil, 400)
// 	}

// 	// get user for JWT
// 	user, err := script.GetActiveUserByStaffID(encStaffID)
// 	if err != nil {
// 		return global.JSONResponseWithErrorV1(c, "404", "User not found", err, 404)
// 	}

// 	// generate JWT
// 	token, err := jwt.GenerateUserToken(
// 		user.ID,
// 		user.StaffID,
// 		user.InstitutionID,
// 		user.Role,
// 	)

// 	if err != nil {
// 		return global.JSONResponseWithErrorV1(c, "500", "Failed to generate token", err, 500)
// 	}

// 	// update login status
// 	err = script.SetUserLoginStatus(user.ID, true)
// 	if err != nil {
// 		return global.JSONResponseWithErrorV1(c, "500", "Failed to update login status", err, 500)
// 	}

// 	// cleanup OTP
// 	_ = script.DeleteOTPByStaffID(encStaffID)

// 	middleware.TouchActivity(user.ID, user.Role)

// 	// response
// 	return global.JSONResponseWithDataV1(
// 		c,
// 		"200",
// 		"Login successful",
// 		map[string]any{
// 			"token": token,
// 		},
// 		200,
// 	)
// }
