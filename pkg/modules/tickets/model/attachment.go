package models

// ── TicketAttachment ──────────────────────────────────────────────────────────
type TicketAttachment struct {
	ID         uint   `gorm:"primaryKey" json:"id"`
	TicketID   string `json:"ticket_id"`
	FileName   string `json:"file_name"`
	FileKey    string `json:"file_key"`
	UploadedBy uint `json:"uploaded_by"`
}

func (TicketAttachment) TableName() string {
	return "ticket_attachments"
}