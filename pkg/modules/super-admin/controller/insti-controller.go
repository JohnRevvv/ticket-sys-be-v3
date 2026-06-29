package controller

import (
	"ideyanale-be/pkg/config"
	global "ideyanale-be/pkg/global/json_response"
	encrypDecryptV1 "ideyanale-be/pkg/middleware/encryption/v1"
	"ideyanale-be/pkg/middleware/jwt"
	InsAdScript "ideyanale-be/pkg/modules/insti-admin/script"
	SAdScript "ideyanale-be/pkg/modules/super-admin/script"
	"strings"

	"github.com/gofiber/fiber/v3"
)

func AddInstitution(c fiber.Ctx) error {

	if err := jwt.RequireRoles(c, "Super-Admin"); err != nil {
		return global.JSONResponseWithErrorV1(c, "403", "Forbidden", err, 403)
	}

	type Req struct {
		InstitutionCode string `json:"institution_code"`
		InstitutionName string `json:"institution_name"`
		Description     string `json:"description"`
	}

	var req Req

	// validate request
	if err := c.Bind().Body(&req); err != nil {
		return global.JSONResponseWithErrorV1(c, "400", "Invalid request", err, 400)
	}

	institutionCode := strings.TrimSpace(req.InstitutionCode)
	institutionName := strings.TrimSpace(req.InstitutionName)
	description := strings.TrimSpace(req.Description)

	// normalize spaces
	codefields := strings.Fields(institutionCode)
	namefields := strings.Fields(institutionName)
	normalizedCode := strings.Join(codefields, " ")
	normalizedName := strings.Join(namefields, " ")

	encinsitutionCode, err := encrypDecryptV1.EncryptV1(normalizedCode, config.SecretKey)
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Encrypt institution code failed", err, 500)
	}

	encinsitutionName, err := encrypDecryptV1.EncryptV1(normalizedName, config.SecretKey)
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Encrypt institution name failed", err, 500)
	}

	encdescription, err := encrypDecryptV1.EncryptV1(description, config.SecretKey)
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Encrypt description failed", err, 500)
	}

	exists, err := SAdScript.InstitutionExists(encinsitutionCode, encinsitutionName)
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Validation failed", err, 500)
	}

	if exists {
		return global.JSONResponseWithErrorV1(c, "409", "institution already exists", nil, 409)
	}

	// save
	institutionID, err := SAdScript.AddInstitution(
		encinsitutionCode,
		encinsitutionName,
		encdescription,
	)
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Add institution failed", err, 500)
	}

	if err := InsAdScript.AddDefaultTicketTypes(uint(institutionID)); err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Failed to create default ticket types", err, 500)
	}

	if err := InsAdScript.AddDefaultCategories(uint(institutionID)); err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Failed to create default categories", err, 500)
	}

		if err := InsAdScript.AddDefaultSubCategories(uint(institutionID)); err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Failed to create default sub-categories", err, 500)
	}

	if err := InsAdScript.AddDefaultPositions(uint(institutionID)); err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Failed to create default positions", err, 500)
	}

	if err := InsAdScript.AddDefaultRoles(uint(institutionID)); err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Failed to create default roles", err, 500)
	}

	return global.JSONResponseV1(c, "200", "Institution added successfully", 200)
}

func GetInstitutions(c fiber.Ctx) error {

	// fetch from script layer (DB only)
	rows, err := SAdScript.GetInstitutions()
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Failed to fetch institutions", err, 500)
	}

	type InstitutionResp struct {
		InstitutionID   uint   `json:"institution_id"`
		InstitutionCode string `json:"institution_code"`
		InstitutionName string `json:"institution_name"`
		Description     string `json:"description"`
	}

	data := make([]InstitutionResp, 0, len(rows))

	for _, r := range rows {

		decryptedCode, err := encrypDecryptV1.DecryptV1(r.InstitutionCode, config.SecretKey)
		if err != nil {
			return global.JSONResponseWithErrorV1(c, "500", "Decrypt institution code failed", err, 500)
		}

		decryptedName, err := encrypDecryptV1.DecryptV1(r.InstitutionName, config.SecretKey)
		if err != nil {
			return global.JSONResponseWithErrorV1(c, "500", "Decrypt institution name failed", err, 500)
		}

		decryptedDesc, err := encrypDecryptV1.DecryptV1(r.Description, config.SecretKey)
		if err != nil {
			return global.JSONResponseWithErrorV1(c, "500", "Decrypt description failed", err, 500)
		}

		data = append(data, InstitutionResp{
			InstitutionID:   r.InstitutionID,
			InstitutionCode: decryptedCode,
			InstitutionName: decryptedName,
			Description:     decryptedDesc,
		})
	}

	return global.JSONResponseWithDataV1(c, "200", "Institutions fetched successfully", data, 200)
}
