package controller

import (
	"ideyanale-be/pkg/config"
	global "ideyanale-be/pkg/global/json_response"
	encrypDecryptV1 "ideyanale-be/pkg/middleware/encryption/v1"
	"ideyanale-be/pkg/middleware/jwt"
	"ideyanale-be/pkg/modules/projects/script"
	"strconv"
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

	// institutionID, ok := inst.(int)
	// if !ok {
	// 	return global.JSONResponseWithErrorV1(c, "500", "Invalid institution id type", nil, 500)
	// }

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

func GetServerByID(c fiber.Ctx) error {


	inst := c.Locals("institution_id")
	if inst == nil {
		return global.JSONResponseWithErrorV1(c, "401", "Unauthorized institution", nil, 401)
	}

	institutionID, ok := inst.(int)
	if !ok {
		return global.JSONResponseWithErrorV1(c, "500", "Invalid institution id type", nil, 500)
	}

	serverID, err := strconv.Atoi(c.Params("server_id"))
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "400", "Invalid server_id", err, 400)
	}

	// SCRIPT CALL (DB ONLY)
	server, err := script.GetServerByServerID(serverID)
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Failed to get server", err, 500)
	}

	if server.ServerID == 0 {
		return global.JSONResponseWithErrorV1(c, "404", "Server not found", nil, 404)
	}

	// OWNERSHIP CHECK
	if server.InstitutionID != uint(institutionID) {
		return global.JSONResponseWithErrorV1(c, "403", "Forbidden", nil, 403)
	}

	// DECRYPT
	decryptedServerName, err := encrypDecryptV1.DecryptV1(server.ServerName, config.SecretKey)
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Decrypt server name failed", err, 500)
	}

	decryptedServerIP, err := encrypDecryptV1.DecryptV1(server.ServerIP, config.SecretKey)
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Decrypt server IP failed", err, 500)
	}

	response := struct {
		ServerID      uint   `json:"server_id"`
		InstitutionID int    `json:"institution_id"`
		ServerName    string `json:"server_name"`
		ServerIP      string `json:"server_ip"`
		Status        string `json:"status"`
	}{
		ServerID:      uint(server.ServerID),
		InstitutionID: int(server.InstitutionID),
		ServerName:    decryptedServerName,
		ServerIP:      decryptedServerIP,
		Status:        server.Status,
	}

	return global.JSONResponseWithDataV1(
		c,
		"200",
		"Server retrieved successfully",
		response,
		200,
	)
}

func GetProjectByID(c fiber.Ctx) error {
	inst := c.Locals("institution_id")
	if inst == nil {
		return global.JSONResponseWithErrorV1(c, "401", "Unauthorized institution", nil, 401)
	}

	institutionID, ok := inst.(int)
	if !ok {
		return global.JSONResponseWithErrorV1(c, "500", "Invalid institution id type", nil, 500)
	}

	projectID, err := strconv.Atoi(c.Params("project_id"))
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "400", "Invalid project_id", err, 400)
	}

	// Get project
	project, err := script.GetProjectByProjectID(projectID)
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Failed to get project", err, 500)
	}

	if project == nil {
		return global.JSONResponseWithErrorV1(c, "404", "Project not found", nil, 404)
	}

	// Get server
	server, err := script.GetServerByServerID(int(project.ServerID))
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Failed to get server", err, 500)
	}

	if server == nil {
		return global.JSONResponseWithErrorV1(c, "404", "Server not found", nil, 404)
	}

	// Ownership check
	if server.InstitutionID != uint(institutionID) {
		return global.JSONResponseWithErrorV1(c, "403", "Forbidden", nil, 403)
	}

	// Decrypt fields
	projectName, err := encrypDecryptV1.DecryptV1(project.ProjectName, config.SecretKey)
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Decrypt project name failed", err, 500)
	}

	environment, err := encrypDecryptV1.DecryptV1(project.Environment, config.SecretKey)
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Decrypt environment failed", err, 500)
	}

	description, err := encrypDecryptV1.DecryptV1(project.Description, config.SecretKey)
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Decrypt description failed", err, 500)
	}

	response := struct {
		ProjectID   uint   `json:"project_id"`
		ServerID    uint   `json:"server_id"`
		ProjectName string `json:"project_name"`
		Environment string `json:"environment"`
		Description string `json:"description"`
		Status      string `json:"status"`
	}{
		ProjectID:   project.ProjectID,
		ServerID:    project.ServerID,
		ProjectName: projectName,
		Environment: environment,
		Description: description,
		Status:      project.Status,
	}

	return global.JSONResponseWithDataV1(
		c,
		"200",
		"Project retrieved successfully",
		response,
		200,
	)
}

