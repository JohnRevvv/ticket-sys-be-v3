package controller

import (
	global "ideyanale-be/pkg/global/json_response"
	jwt "ideyanale-be/pkg/middleware/jwt"
	"ideyanale-be/pkg/modules/roles/model"
	RoleScript "ideyanale-be/pkg/modules/roles/script"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v3"
)

func EditRoleHandler(c fiber.Ctx) error {

	if err := jwt.RequireRoles(c, "Super-Admin", "Insti-Admin"); err != nil {
		return global.JSONResponseWithErrorV1(c, "403", "Forbidden", err, 403)
	}

	roleID, err := strconv.Atoi(c.Params("id"))
	if err != nil || roleID <= 0 {
		return global.JSONResponseWithErrorV1(c, "400", "Invalid role id", nil, 400)
	}

	var req model.Roles
	if err := c.Bind().Body(&req); err != nil {
		return global.JSONResponseWithErrorV1(c, "400", "Invalid request body", err, 400)
	}

	req.RoleName = strings.TrimSpace(req.RoleName)
	if req.RoleName == "" {
		return global.JSONResponseWithErrorV1(c, "400", "Role name is required", nil, 400)
	}

	existingRole, err := RoleScript.GetRoleByID(uint(roleID))
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "404", "Role not found", nil, 404)
	}

	currentRole, _ := c.Locals("role").(string)
	currentInstitution, _ := c.Locals("institution_id").(uint)

	if currentRole == "Insti-Admin" {
		if existingRole.InstitutionID != currentInstitution {
			return global.JSONResponseWithErrorV1(c, "403", "You can only edit roles in your institution", nil, 403)
		}
	}

	nameTaken, err := RoleScript.RoleNameExistsForInstitution(req.RoleName, existingRole.InstitutionID, uint(roleID))
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Validation failed", err, 500)
	}
	if nameTaken {
		return global.JSONResponseWithErrorV1(c, "409", "A role with this name already exists in this institution", nil, 409)
	}

	if err := RoleScript.EditRole(uint(roleID), existingRole.InstitutionID, req); err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Failed to update role", err, 500)
	}

	return global.JSONResponseWithDataV1(c, "200", "Role updated successfully", nil, 200)
}