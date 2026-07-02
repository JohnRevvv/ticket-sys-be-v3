package controller

import (
	"ideyanale-be/pkg/config"
	global "ideyanale-be/pkg/global/json_response"
	encrypDecryptV1 "ideyanale-be/pkg/middleware/encryption/v1"
	jwt "ideyanale-be/pkg/middleware/jwt"
	"ideyanale-be/pkg/modules/insti-admin/script"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v3"
)

// ==========================
// POST CONTROLLERS
// ==========================
func AddTicketType(c fiber.Ctx) error {

	if err := jwt.RequireRoles(c, "Insti-Admin"); err != nil {
		return global.JSONResponseWithErrorV1(c, "403", "Forbidden", err, 403)
	}

	type Req struct {
		TicketTypeName string `json:"ticket_type_name"`
	}

	var req Req
	if err := c.Bind().Body(&req); err != nil {
		return global.JSONResponseWithErrorV1(c, "400", "Invalid request", err, 400)
	}

	trimmedNewTicketTypeName := strings.TrimSpace(req.TicketTypeName)

	if trimmedNewTicketTypeName == "" {
		return global.JSONResponseWithErrorV1(c, "400", "Ticket type name is required", nil, 400)
	}

	// normalize spaces
	fields := strings.Fields(trimmedNewTicketTypeName)
	normalizedName := strings.Join(fields, " ")

	inst := c.Locals("institution_id")
	if inst == nil {
		return global.JSONResponseWithErrorV1(c, "401", "Unauthorized institution", nil, 401)
	}

	institutionID, ok := inst.(int)
	if !ok {
		return global.JSONResponseWithErrorV1(c, "500", "Invalid institution id type", nil, 500)
	}

	allTicketTypes, err := script.GetTicketTypesByInstitutionID(institutionID)
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Fetch existing ticket types failed", err, 500)
	}

	for _, tt := range allTicketTypes {
		decTicketTypeName, err := encrypDecryptV1.DecryptV1(tt.TicketTypeName, config.SecretKey)
		if err != nil {
			return global.JSONResponseWithErrorV1(c, "500", "Decrypt ticket type name failed", err, 500)
		}

		if strings.EqualFold(strings.TrimSpace(decTicketTypeName), normalizedName) {
			return global.JSONResponseWithErrorV1(c, "409", "Ticket type name already exists", nil, 409)
		}
	}

	encTicketTypeName, err := encrypDecryptV1.EncryptV1(normalizedName, config.SecretKey)
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Encrypt ticket type name failed", err, 500)
	}

	if err := script.AddTicketType(encTicketTypeName, institutionID); err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Add ticket type failed", err, 500)
	}

	return global.JSONResponseV1(c, "200", "Ticket Type added successfully", 200)
}

