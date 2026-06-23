package controller

import (
	"time"

	"github.com/gofiber/fiber/v3"

	global "ideyanale-be/pkg/global/json_response"
	script "ideyanale-be/pkg/modules/users/script"
)

func GetUsersByInstitutionID(c fiber.Ctx) error {

	ID := c.Params("institution_id")
	if ID == "" {
		return global.JSONResponseWithErrorV1(c, "400", "Institution ID is required", nil, 400,)
	}

	users, err := script.GetUsersByInstitutionID(ID)
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Failed to fetch users", err, 500,)
	}

	return global.JSONResponseWithDataV1(c, "200",	"Users fetched successfully", users, 200,)
}

func GetUsersByID(c fiber.Ctx) error {

	ID := c.Params("id")
	if ID == "" {
		return global.JSONResponseWithErrorV1(c, "400", "User ID is required", nil, 400)
	}

	user, err := script.GetUserByID(ID)
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Failed to fetch user", err, 500)
	}

	return global.JSONResponseWithDataV1(c, "200", "User fetched successfully", user, 200)
}

// helper (keeps controller clean)
func now() time.Time {
	return time.Now()
}
