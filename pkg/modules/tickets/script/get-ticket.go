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