func AddCategory(c fiber.Ctx) error {

	if err := jwt.RequireRoles(c, "Insti-Admin"); err != nil {
		return global.JSONResponseWithErrorV1(c, "403", "Forbidden", err, 403)
	}

	type Req struct {
		TicketTypeID uint   `json:"ticket_type_id"`
		CategoryName string `json:"category_name"`
	}

	var req Req

	if err := c.Bind().Body(&req); err != nil {
		return global.JSONResponseWithErrorV1(c, "400", "Invalid request", err, 400)
	}

	trimmedCategoryName := strings.TrimSpace(req.CategoryName)

	if trimmedCategoryName == "" {
		return global.JSONResponseWithErrorV1(c, "400", "Category name is required", nil, 400)
	}

	fields := strings.Fields(trimmedCategoryName)
	normalizedName := strings.Join(fields, " ")

	inst := c.Locals("institution_id")
	if inst == nil {
		return global.JSONResponseWithErrorV1(c, "401", "Unauthorized institution", nil, 401)
	}

	institutionID, ok := inst.(int)
	if !ok {
		return global.JSONResponseWithErrorV1(c, "500", "Invalid institution id type", nil, 500)
	}

	// Fetch the ticket type the category is being added under
	ticketType, err := script.GetTicketTypeByID(int(req.TicketTypeID))
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Fetch ticket type failed", err, 500)
	}
	if ticketType == nil {
		return global.JSONResponseWithErrorV1(c, "404", "Ticket type not found", nil, 404)
	}

	// Make sure the ticket type actually belongs to the caller's institution
	if uint(institutionID) != ticketType.InstitutionID {
		return global.JSONResponseWithErrorV1(c, "403", "Ticket type does not belong to your institution", nil, 403)
	}

	// Check for duplicate category names under this ticket type
	allCategories, err := script.GetCategoriesByTicketTypeID(int(req.TicketTypeID))
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Fetch existing categories failed", err, 500)
	}

	for _, cat := range allCategories {
		decCategoryName, err := encrypDecryptV1.DecryptV1(cat.CategoryName, config.SecretKey)
		if err != nil {
			return global.JSONResponseWithErrorV1(c, "500", "Decrypt category name failed", err, 500)
		}

		if strings.EqualFold(strings.TrimSpace(decCategoryName), normalizedName) {
			return global.JSONResponseWithErrorV1(c, "409", "Category name already exists", nil, 409)
		}
	}

	encCategoryName, err := encrypDecryptV1.EncryptV1(normalizedName, config.SecretKey)
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Encrypt category name failed", err, 500)
	}

	if err = script.AddCategory(encCategoryName, int(req.TicketTypeID)); err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Add category failed", err, 500)
	}

	return global.JSONResponseV1(c, "200", "Category added successfully", 200)
}

func AddSubCategory(c fiber.Ctx) error {

	if err := jwt.RequireRoles(c, "Insti-Admin"); err != nil {
		return global.JSONResponseWithErrorV1(c, "403", "Forbidden", err, 403)
	}

	type Req struct {
		CategoryID      uint   `json:"category_id"`
		SubjectName     string `json:"subject_name"`
		SubCategoryName string `json:"sub_category_name"`
		Description     string `json:"description"`
		HasDuration     bool   `json:"has_duration"`
		DurationDays    int    `json:"duration_days"`
	}

	var req Req

	if err := c.Bind().Body(&req); err != nil {
		return global.JSONResponseWithErrorV1(c, "400", "Invalid request", err, 400)
	}

	trimmedSubCategoryName := strings.TrimSpace(req.SubCategoryName)

	if req.CategoryID == 0 {
		return global.JSONResponseWithErrorV1(c, "400", "Category id is required", nil, 400)
	}
	if trimmedSubCategoryName == "" {
		return global.JSONResponseWithErrorV1(c, "400", "Sub category name is required", nil, 400)
	}

	fields := strings.Fields(trimmedSubCategoryName)
	normalizedName := strings.Join(fields, " ")

	inst := c.Locals("institution_id")
	if inst == nil {
		return global.JSONResponseWithErrorV1(c, "401", "Unauthorized institution", nil, 401)
	}

	institutionID, ok := inst.(int)
	if !ok {
		return global.JSONResponseWithErrorV1(c, "500", "Invalid institution id type", nil, 500)
	}

	// Validate the category exists
	category, err := script.GetCategoryByID(int(req.CategoryID))
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Fetch category failed", err, 500)
	}
	if category == nil {
		return global.JSONResponseWithErrorV1(c, "404", "Category not found", nil, 404)
	}

	// Walk up to ticket type to confirm institution ownership
	ticketType, err := script.GetTicketTypeByID(int(category.TicketTypeID))
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Fetch ticket type failed", err, 500)
	}
	if ticketType == nil {
		return global.JSONResponseWithErrorV1(c, "404", "Ticket type not found", nil, 404)
	}
	if uint(institutionID) != ticketType.InstitutionID {
		return global.JSONResponseWithErrorV1(c, "403", "Category does not belong to your institution", nil, 403)
	}

	// Check for duplicate sub category names under this category
	allSubCategories, err := script.GetSubCategoriesByCategoryID(int(req.CategoryID))
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Fetch existing sub categories failed", err, 500)
	}

	for _, sc := range allSubCategories {
		decSubCategoryName, err := encrypDecryptV1.DecryptV1(sc.SubCategoryName, config.SecretKey)
		if err != nil {
			return global.JSONResponseWithErrorV1(c, "500", "Decrypt sub category name failed", err, 500)
		}

		if strings.EqualFold(strings.TrimSpace(decSubCategoryName), normalizedName) {
			return global.JSONResponseWithErrorV1(c, "409", "Sub category name already exists", nil, 409)
		}
	}

	if req.HasDuration {
		switch req.DurationDays {
		case 30, 60, 90:
			// valid
		default:
			return global.JSONResponseWithErrorV1(c, "400", "Duration must be 30, 60, or 90 days", nil, 400)
		}
	} else {
		req.DurationDays = 0
	}

	encSubCategoryName, err := encrypDecryptV1.EncryptV1(normalizedName, config.SecretKey)
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Encrypt sub category name failed", err, 500)
	}

	encSubjectName, err := encrypDecryptV1.EncryptV1(req.SubjectName, config.SecretKey)
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Encrypt subject name failed", err, 500)
	}

	encDescription, err := encrypDecryptV1.EncryptV1(req.Description, config.SecretKey)
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Encrypt subject name failed", err, 500)
	}

	if err := script.AddSubCategory(
		encSubCategoryName,
		encSubjectName,
		encDescription,
		req.HasDuration,
		req.DurationDays,
		int(req.CategoryID),
	); err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Failed to add sub category", err, 500)
	}

	return global.JSONResponseV1(c, "200", "Sub Category added successfully", 200)
}

