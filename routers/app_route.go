package routers

import (
	"net/http"

	global "ideyanale-be/pkg/global/json_response"
	middleware "ideyanale-be/pkg/middleware/autologout"
	jwtMiddleware "ideyanale-be/pkg/middleware/jwt"
	loggerV1 "ideyanale-be/pkg/middleware/logger/v1"
	instiadminController "ideyanale-be/pkg/modules/insti-admin/controller"
	institutionController "ideyanale-be/pkg/modules/institutions/controller"
	projectController "ideyanale-be/pkg/modules/projects/controller"
	superadminController "ideyanale-be/pkg/modules/super-admin/controller"
	ticketController "ideyanale-be/pkg/modules/tickets/controller"
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
	auth.Post("/login-otp", userController.LoginWithOTP)
	auth.Post("/verify-otp", userController.VerifyOTP)

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
	protected.Patch("/change-role-admin/:id", superadminController.ChangeRoleToAdmin)
	protected.Patch("/user/:id/status", superadminController.ChangeUserStatus)
	protected.Get("/users/:institution_id", userController.GetUsersByInstitutionID)
	protected.Post("/logout/super-admin", superadminController.LogoutSuperAdmin)

	//Insti Admin
	protected.Get("/get-user", userController.GetUsersByInstitutionID)
	protected.Post("/add/job-position", instiadminController.AddPosition)
	protected.Get("/job-positions-by-institution", instiadminController.GetPositionsByInstitutionID)
	protected.Post("/add-ticket-types", instiadminController.AddTicketType)
	protected.Post("/add-category", instiadminController.AddCategory)
	protected.Post("/add-sub-category", instiadminController.AddSubCategory)
	protected.Post("/add-new-role", instiadminController.AddRole)
	protected.Patch("/user/change-role/:id", superadminController.ChangeUserRole)

	protected.Patch("/edit-ticket-type-info/:ticket_type_id", instiadminController.EditTicketType)
	protected.Patch("/edit-category-info/:category_id", instiadminController.EditCategory)
	protected.Patch("/edit-sub-category-info/:sub_category_id", instiadminController.EditSubCategory)

	protected.Get("/get-ticket-type/:ticket_type_id", instiadminController.GetTicketTypeByID)
	protected.Get("/get-category/:category_id", instiadminController.GetCategoryByID)
	protected.Get("/get-sub-category/:sub_category_id", instiadminController.GetSubCategoryByID)

	protected.Get("/get-ticket-types", instiadminController.GetAllTicketTypes)
	protected.Get("/get-categories/:ticket_type_id", instiadminController.GetAllCategories)
	protected.Get("/get-sub-categories/:category_id", instiadminController.GetAllSubCategories)

	//User
	protected.Post("/logout", userController.Logout)
	protected.Get("/get-user/details/:id", userController.GetUserByID)

	//Institution
	institution := apiV1.Group("/institution", jwtMiddleware.JWTProtected(), middleware.AutoLogout())
	institution.Post("/create", institutionController.AddInstitution)
	institution.Get("/get", institutionController.GetInstitutions)
	institution.Post("/edit/:institution_id", institutionController.EditInstitution)

	//Ticket
	ticket := apiV1.Group("/ticket", jwtMiddleware.JWTProtected(), middleware.AutoLogout())
	ticket.Post("/create", ticketController.CreateNewTicket)

	//Project
	project := apiV1.Group("/project", jwtMiddleware.JWTProtected(), middleware.AutoLogout())
	project.Post("/server/create", projectController.AddServer)
	project.Post("/create", projectController.AddProject)
	project.Get("/get/server/:server_id", projectController.GetServerByID)
	project.Get("/get/project/:project_id", projectController.GetProjectByID)
	project.Get("/get/servers", projectController.GetServers)
	project.Get("/get/projects/:server_id", projectController.GetProjects)
}
