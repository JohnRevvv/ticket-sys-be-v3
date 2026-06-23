package controller

import (
	global "ideyanale-be/pkg/global/json_response"
	"ideyanale-be/pkg/modules/super-admin/script"
	"ideyanale-be/pkg/modules/users/model"

	"github.com/gofiber/fiber/v3"
)

func ChangeRoleToAdmin(c fiber.Ctx) error {
	ID, ok := c.Locals("id").(int)
	if !ok {
		return global.JSONResponseWithErrorV1(
			c,
			"401",
			"Unauthorized",
			nil,
			401,
		)
	}

	if err := script.ChangeRoleToAdmin(&model.UserDetails{ID: ID}); err != nil {
		return global.JSONResponseWithErrorV1(
			c,
			"500",
			"Failed to change user role",
			err,
			500,
		)
	}

	return global.JSONResponseWithDataV1(
		c,
		"200",
		"User role changed to admin successfully",
		nil,
		200,
	)
}