//==========================
// GET CONTROLLERS
//==========================

func GetTicketTypeByID(c fiber.Ctx) error {
	inst := c.Locals("institution_id")
	if inst == nil {
		return global.JSONResponseWithErrorV1(c, "401", "Unauthorized institution", nil, 401)
	}

	institutionID, ok := inst.(int)
	if !ok {
		return global.JSONResponseWithErrorV1(c, "500", "Invalid institution id type", nil, 500)
	}

	ticketTypeID, err := strconv.Atoi(c.Params("ticket_type_id"))
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "400", "Invalid ticket_type_id", err, 400)
	}

	// SCRIPT CALL (DB ONLY)
	ticketType, err := script.GetTicketTypeByID(ticketTypeID)
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Failed to get ticket type", err, 500)
	}

	if ticketType == nil {
		return global.JSONResponseWithErrorV1(c, "404", "Ticket type not found", nil, 404)
	}

	// OWNERSHIP CHECK
	if int(ticketType.InstitutionID) != institutionID {
		return global.JSONResponseWithErrorV1(c, "403", "Forbidden", nil, 403)
	}

	// DECRYPT (CONTROLLER RESPONSIBILITY)
	decryptedName, err := encrypDecryptV1.DecryptV1(ticketType.TicketTypeName, config.SecretKey)
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Decrypt ticket type name failed", err, 500)
	}

	// map decrypted value (DON'T mutate DB struct directly if reused)
	response := struct {
		TicketTypeID   uint   `json:"ticket_type_id"`
		InstitutionID  uint   `json:"institution_id"`
		TicketTypeName string `json:"ticket_type_name"`
		Status         string `json:"status"`
	}{
		TicketTypeID:   ticketType.TicketTypeID,
		InstitutionID:  ticketType.InstitutionID,
		TicketTypeName: decryptedName,
		Status:         ticketType.Status,
	}

	return global.JSONResponseWithDataV1(
		c,
		"200",
		"Ticket type retrieved successfully",
		response,
		200,
	)
}

