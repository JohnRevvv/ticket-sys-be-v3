package controller

import (
	"strconv"

	"github.com/gofiber/fiber/v3"

	global "ideyanale-be/pkg/global/json_response"
	RoleScript "ideyanale-be/pkg/modules/roles/script"
)

func GetRolesByInstitutionHandler(c fiber.Ctx) error {

	institutionID, err := strconv.Atoi(c.Params("institution_id"))
	if err != nil || institutionID <= 0 {
		return global.JSONResponseWithErrorV1(c, "400", "Invalid institution id", nil, 400)
	}

	roles, err := RoleScript.GetRolesByInstitutionID(uint(institutionID))
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Failed to fetch roles", err, 500)
	}

	if len(roles) == 0 {
		return global.JSONResponseWithErrorV1(c, "404", "No roles found for this institution", nil, 404)
	}

	return global.JSONResponseWithDataV1(c, "200", "Roles fetched successfully", roles, 200)
}
