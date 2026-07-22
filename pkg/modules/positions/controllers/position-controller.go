package controller

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v3"

	global "ideyanale-be/pkg/global/json_response"
	jwt "ideyanale-be/pkg/middleware/jwt"
	PositionScript "ideyanale-be/pkg/modules/positions/scripts"
)

func AddPosition(c fiber.Ctx) error {

	if err := jwt.RequireRoles(c, "Super-Admin", "Insti-Admin"); err != nil {
		return global.JSONResponseWithErrorV1(c, "403", "Forbidden", err, 403)
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

		fmt.Printf("institution_id value=%v type=%T\n", inst, inst)


	institutionID, ok := inst.(uint)
	if !ok {
		return global.JSONResponseWithErrorV1(c, "500", "Invalid institution id type", nil, 500)
	}
	trimmedNewName := strings.TrimSpace(req.PositionName)

	// Fetch existing positions for this institution
	positions, err := PositionScript.GetPositionsByInstitutionID(uint(institutionID))
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Fetch existing positions failed", err, 500)
	}

	// Compare directly
	for _, pos := range positions {
		if strings.EqualFold(strings.TrimSpace(pos.PositionName), trimmedNewName) {
			return global.JSONResponseWithErrorV1(c, "409", "Position name already exists", nil, 409)
		}
	}

	err = PositionScript.AddPosition(trimmedNewName, uint(institutionID))
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Add position failed", err, 500)
	}

	return global.JSONResponseV1(c, "200", "Position added successfully", 200)
}

func GetPositions(c fiber.Ctx) error {

	inst := c.Locals("institution_id")
	if inst == nil {
		return global.JSONResponseWithErrorV1(c, "401", "Unauthorized institution", nil, 401)
	}

	institutionID, ok := inst.(int)
	if !ok {
		return global.JSONResponseWithErrorV1(c, "500", "Invalid institution id type", nil, 500)
	}

	rows, err := PositionScript.GetPositionsByInstitutionID(uint(institutionID))
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Failed to fetch positions", err, 500)
	}

	type PositionResp struct {
		PositionID   uint   `json:"position_id"`
		PositionName string `json:"position_name"`
	}

	var resp []PositionResp
	for _, row := range rows {
		resp = append(resp, PositionResp{
			PositionID:   row.PositionID,
			PositionName: row.PositionName,
		})
	}

	return global.JSONResponseWithDataV1(c, "200", "Positions fetched successfully", resp, 200)
}

func GetPositionsByInstitutionID(c fiber.Ctx) error {

	institutionID, err := strconv.Atoi(c.Params("institution_id"))
	if err != nil || institutionID <= 0 {
		return global.JSONResponseWithErrorV1(c, "400", "Invalid institution id", nil, 400)
	}

	rows, err := PositionScript.GetPositionsByInstitutionID(uint(institutionID))
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Failed to fetch positions", err, 500)
	}

	type PositionResp struct {
		PositionID   uint   `json:"position_id"`
		PositionName string `json:"position_name"`
	}

	var resp []PositionResp
	for _, row := range rows {
		resp = append(resp, PositionResp{
			PositionID:   row.PositionID,
			PositionName: row.PositionName,
		})
	}

	return global.JSONResponseWithDataV1(c, "200", "Positions fetched successfully", resp, 200)
}

func EditPosition(c fiber.Ctx) error {

	if err := jwt.RequireRoles(c, "Super-Admin", "Insti-Admin"); err != nil {
		return global.JSONResponseWithErrorV1(c, "403", "Forbidden", err, 403)
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

	positions, err := PositionScript.GetPositionsByInstitutionID(uint(institutionID))
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Fetch existing positions failed", err, 500)
	}

	trimmedNewName := strings.TrimSpace(req.PositionName)

	for _, pos := range positions {
		if pos.PositionID == uint(req.PositionID) {
			continue
		}

		if strings.EqualFold(strings.TrimSpace(pos.PositionName), trimmedNewName) {
			return global.JSONResponseWithErrorV1(c, "409", "Position name already exists", nil, 409)
		}
	}

	rowsAffected, err := PositionScript.UpdatePosition(
		req.PositionID,
		uint(institutionID),
		trimmedNewName,
	)
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Update position failed", err, 500)
	}

	if rowsAffected == 0 {
		return global.JSONResponseWithErrorV1(c, "404", "Position not found", nil, 404)
	}

	return global.JSONResponseV1(c, "200", "Position updated successfully", 200)
}