func GetCategoryByID(c fiber.Ctx) error {
	inst := c.Locals("institution_id")
	if inst == nil {
		return global.JSONResponseWithErrorV1(c, "401", "Unauthorized institution", nil, 401)
	}

	categoryID, err := strconv.Atoi(c.Params("category_id"))
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "400", "Invalid ticket_type_id", err, 400)
	}

	ticketType, err := script.GetCategoryByID(categoryID)
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Failed to get ticket type", err, 500)
	}

	if ticketType == nil {
		return global.JSONResponseWithErrorV1(c, "404", "Ticket type not found", nil, 404)
	}

	// DECRYPT (CONTROLLER RESPONSIBILITY)
	decName, err := encrypDecryptV1.DecryptV1(ticketType.CategoryName, config.SecretKey)
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Decrypt ticket type name failed", err, 500)
	}

	response := struct {
		CategoryID   uint   `json:"category_id"`
		TicketTypeID uint   `json:"ticket_type_id"`
		CategoryName string `json:"category_name"`
		Status       string `json:"status"`
	}{
		CategoryID:   ticketType.CategoryID,
		TicketTypeID: ticketType.TicketTypeID,
		CategoryName: decName,
		Status:       ticketType.Status,
	}

	// SUCCESS RESPONSE (using your V1 data response)
	return global.JSONResponseWithDataV1(c, "200", "Ticket type retrieved successfully", response, 200)
}

func GetSubCategoryByID(c fiber.Ctx) error {
	inst := c.Locals("institution_id")
	if inst == nil {
		return global.JSONResponseWithErrorV1(c, "401", "Unauthorized institution", nil, 401)
	}

	subcategoryID, err := strconv.Atoi(c.Params("sub_category_id"))
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "400", "Invalid ticket_type_id", err, 400)
	}

	subcategory, err := script.GetSubCategoryByID(subcategoryID)
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Failed to get ticket type", err, 500)
	}

	if subcategory == nil {
		return global.JSONResponseWithErrorV1(c, "404", "Sub-Category not found", nil, 404)
	}

	// DECRYPT
	decSubjectName, err := encrypDecryptV1.DecryptV1(subcategory.SubCategoryName, config.SecretKey)
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Decrypt ticket type name failed", err, 500)
	}

	decSubCategoryName, err := encrypDecryptV1.DecryptV1(subcategory.SubjectName, config.SecretKey)
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Decrypt ticket type name failed", err, 500)
	}

	decDescription, err := encrypDecryptV1.DecryptV1(subcategory.Description, config.SecretKey)
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Decrypt ticket type name failed", err, 500)
	}

	response := struct {
		SubCategoryID   uint   `json:"sub_category_id"`
		CategoryID      uint   `json:"category_id"`
		SubjectName     string `json:"subject_name"`
		SubCategoryName string `json:"sub_category_name"`
		Description     string `json:"description"`
		HasDuration     bool   `json:"has_duration"`
		DurationDays    int    `json:"duration_days"`
		Status          string `json:"status"`
	}{
		SubCategoryID:   subcategory.SubCategoryID,
		CategoryID:      subcategory.CategoryID,
		SubjectName:     decSubjectName,
		SubCategoryName: decSubCategoryName,
		Description:     decDescription,
		HasDuration:     subcategory.HasDuration,
		DurationDays:    subcategory.DurationDays,
		Status:          subcategory.Status,
	}

	// SUCCESS RESPONSE
	return global.JSONResponseWithDataV1(c, "200", "Sub-Category Details retrieved successfully", response, 200)
}

func GetAllTicketTypes(c fiber.Ctx) error {

	inst := c.Locals("institution_id")
	if inst == nil {
		return global.JSONResponseWithErrorV1(c, "401", "Unauthorized institution", nil, 401)
	}

	institutionID, ok := inst.(int)
	if !ok {
		return global.JSONResponseWithErrorV1(c, "500", "Invalid institution id type", nil, 500)
	}

	rows, err := script.GetTicketTypesByInstitutionID(institutionID)
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Failed to fetch ticket type", err, 500)
	}

	type TicketTypeResp struct {
		TicketTypeID   uint   `json:"ticket_type_id"`
		TicketTypeName string `json:"ticket_type_name"`
	}

	var resp []TicketTypeResp
	for _, row := range rows {
		decName, err := encrypDecryptV1.DecryptV1(row.TicketTypeName, config.SecretKey)
		if err != nil {
			return global.JSONResponseWithErrorV1(c, "500", "Decrypt ticket type name failed", err, 500)
		}

		resp = append(resp, TicketTypeResp{
			TicketTypeID:   row.TicketTypeID,
			TicketTypeName: decName,
		})
	}

	return global.JSONResponseWithDataV1(c, "200", "Ticket type fetched successfully", resp, 200)
}

