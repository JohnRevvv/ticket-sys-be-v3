package controller

import (
	"ideyanale-be/pkg/config"
	global "ideyanale-be/pkg/global/json_response"
	encrypDecryptV1 "ideyanale-be/pkg/middleware/encryption/v1"
	"ideyanale-be/pkg/middleware/jwt"
	"ideyanale-be/pkg/modules/projects/script"
	"strings"

	"github.com/gofiber/fiber/v3"
)

func AddServer(c fiber.Ctx) error {

	if err := jwt.RequireRoles(c, "Insti-Admin"); err != nil {
		return global.JSONResponseWithErrorV1(c, "403", "Forbidden", err, 403)
	}

	inst := c.Locals("institution_id")
	if inst == nil {
		return global.JSONResponseWithErrorV1(c, "401", "Unauthorized institution", nil, 401)
	}

	institutionID, ok := inst.(int)
	if !ok {
		return global.JSONResponseWithErrorV1(c, "500", "Invalid institution id type", nil, 500)
	}

	type ServerReq struct {
		ServerName string `gorm:"column:server_name;not null" json:"server_name"`
		ServerIP   string `gorm:"column:server_ip;not null" json:"server_ip"`
	}

	var req ServerReq
	if err := c.Bind().Body(&req); err != nil {
		return global.JSONResponseWithErrorV1(c, "400", "Invalid request", err, 400)
	}

	trimmedServerName := strings.TrimSpace(req.ServerName)
	trimmedServerIP := strings.TrimSpace(req.ServerIP)

	if trimmedServerName == "" {
		return global.JSONResponseWithErrorV1(c, "400", "Server name is required", nil, 400)
	}

	if trimmedServerIP == "" {
		return global.JSONResponseWithErrorV1(c, "400", "Server IP is required", nil, 400)
	}

	// normalize spaces
	servernamefields := strings.Fields(trimmedServerName)
	normalizedServerName := strings.Join(servernamefields, " ")
	normalizedServerIP := strings.Join(strings.Fields(trimmedServerIP), "")

	encServerName, err := encrypDecryptV1.EncryptV1(normalizedServerName, config.SecretKey)
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Encrypt server name failed", err, 500)
	}

	encServerIP, err := encrypDecryptV1.EncryptV1(normalizedServerIP, config.SecretKey)
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Encrypt server IP failed", err, 500)
	}

	serverexists, err := script.ServerExists(uint(institutionID), encServerName, encServerIP)
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Validation failed", err, 500)
	}

	if serverexists {
		return global.JSONResponseWithErrorV1(c, "409", "Server already exists", nil, 409)
	}

	if err := script.AddServerScript(encServerName, encServerIP, institutionID); err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Add server failed", err, 500)
	}

	return global.JSONResponseV1(c, "200", "Server added successfully", 200)
}

func AddProject(c fiber.Ctx) error {

	if err := jwt.RequireRoles(c, "Insti-Admin"); err != nil {
		return global.JSONResponseWithErrorV1(c, "403", "Forbidden", err, 403)
	}

	inst := c.Locals("institution_id")
	if inst == nil {
		return global.JSONResponseWithErrorV1(c, "401", "Unauthorized institution", nil, 401)
	}

	type ProjectReq struct {
		ServerID    uint   `json:"server_id"`
		ProjectName string `json:"project_name"`
		Environment string `json:"environment"`
		Description string `json:"description"`
	}

	var req ProjectReq

	if err := c.Bind().Body(&req); err != nil {
		return global.JSONResponseWithErrorV1(c, "400", "Invalid request", err, 400)
	}

	trimmedProjectName := strings.TrimSpace(req.ProjectName)
	trimmedEnvironment := strings.TrimSpace(req.Environment)

	if req.ServerID == 0 {
		return global.JSONResponseWithErrorV1(c, "400", "Server id is required", nil, 400)
	}

	if trimmedProjectName == "" {
		return global.JSONResponseWithErrorV1(c, "400", "project name is required", nil, 400)
	}
	if trimmedEnvironment == "" {
		return global.JSONResponseWithErrorV1(c, "400", "environment is required", nil, 400)
	}

	encProjectName, err := encrypDecryptV1.EncryptV1(trimmedProjectName, config.SecretKey)
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Encrypt project name failed", err, 500)
	}

	encEnvironment, err := encrypDecryptV1.EncryptV1(trimmedEnvironment, config.SecretKey)
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Encrypt environment failed", err, 500)
	}
	encDescription, err := encrypDecryptV1.EncryptV1(req.Description, config.SecretKey)
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Encrypt description failed", err, 500)
	}

	server, err := script.GetServerByServerID(int(req.ServerID))
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Fetch server failed", err, 500)
	}
	if server.ServerID == 0 {
		return global.JSONResponseWithErrorV1(c, "404", "Server not found", nil, 404)
	}

	projectexists, err := script.ProjectExists(uint(server.ServerID), encProjectName, encEnvironment)
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Validation failed", err, 500)
	}

	if projectexists {
		return global.JSONResponseWithErrorV1(c, "409", "Project already exists", nil, 409)
	}

	if err := script.AddProjectScript(

		encProjectName,
		encEnvironment,
		encDescription,
		int(server.ServerID),
	); err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Failed to add project", err, 500)
	}

	return global.JSONResponseV1(c, "200", "Project added successfully", 200)

}
