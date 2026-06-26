package controller

import (
	"errors"
	"fmt"
	global "ideyanale-be/pkg/global/json_response"
	"ideyanale-be/pkg/middleware/jwt"
	SupAdScript "ideyanale-be/pkg/modules/super-admin/script"
	UserScript "ideyanale-be/pkg/modules/users/script"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v3"
)

func ChangeRoleToAdmin(c fiber.Ctx) error {

	if err := jwt.RequireRoles(c, "super-admin"); err != nil {
		return global.JSONResponseWithErrorV1(c, "403", "Forbidden", err, 403)
	}

	userID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "400", "Invalid ticket_type_id", err, 400)
	}

	if err := SupAdScript.ChangeRoleToAdmin(userID); err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Failed to change user role", err, 500)
	}

	return global.JSONResponseWithDataV1(c, "200", "User role changed to admin successfully", nil, 200)
}

func ChangeUserRole(c fiber.Ctx) error {

	if err := jwt.RequireRoles(c, "super-admin", "insti-admin"); err != nil {
		return global.JSONResponseWithErrorV1(c, "403", "Forbidden", err, 403)
	}

	userID, err := strconv.Atoi(c.Params("id"))
	if err != nil || userID <= 0 {
		return global.JSONResponseWithErrorV1(c, "400", "Invalid user id", nil, 400)
	}

	type Req struct {
		RoleID uint `json:"role_id"`
	}

	var req Req

	if err := c.Bind().Body(&req); err != nil {
		return global.JSONResponseWithErrorV1(c, "400", "Invalid request body", err, 400)
	}

	if req.RoleID == 0 {
		return global.JSONResponseWithErrorV1(c, "400", "Role ID is required", nil, 400)
	}

	// Check if role exists
	role, err := SupAdScript.GetRoleByID(req.RoleID)
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Failed to fetch role", err, 500)
	}

	if role.RoleID == 0 {
		return global.JSONResponseWithErrorV1(c, "404", "Role not found", nil, 404)
	}

	// Get current user's role
	currentRole, _ := c.Locals("role").(string)

	// Get current user's institution
	currentInstitution, _ := c.Locals("institution_id").(int)

	// SUPER ADMIN
	if currentRole == "super-admin" {

		err = SupAdScript.ChangeUserRole(userID, req.RoleID)
		if err != nil {
			return global.JSONResponseWithErrorV1(c, "500", "Failed to update role", err, 500)
		}

		return global.JSONResponseWithDataV1(c, "200", "User role updated successfully", nil, 200)
	}

	// INSTI ADMIN
	if currentRole == "insti-admin" {

		// role must belong to same institution
		if int(role.InstitutionID) != currentInstitution {
			return global.JSONResponseWithErrorV1(c, "403", "You can only assign roles from your institution",
				nil,
				403,
			)
		}

		user, err := UserScript.GetUserByID(userID)
		if err != nil {
			return global.JSONResponseWithErrorV1(
				c,
				"500",
				"Failed to fetch user",
				err,
				500,
			)
		}

		if int(user.InstitutionID) != currentInstitution {
			return global.JSONResponseWithErrorV1(
				c,
				"403",
				"You can only manage users in your institution",
				nil,
				403,
			)
		}

		err = SupAdScript.ChangeUserRole(userID, req.RoleID)
		if err != nil {
			return global.JSONResponseWithErrorV1(
				c,
				"500",
				"Failed to update role",
				err,
				500,
			)
		}

		return global.JSONResponseWithDataV1(
			c,
			"200",
			"User role updated successfully",
			nil,
			200,
		)
	}

	return global.JSONResponseWithErrorV1(
		c,
		"403",
		"Forbidden",
		nil,
		403,
	)
}

func ChangeUserStatus(c fiber.Ctx) error {
	if err := jwt.RequireRoles(c, "super-admin", "insti-admin"); err != nil {
		return global.JSONResponseWithErrorV1(c, "403", "Forbidden", err, 403)
	}

	userID, err := strconv.Atoi(c.Params("id"))
	if err != nil || userID <= 0 {
		return global.JSONResponseWithErrorV1(c, "400", "Invalid user id", err, 400)
	}

	type Req struct {
		Status string `json:"status"`
	}

	var req Req

	if err := c.Bind().Body(&req); err != nil {
		return global.JSONResponseWithErrorV1(c, "400", "Invalid request", err, 400)
	}

	status := strings.TrimSpace(strings.ToLower(req.Status))

	switch status {
	case "pending", "active", "disabled":
		// valid
	default:
		return global.JSONResponseWithErrorV1(
			c, "400", "Invalid status", errors.New("status must be active or disabled"), 400)
	}

	fmt.Println("User ID:", userID)
	fmt.Println("Status:", status)

	if err := SupAdScript.ChangeUserStatus(userID, status); err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Failed to change user status", err, 500)
	}

	return global.JSONResponseWithDataV1(c, "200", "User status updated successfully", nil, 200)
}