func GetAllCategories(c fiber.Ctx) error {

	ticketypeIDParam := c.Params("ticket_type_id")
	if ticketypeIDParam == "" {
		return global.JSONResponseWithErrorV1(c, "400", "ticket_type_id path param is required", nil, 400)
	}

	tickettypeID, err := strconv.Atoi(ticketypeIDParam)
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "400", "Invalid ticket_type_id", err, 400)
	}

	inst := c.Locals("institution_id")
	if inst == nil {
		return global.JSONResponseWithErrorV1(c, "401", "Unauthorized institution", nil, 401)
	}

	institutionID, ok := inst.(int)
	if !ok {
		return global.JSONResponseWithErrorV1(c, "500", "Invalid institution id type", nil, 500)
	}

	ticketType, err := script.GetTicketTypeByID(tickettypeID)
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Fetch ticket type failed", err, 500)
	}
	if ticketType == nil {
		return global.JSONResponseWithErrorV1(c, "404", "Ticket type not found", nil, 404)
	}
	if uint(institutionID) != ticketType.InstitutionID {
		return global.JSONResponseWithErrorV1(c, "403", "Category does not belong to your institution", nil, 403)
	}

	allCategories, err := script.GetCategoriesByTicketTypeID(tickettypeID)
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Fetch sub categories failed", err, 500)
	}

	for i := range allCategories {
		decName, err := encrypDecryptV1.DecryptV1(allCategories[i].CategoryName, config.SecretKey)
		if err != nil {
			return global.JSONResponseWithErrorV1(c, "500", "Decrypt sub category name failed", err, 500)
		}
		allCategories[i].CategoryName = decName
	}

	return global.JSONResponseWithDataV1(c, "200", "Successfull", allCategories, 200)
}

func GetAllSubCategories(c fiber.Ctx) error {

	categoryIDParam := c.Params("category_id")
	if categoryIDParam == "" {
		return global.JSONResponseWithErrorV1(c, "400", "category_id is required", nil, 400)
	}

	categoryID, err := strconv.Atoi(categoryIDParam)
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "400", "Invalid category_id", err, 400)
	}

	category, err := script.GetCategoryByID(categoryID)
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Fetch category failed", err, 500)
	}
	if category == nil {
		return global.JSONResponseWithErrorV1(c, "404", "Category not found", nil, 404)
	}

	ticketType, err := script.GetTicketTypeByID(int(category.TicketTypeID))
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Fetch ticket type failed", err, 500)
	}
	if ticketType == nil {
		return global.JSONResponseWithErrorV1(c, "404", "Ticket type not found", nil, 404)
	}

	rows, err := script.GetSubCategoriesByCategoryID(categoryID)
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Fetch sub categories failed", err, 500)
	}

	type SubCategoryResponse struct {
		SubCategoryID   uint   `json:"sub_category_id"`
		CategoryID      uint   `json:"category_id"`
		SubjectName     string `json:"subject_name"`
		SubCategoryName string `json:"sub_category_name"`
		Description     string `json:"description"`
		HasDuration     bool   `json:"has_duration"`
		DurationDays    int    `json:"duration_days"`
		Status          string `json:"status"`
	}

	var response []SubCategoryResponse

	for _, row := range rows {

		subName, err := encrypDecryptV1.DecryptV1(row.SubCategoryName, config.SecretKey)
		if err != nil {
			return global.JSONResponseWithErrorV1(c, "500", "Decrypt sub_category_name failed", err, 500)
		}

		subjectName, err := encrypDecryptV1.DecryptV1(row.SubjectName, config.SecretKey)
		if err != nil {
			return global.JSONResponseWithErrorV1(c, "500", "Decrypt subject_name failed", err, 500)
		}

		description, err := encrypDecryptV1.DecryptV1(row.Description, config.SecretKey)
		if err != nil {
			return global.JSONResponseWithErrorV1(c, "500", "Decrypt description failed", err, 500)
		}

		response = append(response, SubCategoryResponse{
			SubCategoryID:   row.SubCategoryID,
			CategoryID:      row.CategoryID,
			SubjectName:     subjectName,
			SubCategoryName: subName,
			Description:     description,
			HasDuration:     row.HasDuration,
			DurationDays:    row.DurationDays,
			Status:          row.Status,
		})
	}

	return global.JSONResponseWithDataV1(c, "200", "Successful", response, 200)
}

