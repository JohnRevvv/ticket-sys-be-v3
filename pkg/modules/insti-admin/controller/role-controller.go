package controller

import (
	global "ideyanale-be/pkg/global/json_response"
	jwt "ideyanale-be/pkg/middleware/jwt"
	InsAdScript "ideyanale-be/pkg/modules/insti-admin/script"
	"strings"

	"github.com/gofiber/fiber/v3"
)

func AddRole(c fiber.Ctx) error {

	if err := jwt.RequireRoles(c, "super-admin", "insti-admin"); err != nil {
		return global.JSONResponseWithErrorV1(c, "403", "Forbidden", err, 403)
	}

	type Req struct {
		RoleName         string `json:"role_name"`
		CanCreateTicket  bool   `json:"can_create"`
		CanEndorseTicket bool   `json:"can_endorse"`
		CanApproveTicket bool   `json:"can_approve"`
		CanResolveTicket bool   `json:"can_resolve"`
		CanAudit         bool   `json:"can_audit"`
	}

	var req Req

	if err := c.Bind().Body(&req); err != nil {
		return global.JSONResponseWithErrorV1(c, "400", "Invalid request", err, 400)
	}

	inst := c.Locals("institution_id")
	if inst == nil {
		return global.JSONResponseWithErrorV1(c, "401", "Unauthorized institution", nil, 401)
	}

	institutionID, ok := inst.(int)
	if !ok {
		return global.JSONResponseWithErrorV1(c, "500", "Invalid institution id type", nil, 500)
	}

	fields := strings.Fields(req.RoleName)
	normalizedRoleName := strings.Join(fields, " ")

	if normalizedRoleName == "" {
		return global.JSONResponseWithErrorV1(c, "400", "Role name is required", nil, 400)
	}

	existingRole, err := InsAdScript.ExistingRole(institutionID, normalizedRoleName,)
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Failed to check existing role", err, 500)
	}

	if existingRole {
		return global.JSONResponseWithErrorV1(c, "409", "Role name already exists", nil, 409)
	}

	err = InsAdScript.AddRole(
		normalizedRoleName,
		institutionID,
		req.CanCreateTicket,
		req.CanEndorseTicket,
		req.CanApproveTicket,
		req.CanResolveTicket,
		req.CanAudit,
	)
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Add role failed", err, 500)
	}

	return global.JSONResponseV1(c, "200", "Role added successfully", 200)
}
