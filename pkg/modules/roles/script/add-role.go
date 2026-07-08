package script

import (
	"ideyanale-be/pkg/config"

	"gorm.io/gorm"
)

type DefaultRole struct {
	RoleName         string
	CanCreateTicket  bool
	CanEndorseTicket bool
	CanApproveTicket bool
	CanResolveTicket bool
	CanAudit         bool
}

var DefaultRoles = []DefaultRole{
	{
		RoleName:         "Insti-Admin",
		CanApproveTicket: true,
		CanEndorseTicket: true,
		CanAudit:         true,
		CanResolveTicket: true,
		CanCreateTicket:  true,
	},
	{
		RoleName:         "Approver",
		CanApproveTicket: true,
	},
	{
		RoleName:         "Endorser",
		CanEndorseTicket: true,
	},
	{
		RoleName: "Auditor",
		CanAudit: true,
	},
	{
		RoleName:         "Resolver",
		CanResolveTicket: true,
		CanAudit: true,
	},
	{
		RoleName:        "User",
		CanCreateTicket: true,
	},
}

func AddDefaultRoles(institutionID uint) error {
	return config.DBConnList[0].Transaction(func(tx *gorm.DB) error {

		for _, role := range DefaultRoles {

			if err := tx.Exec(`
				INSERT INTO roles (
					institution_id,
					role_name,
					can_create,
					can_endorse,
					can_approve,
					can_resolve,
					can_audit,
					created_at,
					updated_at
				)
				VALUES (?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
			`,
				institutionID,
				role.RoleName,
				role.CanCreateTicket,
				role.CanEndorseTicket,
				role.CanApproveTicket,
				role.CanResolveTicket,
				role.CanAudit,
			).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func AddRole(roleName string, institutionID uint, canCreate, canEndorse, canApprove, canResolve, canAudit bool) error {

	return config.DBConnList[0].Exec(`
		INSERT INTO roles (
			institution_id,
			role_name,
			can_create,
			can_endorse,
			can_approve,
			can_resolve,
			can_audit,
			created_at,
			updated_at
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
	`,
		institutionID,
		roleName,
		canCreate,
		canEndorse,
		canApprove,
		canResolve,
		canAudit,
	).Error
}