//==========================
// EDIT CONTROLLERS
//==========================

func EditTicketType(c fiber.Ctx) error {
	if err := jwt.RequireRoles(c, "Insti-Admin"); err != nil {
		return global.JSONResponseWithErrorV1(c, "403", "Forbidden", err, 403)
	}

	ticketTypeID, err := strconv.Atoi(c.Params("ticket_type_id"))
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "400", "Invalid ticket_type_id", err, 400)
	}

	// SCRIPT CALL (DB ONLY)
	ticketType, err := script.GetTicketTypeByID(ticketTypeID)
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Failed to get ticket type", err, 500)
	}

	if ticketType == nil {
		return global.JSONResponseWithErrorV1(c, "404", "Ticket type not found", nil, 404)
	}

	institutionID, ok := c.Locals("institution_id").(int)
	if !ok {
		return global.JSONResponseWithErrorV1(c, "401", "Unauthorized institution", nil, 401)
	}

	if ticketType.InstitutionID != uint(institutionID) {
		return global.JSONResponseWithErrorV1(c, "403", "Ticket type does not belong to your institution", nil, 403)
	}

	type Req struct {
		TicketTypeName string `json:"ticket_type_name"`
		Status         string `json:"status"`
	}

	var req Req
	if err := c.Bind().Body(&req); err != nil {
		return global.JSONResponseWithErrorV1(c, "400", "Invalid request", err, 400)
	}

	req.TicketTypeName = strings.TrimSpace(req.TicketTypeName)
	req.Status = strings.TrimSpace(req.Status)

	if req.TicketTypeName == "" {
		return global.JSONResponseWithErrorV1(c, "400", "Ticket type name is required", nil, 400)
	}

	// normalize spaces
	fields := strings.Fields(req.TicketTypeName)
	normalizedName := strings.Join(fields, " ")

	// Check duplicate ticket type names within the same institution
	allTicketTypes, err := script.GetTicketTypesByInstitutionID(int(ticketType.InstitutionID))
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Fetch existing ticket types failed", err, 500)
	}

	for _, tt := range allTicketTypes {
		if tt.TicketTypeID == ticketType.TicketTypeID {
			continue
		}

		name, err := encrypDecryptV1.DecryptV1(tt.TicketTypeName, config.SecretKey)
		if err != nil {
			return global.JSONResponseWithErrorV1(c, "500", "Decrypt ticket type name failed", err, 500)
		}

		if strings.EqualFold(strings.TrimSpace(name), normalizedName) {
			return global.JSONResponseWithErrorV1(c, "409", "Ticket type name already exists", nil, 409)
		}
	}

	if req.Status == "" {
		req.Status = ticketType.Status
	}

	encTicketTypeName, err := encrypDecryptV1.EncryptV1(normalizedName, config.SecretKey)
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Encrypt ticket type name failed", err, 500)
	}

	err = script.EditTicketType(
		int(ticketTypeID),
		encTicketTypeName,
		req.Status,
	)
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Edit ticket type failed", err, 500)
	}

	return global.JSONResponseV1(c, "200", "Ticket Type updated successfully", 200)
}

