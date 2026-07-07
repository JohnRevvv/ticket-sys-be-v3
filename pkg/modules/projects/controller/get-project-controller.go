package controller

import (
	"ideyanale-be/pkg/config"
	global "ideyanale-be/pkg/global/json_response"
	encrypDecryptV1 "ideyanale-be/pkg/middleware/encryption/v1"
	"ideyanale-be/pkg/modules/projects/script"
	"strconv"

	"github.com/gofiber/fiber/v3"
)

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