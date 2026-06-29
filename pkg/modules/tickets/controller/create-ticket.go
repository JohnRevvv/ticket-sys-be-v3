package controller

import (
	"ideyanale-be/pkg/config"
	global "ideyanale-be/pkg/global/json_response"
	encrypDecryptV1 "ideyanale-be/pkg/middleware/encryption/v1"
	SAdScript "ideyanale-be/pkg/modules/super-admin/script"
	ticketModel "ideyanale-be/pkg/modules/tickets/model"
	ticketScript "ideyanale-be/pkg/modules/tickets/script"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
)

type CreateTicketRequest struct {
	TicketTypeID    uint       `json:"ticket_type_id"`
	CategoryID      uint       `json:"category_id"`
	SubCategoryID   uint       `json:"subcategory_id"`
	Subject         string     `json:"subject"`
	Description     string     `json:"description"`
	DueDate         *time.Time `json:"due_date"`
	InstitutionPool uint       `json:"institution_pool"`
	EndorserID      uint       `json:"endorser_id"`
}

func CreateNewTicket(c fiber.Ctx) error {
	inst := c.Locals("institution_id")
	if inst == nil {
		return global.JSONResponseWithErrorV1(c, "401", "Unauthorized institution", nil, 401)
	}

	institutionID, ok := inst.(int)
	if !ok {
		return global.JSONResponseWithErrorV1(c, "500", "Invalid institution id type", nil, 500)
	}

	submitter := c.Locals("id")
	if submitter == nil {
		return global.JSONResponseWithErrorV1(c, "401", "Unauthorized user", nil, 401)
	}

	submitterID, ok := submitter.(int)
	if !ok {
		return global.JSONResponseWithErrorV1(c, "500", "Invalid user id", nil, 500)
	}

	var req CreateTicketRequest

	if err := c.Bind().Body(&req); err != nil {
		return global.JSONResponseWithErrorV1(c, "400", "Invalid request body", err, 400)
	}

	// =====================
	// Required Fields
	// =====================

	if req.InstitutionPool == 0 {
		return global.JSONResponseWithErrorV1(c, "400", "Institution pool is required", nil, 400)
	}

	if req.TicketTypeID == 0 {
		return global.JSONResponseWithErrorV1(c, "400", "Ticket type is required", nil, 400)
	}

	if req.CategoryID == 0 {
		return global.JSONResponseWithErrorV1(c, "400", "Category is required", nil, 400)
	}

	if req.SubCategoryID == 0 {
		return global.JSONResponseWithErrorV1(c, "400", "Subcategory is required", nil, 400)
	}

	if strings.TrimSpace(req.Subject) == "" {
		return global.JSONResponseWithErrorV1(c, "400", "Subject is required", nil, 400)
	}

	if strings.TrimSpace(req.Description) == "" {
		return global.JSONResponseWithErrorV1(c, "400", "Description is required", nil, 400)
	}

	if req.DueDate == nil {
		return global.JSONResponseWithErrorV1(c, "400", "Due date is required", nil, 400)
	}

	if req.EndorserID == 0 {
		return global.JSONResponseWithErrorV1(c, "400", "Endorser is required", nil, 400)
	}

	// =====================
	// Check Institution Pool
	// =====================

	pool, err := SAdScript.GetInstitutionByID(req.InstitutionPool)
	if err != nil || pool == nil {
		return global.JSONResponseWithErrorV1(c, "404", "Institution pool not found", nil, 404)
	}

	// =====================
	// Check Ticket Type
	// =====================

	exist, err := ticketScript.IsTicketTypeBelongsToInstitution(
		req.InstitutionPool,
		req.TicketTypeID,
	)
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Failed to validate ticket type", err, 500)
	}

	if !exist {
		return global.JSONResponseWithErrorV1(c, "400", "Invalid ticket type for selected institution pool", nil, 400)
	}

	// =====================
	// Check Category
	// =====================

	exist, err = ticketScript.IsCategoryBelongsToTicketType(
		req.TicketTypeID,
		req.CategoryID,
	)
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Failed to validate category", err, 500)
	}

	if !exist {
		return global.JSONResponseWithErrorV1(c, "400", "Invalid category", nil, 400)
	}

	// =====================
	// Check SubCategory
	// =====================

	exist, err = ticketScript.IsSubCategoryBelongsToCategory(
		req.CategoryID,
		req.SubCategoryID,
	)
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Failed to validate subcategory", err, 500)
	}

	if !exist {
		return global.JSONResponseWithErrorV1(c, "400", "Invalid subcategory", nil, 400)
	}

	// =====================
	// Check Endorser
	// =====================

	isValid, err := ticketScript.IsValidEndorser(
	uint(institutionID),
	req.EndorserID,
)
if err != nil {
	return global.JSONResponseWithErrorV1(
		c,
		"500",
		"Failed to validate endorser",
		err,
		500,
	)
}

if !isValid {
	return global.JSONResponseWithErrorV1(
		c,
		"404",
		"Selected user cannot endorse tickets",
		nil,
		404,
	)
}

	// =====================
	// Encrypt
	// =====================

	encSubject, err := encrypDecryptV1.EncryptV1(req.Subject, config.SecretKey)
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Failed to encrypt subject", err, 500)
	}

	encDescription, err := encrypDecryptV1.EncryptV1(req.Description, config.SecretKey)
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Failed to encrypt description", err, 500)
	}

	ticket := ticketModel.Ticket{
		TicketID:        ticketScript.GenerateTicketID(),
		InstitutionID:   uint(institutionID),
		TicketTypeID:    req.TicketTypeID,
		CategoryID:      req.CategoryID,
		SubCategoryID:   req.SubCategoryID,
		Subject:         encSubject,
		Description:     encDescription,
		DueDate:         req.DueDate,
		InstitutionPool: req.InstitutionPool,
		SubmitterID:     uint(submitterID),
		EndorserID:      req.EndorserID,
	}

	if err := ticketScript.CreateTicket(&ticket); err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Failed to create ticket", err, 500)
	}

	return global.JSONResponseWithDataV1(c, "200", "Ticket created successfully", ticket, 200)
}
