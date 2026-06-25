package controller

import (
	"strings"

	"github.com/gofiber/fiber/v3"

	"ideyanale-be/pkg/config"
	global "ideyanale-be/pkg/global/json_response"
	encrypDecryptV1 "ideyanale-be/pkg/middleware/encryption/v1"
	jwt "ideyanale-be/pkg/middleware/jwt"
	InsAdScript "ideyanale-be/pkg/modules/insti-admin/script"
)



func AddPosition(c fiber.Ctx) error {

	if err := jwt.RequireRoles(c, "super-admin", "insti-admin"); err != nil {
		return global.JSONResponseWithErrorV1(
			c,
			"403",
			"Forbidden",
			err,
			403,
		)
	}

	type Req struct {
		PositionName string `json:"position_name"`
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

	// Fetch existing positions for this institution
	positions, err := InsAdScript.GetPositionsByInstitutionID(institutionID)
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Fetch existing positions failed", err, 500)
	}

	// Decrypt and compare in the controller
	trimmedNewName := strings.TrimSpace(req.PositionName)
	for _, pos := range positions {
		decName, err := encrypDecryptV1.DecryptV1(pos.JobName, config.SecretKey)
		if err != nil {
			return global.JSONResponseWithErrorV1(c, "500", "Decrypt position name failed", err, 500)
		}

		if strings.EqualFold(strings.TrimSpace(decName), trimmedNewName) {
			return global.JSONResponseWithErrorV1(c, "409", "Position name already exists", nil, 409)
		}
	}

	encPositionName, err := encrypDecryptV1.EncryptV1(req.PositionName, config.SecretKey)
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Encrypt position name failed", err, 500)
	}

	err = InsAdScript.AddPosition(encPositionName, institutionID)
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Add position failed", err, 500)
	}

	return global.JSONResponseV1(c, "200", "Position added successfully", 200)
}

func GetPositionsByInstitutionID(c fiber.Ctx) error {

	inst := c.Locals("institution_id")
	if inst == nil {
		return global.JSONResponseWithErrorV1(c, "401", "Unauthorized institution", nil, 401)
	}

	institutionID, ok := inst.(int)
	if !ok {
		return global.JSONResponseWithErrorV1(c, "500", "Invalid institution id type", nil, 500)
	}

	rows, err := InsAdScript.GetPositionsByInstitutionID(institutionID)
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Failed to fetch positions", err, 500)
	}

	type PositionResp struct {
		PositionID   uint   `json:"position_id"`
		PositionName string `json:"position_name"`
	}

	var resp []PositionResp
	for _, row := range rows {
		decName, err := encrypDecryptV1.DecryptV1(row.JobName, config.SecretKey)
		if err != nil {
			return global.JSONResponseWithErrorV1(c, "500", "Decrypt position name failed", err, 500)
		}

		resp = append(resp, PositionResp{
			PositionID:   row.PositionID,
			PositionName: decName,
		})
	}

	return global.JSONResponseWithDataV1(c, "200", "Positions fetched successfully", resp, 200)
}

func EditPosition(c fiber.Ctx) error {

	if err := jwt.RequireRoles(c, "super-admin", "insti-admin"); err != nil {
		return global.JSONResponseWithErrorV1(c, "403", "Forbidden", err, 403,)
	}

	type Req struct {
		PositionID   int    `json:"position_id"`
		PositionName string `json:"position_name"`
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

	// Fetch existing positions for this institution
	positions, err := InsAdScript.GetPositionsByInstitutionID(institutionID)
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Fetch existing positions failed", err, 500)
	}

	// Decrypt and compare in the controller, skipping the position being edited
	trimmedNewName := strings.TrimSpace(req.PositionName)
	for _, pos := range positions {
		if pos.PositionID == uint(req.PositionID) {
			continue
		}

		decName, err := encrypDecryptV1.DecryptV1(pos.JobName, config.SecretKey)
		if err != nil {
			return global.JSONResponseWithErrorV1(c, "500", "Decrypt position name failed", err, 500)
		}

		if strings.EqualFold(strings.TrimSpace(decName), trimmedNewName) {
			return global.JSONResponseWithErrorV1(c, "409", "Position name already exists", nil, 409)
		}
	}

	encPositionName, err := encrypDecryptV1.EncryptV1(req.PositionName, config.SecretKey)
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Encrypt position name failed", err, 500)
	}

	rowsAffected, err := InsAdScript.UpdatePosition(req.PositionID, institutionID, encPositionName)
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Update position failed", err, 500)
	}

	if rowsAffected == 0 {
		return global.JSONResponseWithErrorV1(c, "404", "Position not found", nil, 404)
	}

	return global.JSONResponseV1(c, "200", "Position updated successfully", 200)
}