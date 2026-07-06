package script

import (
	"ideyanale-be/pkg/config"
)

func AddServerScript(servername, serverip string, institutionID int) error {
	return config.DBConnList[0].Exec(`
		INSERT INTO servers (
			server_name,
			server_ip,
			institution_id
		)
		VALUES (?, ?, ?)
	`,
		servername,
		serverip,
		institutionID,
	).Error
}

func AddProjectScript(projectname, environment, description string, serverID int) error {
	return config.DBConnList[0].Exec(`
		INSERT INTO projects (
			project_name,
			environment,
			description,
			server_id
		)
		VALUES (?, ?, ?, ?)
	`,
		projectname,
		environment,
		description,
		serverID,
	).Error
}