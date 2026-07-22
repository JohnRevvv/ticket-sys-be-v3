package controller

import (
	"strings"

	global "ideyanale-be/pkg/global/json_response"
	"ideyanale-be/pkg/middleware/jwt"
	ticketScript "ideyanale-be/pkg/modules/tickets/script"

	"github.com/gofiber/fiber/v3"
)

type TicketActionRequest struct {
	Action string `json:"action"`
}

func ProcessTicket(c fiber.Ctx) error {

	type Req struct {
		Action string `json:"action"`
	}

	var req Req

	if err := c.Bind().Body(&req); err != nil {
		return global.JSONResponseWithErrorV1(c, "400", "Invalid request body", err, 400)
	}

	ticketID := c.Params("ticket_id")
	if ticketID == "" {
		return global.JSONResponseWithErrorV1(c, "400", "Ticket ID is required", nil, 400)
	}

	userID, ok := c.Locals("id").(int)
	if !ok {
		return global.JSONResponseWithErrorV1(c, "401", "Unauthorized user", nil, 401)
	}

	institutionID, ok := c.Locals("institution_id").(uint)
	if !ok {
		return global.JSONResponseWithErrorV1(c, "401", "Unauthorized institution", nil, 401)
	}

	ticket, err := ticketScript.GetTicketByTicketID(ticketID)
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "404", "Ticket not found", err, 404)
	}

	switch strings.ToLower(req.Action) {

	case "endorse":

		if err := jwt.RequirePermission(c, func(p jwt.Permissions) bool {
			return p.CanEndorseTicket
		}); err != nil {
			return global.JSONResponseWithErrorV1(
				c,
				"403",
				"You cannot endorse tickets",
				nil,
				403,
			)
		}

		// Cannot endorse own ticket
		if ticket.SubmitterID == uint(userID) {
			return global.JSONResponseWithErrorV1(
				c,
				"403",
				"Submitter cannot endorse their own ticket",
				nil,
				403,
			)
		}

		// Must be assigned endorser
		if ticket.EndorserID != uint(userID) {
			return global.JSONResponseWithErrorV1(
				c,
				"403",
				"Ticket is not assigned to you",
				nil,
				403,
			)
		}

		if strings.ToLower(ticket.Status) != "for endorsement" {
			return global.JSONResponseWithErrorV1(
				c,
				"400",
				"Ticket is not for endorsement",
				nil,
				400,
			)
		}

		if err := ticketScript.EndorseTicket(ticket.TicketID); err != nil {
			return global.JSONResponseWithErrorV1(
				c,
				"500",
				"Failed to endorse ticket",
				err,
				500,
			)
		}

	case "approve":

		if err := jwt.RequirePermission(c, func(p jwt.Permissions) bool {
			return p.CanApproveTicket
		}); err != nil {
			return global.JSONResponseWithErrorV1(
				c,
				"403",
				"You cannot approve tickets",
				nil,
				403,
			)
		}

		// Cannot approve own ticket
		if ticket.SubmitterID == uint(userID) {
			return global.JSONResponseWithErrorV1(
				c,
				"403",
				"Submitter cannot approve their own ticket",
				nil,
				403,
			)
		}

		// Cannot approve ticket you endorsed
		if ticket.EndorserID == uint(userID) {
			return global.JSONResponseWithErrorV1(
				c,
				"403",
				"Endorser cannot approve the same ticket",
				nil,
				403,
			)
		}

		if ticket.InstitutionPool != institutionID {
			return global.JSONResponseWithErrorV1(
				c,
				"403",
				"Ticket does not belong to your institution pool",
				nil,
				403,
			)
		}

		if strings.ToLower(ticket.Status) != "for approval" {
			return global.JSONResponseWithErrorV1(
				c,
				"400",
				"Ticket is not for approval",
				nil,
				400,
			)
		}

		if err := ticketScript.ApproveTicket(
			ticket.TicketID,
			uint(userID),
		); err != nil {
			return global.JSONResponseWithErrorV1(
				c,
				"500",
				"Failed to approve ticket",
				err,
				500,
			)
		}

	case "resolve":

		if err := jwt.RequirePermission(c, func(p jwt.Permissions) bool {
			return p.CanResolveTicket
		}); err != nil {
			return global.JSONResponseWithErrorV1(
				c,
				"403",
				"You cannot resolve tickets",
				nil,
				403,
			)
		}

		// Cannot resolve own ticket
		if ticket.SubmitterID == uint(userID) {
			return global.JSONResponseWithErrorV1(
				c,
				"403",
				"Submitter cannot resolve their own ticket",
				nil,
				403,
			)
		}

		if ticket.InstitutionPool != institutionID {
			return global.JSONResponseWithErrorV1(
				c,
				"403",
				"Ticket does not belong to your institution pool",
				nil,
				403,
			)
		}

		if strings.ToLower(ticket.Status) != "for resolution" {
			return global.JSONResponseWithErrorV1(
				c,
				"400",
				"Ticket is not for resolution",
				nil,
				400,
			)
		}

		if err := ticketScript.ResolveTicket(
			ticket.TicketID,
			uint(userID),
		); err != nil {
			return global.JSONResponseWithErrorV1(
				c,
				"500",
				"Failed to resolve ticket",
				err,
				500,
			)
		}

	default:
		return global.JSONResponseWithErrorV1(
			c,
			"400",
			"Invalid action. Allowed values: endorse, approve, resolve",
			nil,
			400,
		)
	}

	updatedTicket, err := ticketScript.GetTicketByTicketID(ticketID)
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Failed to fetch updated ticket", err, 500)
	}

	return global.JSONResponseWithDataV1(
		c,
		"200",
		"Ticket processed successfully",
		updatedTicket,
		200,
	)
}
