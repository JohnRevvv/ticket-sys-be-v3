package script

import (
	"ideyanale-be/pkg/config"
	ticketModel "ideyanale-be/pkg/modules/tickets/model"
)

func CreateAttachment(attachment *ticketModel.TicketAttachment) error {
	return config.DBConnList[0].Exec(`
		INSERT INTO ticket_attachments (
			ticket_id,
			file_name,
			file_key,
			uploaded_by
		)
		VALUES (?, ?, ?, ?)
	`,
		attachment.TicketID,
		attachment.FileName,
		attachment.FileKey,
		attachment.UploadedBy,
	).Error
}
