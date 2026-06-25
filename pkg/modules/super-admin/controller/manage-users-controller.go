package controller

import (
	"errors"
	global "ideyanale-be/pkg/global/json_response"
	"ideyanale-be/pkg/middleware/jwt"
	"ideyanale-be/pkg/modules/super-admin/script"
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

	if err := script.ChangeRoleToAdmin(userID); err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Failed to change user role", err, 500)
	}

	return global.JSONResponseWithDataV1(c, "200", "User role changed to admin successfully", nil, 200)
}

func ChangeUserStatus(c fiber.Ctx) error {
	if err := jwt.RequireRoles(c, "super-admin", "insti-admin"); err != nil {
		return global.JSONResponseWithErrorV1(c, "403", "Forbidden", err, 403)
	}

	userID, err := strconv.Atoi(c.Params("id"))
	if err != nil || userID <= 0 {
		return global.JSONResponseWithErrorV1(c, "400", "Invalid user id", err, 400,)
	}

	type Req struct {
		Status string `json:"status"`
	}

	var req Req

	if err := c.Bind().Body(&req); err != nil {
		return global.JSONResponseWithErrorV1(c, "400", "Invalid request", err, 400,)
	}

	status := strings.TrimSpace(strings.ToLower(req.Status))

	switch status {
	case "active", "disabled":
		// valid
	default:
		return global.JSONResponseWithErrorV1(
			c, "400", "Invalid status", errors.New("status must be active or disabled"), 400,)
	}

	if err := script.ChangeUserStatus(userID, status); err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Failed to change user status", err,500,)
	}

	return global.JSONResponseWithDataV1(c, "200", "User status updated successfully", nil, 200,)
}
