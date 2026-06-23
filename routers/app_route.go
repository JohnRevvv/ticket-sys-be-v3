package routers

import (
	"net/http"

	global "ideyanale-be/pkg/global/json_response"
	middleware "ideyanale-be/pkg/middleware/autologout"
	jwtMiddleware "ideyanale-be/pkg/middleware/jwt"
	loggerV1 "ideyanale-be/pkg/middleware/logger/v1"
	instiadminController "ideyanale-be/pkg/modules/insti-admin/controller"
	superadminController "ideyanale-be/pkg/modules/super-admin/controller"
	"ideyanale-be/pkg/modules/users/controller"
	userController "ideyanale-be/pkg/modules/users/controller"
	crlDataEncryptionV1 "ideyanale-be/pkg/services/data_encryption/controller/v1"

	"github.com/gofiber/fiber/v3"
)

func AppRoutes(app *fiber.App) {

	app.Get("/", func(c fiber.Ctx) error {
		loggerV1.SystemLogger("API Health Check", "HealthCheck", "api_health", "HealthCheck", "Success", nil, "API is running...")
		return global.JSONResponseV1(c, "200", "API is running...", http.StatusOK)
	})

	apiV1 := app.Group("/api/v1")

	apiV1.Get("/", func(c fiber.Ctx) error {
		loggerV1.SystemLogger("API V1 Health Check", "system", "api_v1_health", "HealthCheck", "Success", nil, "API version 1 is running...")
		return global.JSONResponseV1(c, "200", "API version 1 is running...", http.StatusOK)
	})

	// =========================
	// UTILITY (PUBLIC)
	// =========================
	utility := apiV1.Group("/utility")
	utility.Post("/encrypt-data", crlDataEncryptionV1.EncrypDecryptV1)
	utility.Post("/decrypt-data", crlDataEncryptionV1.DecryptDataV1)

	// =========================
	// AUTH (PUBLIC)
	// =========================
	auth := apiV1.Group("/auth")
	auth.Post("/login-otp", controller.LoginWithOTP)
	auth.Post("/verify-otp", controller.VerifyOTP)

	// =========================
	// USER (PUBLIC OR PARTIAL)
	// =========================
	public := apiV1.Group("/public")

	public.Post("/register/super-admin", superadminController.CreateSuperAdmin)
	public.Post("/register-user", userController.RegisterUser)
	public.Post("/login/super-admin", superadminController.LoginSuperAdmin)

	// protected := apiV1.Group("/protected", jwtMiddleware.JWTProtected())
	protected := apiV1.Group("/protected", jwtMiddleware.JWTProtected(), middleware.AutoLogout())
	

	//Super Admin
	protected.Post("/add-institution", superadminController.AddInstitution)
	protected.Get("/institutions", superadminController.GetInstitutions)
	protected.Patch("/change-role-admin/:id", superadminController.ChangeRoleToAdmin)
	protected.Get("/users/:institution_id", userController.GetUsersByInstitutionID)


	//Insti Admin
	protected.Get("/get-user", userController.GetUsersByInstitutionID)
	protected.Post("/add/job-position", instiadminController.AddPosition)
	protected.Get("/job-positions-by-institution", instiadminController.GetPositionsByInstitutionID)
	protected.Post("/add-ticket-types", instiadminController.AddTicketType)
	protected.Post("/add-category", instiadminController.AddCategory)
	protected.Post("/add-sub-category", instiadminController.AddSubCategory)

	protected.Get("/get-ticket-types", instiadminController.GetTicketTypeByInstitutionID)
	protected.Get("/get-categories", instiadminController.GetCategory)
	protected.Get("/get-sub-categories", instiadminController.GetSubCategory)
	

	//User
	protected.Post("/logout", userController.Logout)
	protected.Get("/get-user/details/:id", userController.GetUsersByID)

}