func EditCategory(c fiber.Ctx) error {
	if err := jwt.RequireRoles(c, "Insti-Admin"); err != nil {
		return global.JSONResponseWithErrorV1(c, "403", "Forbidden", err, 403)
	}

	categoryID, err := strconv.Atoi(c.Params("category_id"))
	if err != nil || categoryID <= 0 {
		return global.JSONResponseWithErrorV1(c, "400", "Invalid category_id", nil, 400)
	}

	// Get category
	category, err := script.GetCategoryByID(categoryID)
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Failed to get category", err, 500)
	}

	if category == nil {
		return global.JSONResponseWithErrorV1(c, "404", "Category not found", nil, 404)
	}

	// Get institution id from JWT
	institutionID, ok := c.Locals("institution_id").(int)
	if !ok {
		return global.JSONResponseWithErrorV1(c, "401", "Unauthorized institution", nil, 401)
	}

	// Verify ownership through ticket type
	ticketType, err := script.GetTicketTypeByID(int(category.TicketTypeID))
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Failed to get ticket type", err, 500)
	}

	if ticketType == nil {
		return global.JSONResponseWithErrorV1(c, "404", "Ticket type not found", nil, 404)
	}

	if ticketType.InstitutionID != uint(institutionID) {
		return global.JSONResponseWithErrorV1(c, "403", "Category does not belong to your institution", nil, 403)
	}

	type Req struct {
		CategoryName string `json:"category_name"`
		Status       string `json:"status"`
	}

	var req Req
	if err := c.Bind().Body(&req); err != nil {
		return global.JSONResponseWithErrorV1(c, "400", "Invalid request", err, 400)
	}

	req.CategoryName = strings.TrimSpace(req.CategoryName)
	req.Status = strings.TrimSpace(req.Status)

	if req.CategoryName == "" {
		return global.JSONResponseWithErrorV1(c, "400", "Category name is required", nil, 400)
	}

	// normalize spaces
	fields := strings.Fields(req.CategoryName)
	normalizedName := strings.Join(fields, " ")

	// Check duplicates within the same ticket type
	allCategories, err := script.GetCategoriesByTicketTypeID(int(category.TicketTypeID))
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Failed to fetch categories", err, 500)
	}

	for _, cat := range allCategories {
		// Skip current category
		if cat.CategoryID == category.CategoryID {
			continue
		}

		name, err := encrypDecryptV1.DecryptV1(
			cat.CategoryName,
			config.SecretKey,
		)
		if err != nil {
			return global.JSONResponseWithErrorV1(c, "500", "Decrypt category name failed", err, 500)
		}

		if strings.EqualFold(strings.TrimSpace(name), normalizedName) {
			return global.JSONResponseWithErrorV1(c, "409", "Category name already exists", nil, 409)
		}
	}

	if req.Status == "" {
		req.Status = category.Status
	}

	encCategoryName, err := encrypDecryptV1.EncryptV1(normalizedName, config.SecretKey)
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Encrypt category name failed", err, 500)
	}

	err = script.EditCategory(
		categoryID,
		encCategoryName,
		req.Status,
	)
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Edit category failed", err, 500)
	}

	return global.JSONResponseV1(c, "200", "Category updated successfully", 200)
}

