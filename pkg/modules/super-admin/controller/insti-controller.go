package controller

import (
	"ideyanale-be/pkg/config"
	global "ideyanale-be/pkg/global/json_response"
	encrypDecryptV1 "ideyanale-be/pkg/middleware/encryption/v1"
	"ideyanale-be/pkg/modules/super-admin/script"
	"strings"

	"github.com/gofiber/fiber/v3"
)

func AddInstitution(c fiber.Ctx) error {

	type Req struct {
		InstitutionName string `json:"institution_name"`
		Description     string `json:"description"`
	}

	var req Req

	// validate request
	if err := c.Bind().Body(&req); err != nil {
		return global.JSONResponseWithErrorV1(c, "400", "Invalid request", err, 400)
	}

	institutionName := strings.TrimSpace(req.InstitutionName)
	description := strings.TrimSpace(req.Description)

	encinsitutionName, err := encrypDecryptV1.EncryptV1(institutionName, config.SecretKey)
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Encrypt institution name failed", err, 500)
	}

	encdescription, err := encrypDecryptV1.EncryptV1(description, config.SecretKey)
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Encrypt description failed", err, 500)
	}

	// save
	err = script.AddInstitution(
		encinsitutionName,
		encdescription,
	)

	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Add institution failed", err, 500)
	}

	return global.JSONResponseV1(c, "200", "Institution added successfully", 200)

}

func GetInstitutions(c fiber.Ctx) error {

	// fetch from script layer (DB only)
	rows, err := script.GetInstitutions()
	if err != nil {
		return global.JSONResponseWithErrorV1(
			c,
			"500",
			"Failed to fetch institutions",
			err,
			500,
		)
	}

	type InstitutionResp struct {
		InstitutionID   uint   `json:"institution_id"`
		InstitutionName string `json:"institution_name"`
		Description     string `json:"description"`
	}

	data := make([]InstitutionResp, 0, len(rows))

	for _, r := range rows {

		decryptedName, err := encrypDecryptV1.DecryptV1(r.InstitutionName, config.SecretKey)
		if err != nil {
			return global.JSONResponseWithErrorV1(
				c,
				"500",
				"Decrypt institution name failed",
				err,
				500,
			)
		}

		decryptedDesc, err := encrypDecryptV1.DecryptV1(r.Description, config.SecretKey)
		if err != nil {
			return global.JSONResponseWithErrorV1(
				c,
				"500",
				"Decrypt description failed",
				err,
				500,
			)
		}

		data = append(data, InstitutionResp{
			InstitutionID:   r.InstitutionID,
			InstitutionName: decryptedName,
			Description:     decryptedDesc,
		})
	}

	return global.JSONResponseWithDataV1(c, "200", "Institutions fetched successfully", data, 200)
}