func GetServers(c fiber.Ctx) error {
	inst := c.Locals("institution_id")
	if inst == nil {
		return global.JSONResponseWithErrorV1(c, "401", "Unauthorized institution", nil, 401)
	}

	institutionID, ok := inst.(int)
	if !ok {
		return global.JSONResponseWithErrorV1(c, "500", "Invalid institution id type", nil, 500)
	}

	// SCRIPT CALL (DB ONLY)
	servers, err := script.GetServersByInstitutionID(institutionID)
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Failed to retrieve servers", err, 500)
	}

	if len(servers) == 0 {
		return global.JSONResponseWithErrorV1(c, "404", "No servers found", nil, 404)
	}

	type ServerResponse struct {
		ServerID      uint   `json:"server_id"`
		ServerName    string `json:"server_name"`
		ServerIP      string `json:"server_ip"`
		InstitutionID uint   `json:"institution_id"`
	}

	response := make([]ServerResponse, 0, len(servers))

	for _, server := range servers {
		serverName, err := encrypDecryptV1.DecryptV1(server.ServerName, config.SecretKey)
		if err != nil {
			return global.JSONResponseWithErrorV1(c, "500", "Failed to decrypt server name", err, 500)
		}

		serverIP, err := encrypDecryptV1.DecryptV1(server.ServerIP, config.SecretKey)
		if err != nil {
			return global.JSONResponseWithErrorV1(c, "500", "Failed to decrypt server IP", err, 500)
		}

		response = append(response, ServerResponse{
			ServerID:      server.ServerID,
			ServerName:    serverName,
			ServerIP:      serverIP,
			InstitutionID: server.InstitutionID,
		})
	}

	return global.JSONResponseWithDataV1(
		c,
		"200",
		"Servers retrieved successfully",
		response,
		200,
	)
}

func GetProjects(c fiber.Ctx) error {
	inst := c.Locals("institution_id")
	if inst == nil {
		return global.JSONResponseWithErrorV1(c, "401", "Unauthorized institution", nil, 401)
	}

	institutionID, ok := inst.(int)
	if !ok {
		return global.JSONResponseWithErrorV1(c, "500", "Invalid institution id type", nil, 500)
	}

	serverID, err := strconv.Atoi(c.Params("server_id"))
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "400", "Invalid server_id", err, 400)
	}

	// Verify that the server belongs to the institution
	servers, err := script.GetServersByInstitutionID(institutionID)
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Failed to retrieve servers", err, 500)
	}

	found := false
	for _, server := range servers {
		if server.ServerID == uint(serverID) {
			found = true
			break
		}
	}

	if !found {
		return global.JSONResponseWithErrorV1(c, "404", "Server not found", nil, 404)
	}

	// Get projects
	projects, err := script.GetAllProjectsByServerID(serverID)
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Failed to retrieve projects", err, 500)
	}

	type ProjectResponse struct {
		ProjectID   uint   `json:"project_id"`
		ServerID    uint   `json:"server_id"`
		ProjectName string `json:"project_name"`
		Environment string `json:"environment"`
		Description string `json:"description"`
		Status      string `json:"status"`
	}

	response := make([]ProjectResponse, 0, len(projects))

	for _, project := range projects {
		projectName, err := encrypDecryptV1.DecryptV1(project.ProjectName, config.SecretKey)
		if err != nil {
			return global.JSONResponseWithErrorV1(c, "500", "Failed to decrypt project name", err, 500)
		}

		environment, err := encrypDecryptV1.DecryptV1(project.Environment, config.SecretKey)
		if err != nil {
			return global.JSONResponseWithErrorV1(c, "500", "Failed to decrypt environment", err, 500)
		}

		description, err := encrypDecryptV1.DecryptV1(project.Description, config.SecretKey)
		if err != nil {
			return global.JSONResponseWithErrorV1(c, "500", "Failed to decrypt description", err, 500)
		}

		response = append(response, ProjectResponse{
			ProjectID:   project.ProjectID,
			ServerID:    project.ServerID,
			ProjectName: projectName,
			Environment: environment,
			Description: description,
			Status:      project.Status,
		})
	}

	return global.JSONResponseWithDataV1(
		c,
		"200",
		"Projects retrieved successfully",
		response,
		200,
	)
}
