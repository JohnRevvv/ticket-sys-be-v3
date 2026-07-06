package script

import (
	"ideyanale-be/pkg/config"
	"ideyanale-be/pkg/modules/projects/model"
	ProjModel "ideyanale-be/pkg/modules/projects/model"
)

func GetServersByInstitutionID(institutionID int) ([]ProjModel.ServerModel, error) {
	var allServers []ProjModel.ServerModel

	err := config.DBConnList[0].Raw(`
		SELECT
			server_id,
			server_name,
			server_ip,
			institution_id
		FROM servers
		WHERE institution_id = ?
	`, institutionID).Scan(&allServers).Error

	return allServers, err
}

func GetServerByServerID(serverID int) (*model.ServerModel, error) {
	var server model.ServerModel

	err := config.DBConnList[0].Raw(`
		SELECT
			server_id,
			institution_id,
			server_name,
			server_ip,
			status
		FROM servers
		WHERE server_id = ?
	`, serverID).Scan(&server).Error

	if err != nil {
		return nil, err
	}

	if server.ServerID == 0 {
		return nil, nil
	}

	return &server, nil
}

func GetAllProjectsByServerID(serverID int) ([]model.ProjectModel, error) {
	var projects []model.ProjectModel

	err := config.DBConnList[0].Raw(`
		SELECT
			project_id,
			server_id,
			project_name,
			environment,
			description,
			status
		FROM projects
		WHERE server_id = ?
	`, serverID).Scan(&projects).Error

	return projects, err
}

func GetProjectByProjectID(projectID int) (*model.ProjectModel, error) {
	var project model.ProjectModel

	err := config.DBConnList[0].Raw(`
		SELECT
			project_id,
			server_id,
			project_name,
			environment,
			description,
			status
		FROM projects
		WHERE project_id = ?
	`, projectID).Scan(&project).Error

	if err != nil {
		return nil, err
	}

	if project.ProjectID == 0 {
		return nil, nil
	}

	return &project, nil
}

func ServerExists(institutionID uint, serverName, serverIP string) (bool, error) {
	var serverexists bool

	err := config.DBConnList[0].Raw(`
		SELECT EXISTS (
			SELECT 1
			FROM servers
			WHERE institution_id = ?
			  AND server_name = ?
			  AND server_ip = ?
		)
	`,
		institutionID,
		serverName,
		serverIP,
	).Scan(&serverexists).Error

	if err != nil {
		return false, err
	}

	return serverexists, nil
}

func ProjectExists(serverID uint, projectName, environment string) (bool, error) {
	var projectexists bool

	err := config.DBConnList[0].Raw(`
		SELECT EXISTS (
			SELECT 1
			FROM projects
			WHERE server_id = ?
			  AND project_name = ?
			  AND environment = ?
		)
	`,
		serverID,
		projectName,
		environment,
	).Scan(&projectexists).Error

	if err != nil {
		return false, err
	}

	return projectexists, nil
}
