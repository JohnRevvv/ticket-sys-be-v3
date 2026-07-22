package script

import (
	"time"
	"ideyanale-be/pkg/config"
	Ticketmodel"ideyanale-be/pkg/modules/tickets/model"
)

func EndorseTicket(ticketID string) error {

	now := time.Now()

	err := config.DBConnList[0].
		Model(&Ticketmodel.Ticket{}).
		Where("ticket_id = ?", ticketID).
		Updates(map[string]interface{}{
			"status":      "for approval",
			"endorsed_at": now,
		}).Error

	return err
}

func ApproveTicket(ticketID string, approverID uint) error {

	now := time.Now()

	err := config.DBConnList[0].
		Model(&Ticketmodel.Ticket{}).
		Where("ticket_id = ?", ticketID).
		Updates(map[string]interface{}{
			"approver_id": approverID,
			"approved_at": now,
			"status":      "for resolution",
		}).Error

	return err
}

func ResolveTicket(ticketID string, resolverID uint) error {

	now := time.Now()

	err := config.DBConnList[0].
		Model(&Ticketmodel.Ticket{}).
		Where("ticket_id = ?", ticketID).
		Updates(map[string]interface{}{
			"resolver_id": resolverID,
			"resolved_at": now,
			"status":      "resolved",
		}).Error

	return err
}

