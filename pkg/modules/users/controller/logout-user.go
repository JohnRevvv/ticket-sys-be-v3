package controller

import (
	global "ideyanale-be/pkg/global/json_response"
	"ideyanale-be/pkg/modules/users/script"

	"github.com/gofiber/fiber/v3"
)


func Logout(c fiber.Ctx) error {

	// get userID from context (set by JWT middleware)
	ID, ok := c.Locals("id").(int)
	if !ok {
		return global.JSONResponseWithErrorV1(c, "401", "Unauthorized", nil, 401,)
	}

	// update DB logout status
	if err := script.LogoutUser(ID); err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Failed to logout user", err, 500,)
	}

	return global.JSONResponseWithDataV1(c, "200", "Logout successful", nil, 200,)
}