func EditSubCategory(c fiber.Ctx) error {
	if err := jwt.RequireRoles(c, "Insti-Admin"); err != nil {
		return global.JSONResponseWithErrorV1(c, "403", "Forbidden", err, 403)
	}

	subCategoryID, err := strconv.Atoi(c.Params("sub_category_id"))
	if err != nil || subCategoryID <= 0 {
		return global.JSONResponseWithErrorV1(c, "400", "Invalid sub_category_id", nil, 400)
	}

	subCategory, err := script.GetSubCategoryByID(subCategoryID)
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Failed to get sub category", err, 500)
	}
	if subCategory == nil {
		return global.JSONResponseWithErrorV1(c, "404", "Sub category not found", nil, 404)
	}

	institutionID, ok := c.Locals("institution_id").(int)
	if !ok {
		return global.JSONResponseWithErrorV1(c, "401", "Unauthorized institution", nil, 401)
	}

	category, err := script.GetCategoryByID(int(subCategory.CategoryID))
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Failed to get category", err, 500)
	}
	if category == nil {
		return global.JSONResponseWithErrorV1(c, "404", "Category not found", nil, 404)
	}

	ticketType, err := script.GetTicketTypeByID(int(category.TicketTypeID))
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Failed to get ticket type", err, 500)
	}
	if ticketType == nil {
		return global.JSONResponseWithErrorV1(c, "404", "Ticket type not found", nil, 404)
	}

	if ticketType.InstitutionID != uint(institutionID) {
		return global.JSONResponseWithErrorV1(c, "403", "Sub category does not belong to your institution", nil, 403)
	}

	type Req struct {
		SubCategoryName string `json:"sub_category_name"`
		SubjectName     string `json:"subject_name"`
		Description     string `json:"description"`
		HasDuration     bool   `json:"has_duration"`
		DurationDays    int    `json:"duration_days"`
		Status          string `json:"status"`
	}

	var req Req
	if err := c.Bind().Body(&req); err != nil {
		return global.JSONResponseWithErrorV1(c, "400", "Invalid request", err, 400)
	}

	req.SubCategoryName = strings.TrimSpace(req.SubCategoryName)
	req.SubjectName = strings.TrimSpace(req.SubjectName)
	req.Description = strings.TrimSpace(req.Description)
	req.Status = strings.TrimSpace(req.Status)

	if req.SubCategoryName == "" {
		return global.JSONResponseWithErrorV1(c, "400", "Sub category name is required", nil, 400)
	}

	if req.SubjectName == "" {
		return global.JSONResponseWithErrorV1(c, "400", "Subject name is required", nil, 400)
	}

	if req.HasDuration {
		switch req.DurationDays {
		case 30, 60, 90:
			// valid
		default:
			return global.JSONResponseWithErrorV1(
				c,
				"400",
				"Duration must be 30, 60, or 90 days",
				nil,
				400,
			)
		}
	} else {
		req.DurationDays = 0
	}

	// normalize spaces
	fields := strings.Fields(req.SubCategoryName)
	normalizedName := strings.Join(fields, " ")

	subCategories, err := script.GetSubCategoriesByCategoryID(int(subCategory.CategoryID))
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Failed to fetch sub categories", err, 500)
	}

	for _, sc := range subCategories {
		if sc.SubCategoryID == subCategory.SubCategoryID {
			continue
		}

		name, err := encrypDecryptV1.DecryptV1(sc.SubCategoryName, config.SecretKey)
		if err != nil {
			return global.JSONResponseWithErrorV1(c, "500", "Decrypt sub category name failed", err, 500)
		}

		if strings.EqualFold(strings.TrimSpace(name), normalizedName) {
			return global.JSONResponseWithErrorV1(c, "409", "Sub category name already exists", nil, 409)
		}
	}

	if req.Status == "" {
		req.Status = subCategory.Status
	}

	encSubCategoryName, err := encrypDecryptV1.EncryptV1(normalizedName, config.SecretKey)
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Encrypt sub category name failed", err, 500)
	}
	encSubjectName, err := encrypDecryptV1.EncryptV1(req.SubjectName, config.SecretKey)
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Encrypt subject name failed", err, 500)
	}
	encDescription, err := encrypDecryptV1.EncryptV1(req.Description, config.SecretKey)
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Encrypt description failed", err, 500)
	}

	err = script.EditSubCategory(
		subCategoryID,
		encSubCategoryName,
		encSubjectName,
		encDescription,
		req.HasDuration,
		req.DurationDays,
		req.Status,
	)

	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Edit sub category failed", err, 500)
	}

	return global.JSONResponseV1(c, "200", "Sub category updated successfully", 200)
}
