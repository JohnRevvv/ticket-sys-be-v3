package script

import (
	"fmt"
	"ideyanale-be/pkg/config"
	"ideyanale-be/pkg/modules/roles/model"
)


func EditRole(roleID uint, institutionID uint, req model.Roles) error {
	result := config.DBConnList[0].Exec(`
		UPDATE roles
		SET
			role_name    = ?,
			can_create   = ?,
			can_endorse  = ?,
			can_approve  = ?,
			can_resolve  = ?,
			can_audit    = ?,
			updated_at   = NOW()
		WHERE role_id = ?
		AND institution_id = ?
		AND deleted_at IS NULL
	`,
		req.RoleName,
		req.CanCreateTicket,
		req.CanEndorseTicket,
		req.CanApproveTicket,
		req.CanResolveTicket,
		req.CanAudit,
		roleID,
		institutionID,
	)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("role not found or does not belong to this institution")
	}

	return nil
}