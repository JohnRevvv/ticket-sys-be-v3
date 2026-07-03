package controller

import (
	"ideyanale-be/pkg/config"
	global "ideyanale-be/pkg/global/json_response"
	encrypDecryptV1 "ideyanale-be/pkg/middleware/encryption/v1"
	"ideyanale-be/pkg/middleware/jwt"
	InsAdScript "ideyanale-be/pkg/modules/insti-admin/script"
	InstiScript "ideyanale-be/pkg/modules/institution/script"
	services "ideyanale-be/pkg/services/s3_service"
	"regexp"
	"strconv"
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

	exists, err := InstiScript.InstitutionExists(encinsitutionCode, encinsitutionName)
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Validation failed", err, 500)
	}

	if exists {
		return global.JSONResponseWithErrorV1(c, "409", "institution already exists", nil, 409)
	}

	// save
	institutionID, err := InstiScript.AddInstitution(
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
	rows, err := InstiScript.GetInstitutions()
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

func EditInstitution(c fiber.Ctx) error {
	if err := jwt.RequireRoles(c, "Super-Admin", "Insti-Admin"); err != nil {
		return global.JSONResponseWithErrorV1(c, "403", "Forbidden", err, 403)
	}

	userID, ok := c.Locals("id").(int)
	if !ok {
		return global.JSONResponseWithErrorV1(c, "401", "Invalid user ID", nil, 401)
	}

	institutionID, err := strconv.Atoi(c.Params("institution_id"))
	if err != nil || institutionID <= 0 {
		return global.JSONResponseWithErrorV1(c, "400", "Invalid institution_id", err, 400)
	}

	institutionCode := strings.TrimSpace(c.FormValue("institution_code"))
	institutionName := strings.TrimSpace(c.FormValue("institution_name"))
	description := strings.TrimSpace(c.FormValue("description"))
	institutionColor := strings.TrimSpace(c.FormValue("institution_color"))

	// Validate HEX color
	hexRegex := regexp.MustCompile(`^#(?:[0-9a-fA-F]{3}|[0-9a-fA-F]{6})$`)
	if institutionColor != "" && !hexRegex.MatchString(institutionColor) {
		return global.JSONResponseWithErrorV1(c, "400", "Invalid institution color", nil, 400)
	}

	institutionCode = strings.Join(strings.Fields(institutionCode), " ")
	institutionName = strings.Join(strings.Fields(institutionName), " ")

	encCode, err := encrypDecryptV1.EncryptV1(institutionCode, config.SecretKey)
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Encrypt institution code failed", err, 500)
	}

	encName, err := encrypDecryptV1.EncryptV1(institutionName, config.SecretKey)
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Encrypt institution name failed", err, 500)
	}

	encDescription, err := encrypDecryptV1.EncryptV1(description, config.SecretKey)
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Encrypt description failed", err, 500)
	}

	if err := InstiScript.UpdateInstitution(uint(institutionID), encCode, encName, encDescription, institutionColor); err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Update institution failed", err, 500)
	}


	var s3Service, _ = services.NewS3Service()
	// Optional logo upload
	fileHeader, err := c.FormFile("logo")
	if err == nil {

		oldLogo, _ := InstiScript.GetInstitutionLogo(uint(institutionID))

		fileName, fileKey, err := s3Service.UploadLogo(fileHeader)
		if err != nil {
			return global.JSONResponseWithErrorV1(c, "500", err.Error(), err, 500)
		}

		if oldLogo != nil && oldLogo.FileKey != "" {
			_ = s3Service.Delete(oldLogo.FileKey)
		}

		if err := InstiScript.UpsertInstitutionLogo(uint(institutionID), fileName, fileKey, uint(userID)); err != nil {
			return global.JSONResponseWithErrorV1(c, "500", "Failed to save logo", err, 500)
		}
	}

	return global.JSONResponseV1(c, "200", "Institution updated successfully", 200)
}
