package controller

import (
	"ideyanale-be/pkg/config"
	global "ideyanale-be/pkg/global/json_response"
	encrypDecryptV1 "ideyanale-be/pkg/middleware/encryption/v1"
	InstiScript "ideyanale-be/pkg/modules/institution/script"
	ticketModel "ideyanale-be/pkg/modules/tickets/model"
	ticketScript "ideyanale-be/pkg/modules/tickets/script"
	"ideyanale-be/pkg/services/s3_service"
	"strconv"
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

	ticketTypeID, _ := strconv.ParseUint(c.FormValue("ticket_type_id"), 10, 64)
	categoryID, _ := strconv.ParseUint(c.FormValue("category_id"), 10, 64)
	subCategoryID, _ := strconv.ParseUint(c.FormValue("subcategory_id"), 10, 64)
	institutionPool, _ := strconv.ParseUint(c.FormValue("institution_pool"), 10, 64)
	endorserID, _ := strconv.ParseUint(c.FormValue("endorser_id"), 10, 64)

	req.TicketTypeID = uint(ticketTypeID)
	req.CategoryID = uint(categoryID)
	req.SubCategoryID = uint(subCategoryID)
	req.InstitutionPool = uint(institutionPool)
	req.EndorserID = uint(endorserID)
	req.Subject = c.FormValue("subject")
	req.Description = c.FormValue("description")

	if dueDate := c.FormValue("due_date"); dueDate != "" {
		parsedTime, err := time.Parse(time.RFC3339, dueDate)
		if err != nil {
			return global.JSONResponseWithErrorV1(c, "400", "Invalid due date format", err, 400)
		}
		req.DueDate = &parsedTime
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

	pool, err := InstiScript.GetInstitutionByID(req.InstitutionPool)
	if err != nil || pool == nil {
		return global.JSONResponseWithErrorV1(c, "404", "Institution pool not found", nil, 404)
	}

	// =====================
	// Check Ticket Type
	// =====================

	exist, err := ticketScript.IsTicketTypeBelongsToInstitution(req.InstitutionPool, req.TicketTypeID)
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Failed to validate ticket type", err, 500)
	}

	if !exist {
		return global.JSONResponseWithErrorV1(c, "400", "Invalid ticket type for selected institution pool", nil, 400)
	}

	// =====================
	// Check Category
	// =====================

	exist, err = ticketScript.IsCategoryBelongsToTicketType(req.TicketTypeID, req.CategoryID)
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Failed to validate category", err, 500)
	}

	if !exist {
		return global.JSONResponseWithErrorV1(c, "400", "Invalid category", nil, 400)
	}

	// =====================
	// Check SubCategory
	// =====================

	exist, err = ticketScript.IsSubCategoryBelongsToCategory(req.CategoryID, req.SubCategoryID)
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Failed to validate subcategory", err, 500)
	}

	if !exist {
		return global.JSONResponseWithErrorV1(c, "400", "Invalid subcategory", nil, 400)
	}

	// =====================
	// Check Endorser
	// =====================

	isValid, err := ticketScript.IsValidEndorser(uint(institutionID), req.EndorserID)
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Failed to validate endorser", err, 500)
	}

	if !isValid {
		return global.JSONResponseWithErrorV1(c, "404", "Selected user cannot endorse tickets", nil, 404)
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

	// ====================================
	// Upload Attachments (max 5)
	// ====================================

	multipartForm, err := c.MultipartForm()
	if err == nil {

		files := multipartForm.File["file"] // field name = file

		if len(files) > 5 {
			return global.JSONResponseWithErrorV1(
				c,
				"400",
				"Maximum of 5 attachments allowed",
				nil,
				400,
			)
		}

		s3Service, err := services.NewS3Service()
		if err != nil {
			return global.JSONResponseWithErrorV1(c, "500", "Failed to initialize S3", err, 500)
		}

		for _, fileHeader := range files {

			fileName, fileKey, err := s3Service.Upload(fileHeader, ticket.TicketID)
			if err != nil {
				return global.JSONResponseWithErrorV1(c, "500", "Failed to upload file", err, 500,)
			}

			attachment := ticketModel.TicketAttachment{
				TicketID:   ticket.TicketID,
				FileName:   fileName,
				FileKey:    fileKey,
				UploadedBy: uint(submitterID),
			}

			if err := ticketScript.CreateAttachment(&attachment); err != nil {
				return global.JSONResponseWithErrorV1(
					c,
					"500",
					"Failed to save attachment",
					err,
					500,
				)
			}
		}
	}

	ticketWithAttachments, err := ticketScript.GetTicketByTicketID(ticket.TicketID)
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Failed to fetch ticket", err, 500)
	}

	return global.JSONResponseWithDataV1(c, "200", "Ticket created successfully", ticketWithAttachments, 200)
}
