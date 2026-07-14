package script

import (
	"ideyanale-be/pkg/config"
	ticketModel "ideyanale-be/pkg/modules/tickets/model"
)

func GetTicketByTicketID(ticketID string) (*ticketModel.Ticket, error) {
	var ticket ticketModel.Ticket

	err := config.DBConnList[0].
		Preload("TicketAttachment").
		Where("ticket_id = ?", ticketID).
		First(&ticket).Error

	if err != nil {
		return nil, err
	}

	return &ticket, nil
}

func GetTicketsByUserID(userID uint) ([]ticketModel.Ticket, error) {
	var tickets []ticketModel.Ticket

	// Get all tickets submitted by the user
	err := config.DBConnList[0].Raw(`
		SELECT
			id,
			ticket_id,
			institution_id,
			ticket_type_id,
			category_id,
			subcategory_id,
			subject,
			description,
			due_date,
			institution_pool,
			submitter_id,
			resolver_id,
			endorser_id,
			approver_id,
			status,
			created_at,
			updated_at,
			cancelled_by,
			cancelled_at,
			cancellation_reason,
			started_at,
			resolved_at,
			resolution_minutes,
			resolution_time,
			on_hold,
			hold_at,
			closed_by,
			closed_at,
			endorsed_at,
			approved_at,
			close_token,
			close_token_used
		FROM tickets
		WHERE submitter_id = ?
		ORDER BY created_at DESC
	`, userID).Scan(&tickets).Error
	if err != nil {
		return nil, err
	}

	// Load attachments for each ticket
	for i := range tickets {
		var attachments []ticketModel.TicketAttachment

		err := config.DBConnList[0].Raw(`
			SELECT
				id,
				ticket_id,
				file_name,
				file_key,
				uploaded_by
			FROM ticket_attachments
			WHERE ticket_id = ?
			ORDER BY id ASC
		`, tickets[i].TicketID).Scan(&attachments).Error
		if err != nil {
			return nil, err
		}

		tickets[i].TicketAttachment = attachments
	}

	return tickets, nil
}

func GetAllTickets() ([]ticketModel.Ticket, error) {
	var tickets []ticketModel.Ticket

	err := config.DBConnList[0].
		Preload("TicketAttachment").
		Find(&tickets).Error

	if err != nil {
		return nil, err
	}

	return tickets, nil
}