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

//==========================
// POST CONTROLLERS
//==========================
func AddTicketType(c fiber.Ctx) error {

	if err := jwt.RequireRoles(c, "insti-admin"); err != nil {
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

		if strings.EqualFold(strings.TrimSpace(decTicketTypeName), trimmedNewTicketTypeName) {
			return global.JSONResponseWithErrorV1(c, "409", "Ticket type name already exists", nil, 409)
		}
	}

	encTicketTypeName, err := encrypDecryptV1.EncryptV1(trimmedNewTicketTypeName, config.SecretKey)
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Encrypt ticket type name failed", err, 500)
	}

	if err := script.AddPosition(encTicketTypeName, institutionID); err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Add ticket type failed", err, 500)
	}

	return global.JSONResponseV1(c, "200", "Ticket Type added successfully", 200)
}

func AddCategory(c fiber.Ctx) error {

	if err := jwt.RequireRoles(c, "insti-admin"); err != nil {
		return global.JSONResponseWithErrorV1(
			c,
			"403",
			"Forbidden",
			err,
			403,
		)
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

		if strings.EqualFold(strings.TrimSpace(decCategoryName), trimmedCategoryName) {
			return global.JSONResponseWithErrorV1(c, "409", "Category name already exists", nil, 409)
		}
	}

	encCategoryName, err := encrypDecryptV1.EncryptV1(req.CategoryName, config.SecretKey)
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Encrypt category name failed", err, 500)
	}

	err = script.AddCategory(encCategoryName, int(req.TicketTypeID))
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Add category failed", err, 500)
	}

	return global.JSONResponseV1(c, "200", "Category added successfully", 200)
}

func AddSubCategory(c fiber.Ctx) error {

	if err := jwt.RequireRoles(c, "insti-admin"); err != nil {
		return global.JSONResponseWithErrorV1(c, "403", "Forbidden", err, 403)
	}

	type Req struct {
		CategoryID      uint   `json:"category_id"`
		SubjectName     string `json:"subject_name"`
		SubCategoryName string `json:"sub_category_name"`
		Description     string `json:"description"`
	}

	var req Req

	if err := c.Bind().Body(&req); err != nil {
		return global.JSONResponseWithErrorV1(c, "400", "Invalid request", err, 400)
	}

	trimmedSubCategoryName := strings.TrimSpace(req.SubCategoryName)
	trimmedSubjectName := strings.TrimSpace(req.SubjectName)

	if req.CategoryID == 0 {
		return global.JSONResponseWithErrorV1(c, "400", "Category id is required", nil, 400)
	}
	if trimmedSubCategoryName == "" {
		return global.JSONResponseWithErrorV1(c, "400", "Sub category name is required", nil, 400)
	}
	if trimmedSubjectName == "" {
		return global.JSONResponseWithErrorV1(c, "400", "Subject name is required", nil, 400)
	}

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

		if strings.EqualFold(strings.TrimSpace(decSubCategoryName), trimmedSubCategoryName) {
			return global.JSONResponseWithErrorV1(c, "409", "Sub category name already exists", nil, 409)
		}
	}

	encSubCategoryName, err := encrypDecryptV1.EncryptV1(req.SubCategoryName, config.SecretKey)
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Encrypt sub category name failed", err, 500)
	}

	err = script.AddSubCategory(encSubCategoryName, trimmedSubjectName, req.Description, int(req.CategoryID))
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Add sub category failed", err, 500)
	}

	return global.JSONResponseV1(c, "200", "Sub Category added successfully", 200)
}

//==========================
// GET CONTROLLERS
//==========================

func GetTicketTypeByInstitutionID(c fiber.Ctx) error {

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
		PositionID   uint   `json:"position_id"`
		PositionName string `json:"position_name"`
	}

	var resp []TicketTypeResp
	for _, row := range rows {
		decName, err := encrypDecryptV1.DecryptV1(row.TicketTypeName, config.SecretKey)
		if err != nil {
			return global.JSONResponseWithErrorV1(c, "500", "Decrypt ticket type name failed", err, 500)
		}

		resp = append(resp, TicketTypeResp{
			PositionID:   row.TicketTypeID,
			PositionName: decName,
		})
	}

	return global.JSONResponseWithDataV1(c, "200", "Ticket type fetched successfully", resp, 200)
}

func GetCategory(c fiber.Ctx) error {

	if err := jwt.RequireRoles(c, "insti-admin"); err != nil {
		return global.JSONResponseWithErrorV1(c, "403", "Forbidden", err, 403)
	}

	ticketypeIDParam := c.Query("ticket_type_id")
	if ticketypeIDParam  == "" {
		return global.JSONResponseWithErrorV1(c, "400", "category_id query param is required", nil, 400)
	}

	tickettypeID, err := strconv.Atoi(ticketypeIDParam )
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

	tickettype, err := script.GetTicketTypeByID(tickettypeID)
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Fetch category failed", err, 500)
	}
	if tickettype == nil {
		return global.JSONResponseWithErrorV1(c, "404", "Category not found", nil, 404)
	}

	ticketType, err := script.GetTicketTypeByID(int(tickettype.TicketTypeID))
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

	return global.JSONResponseWithDataV1(c, "200","Successfull", allCategories, 200)
}

func GetSubCategory(c fiber.Ctx) error {

	if err := jwt.RequireRoles(c, "insti-admin"); err != nil {
		return global.JSONResponseWithErrorV1(c, "403", "Forbidden", err, 403)
	}

	categoryIDParam := c.Query("category_id")
	if categoryIDParam == "" {
		return global.JSONResponseWithErrorV1(c, "400", "category_id query param is required", nil, 400)
	}

	categoryID, err := strconv.Atoi(categoryIDParam)
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "400", "Invalid category_id", err, 400)
	}

	inst := c.Locals("institution_id")
	if inst == nil {
		return global.JSONResponseWithErrorV1(c, "401", "Unauthorized institution", nil, 401)
	}

	institutionID, ok := inst.(int)
	if !ok {
		return global.JSONResponseWithErrorV1(c, "500", "Invalid institution id type", nil, 500)
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
	if uint(institutionID) != ticketType.InstitutionID {
		return global.JSONResponseWithErrorV1(c, "403", "Category does not belong to your institution", nil, 403)
	}

	allSubCategories, err := script.GetSubCategoriesByCategoryID(categoryID)
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Fetch sub categories failed", err, 500)
	}

	for i := range allSubCategories {
		decName, err := encrypDecryptV1.DecryptV1(allSubCategories[i].SubCategoryName, config.SecretKey)
		if err != nil {
			return global.JSONResponseWithErrorV1(c, "500", "Decrypt sub category name failed", err, 500)
		}
		allSubCategories[i].SubCategoryName = decName
	}

	return global.JSONResponseWithDataV1(c, "200","Successfull", allSubCategories, 200)
}

//==========================
// EDIT CONTROLLERS
//==========================

func EditSubCategory(c fiber.Ctx) error {

	if err := jwt.RequireRoles(c, "insti-admin"); err != nil {
		return global.JSONResponseWithErrorV1(c, "403", "Forbidden", err, 403)
	}

	type Req struct {
		SubCategoryID   uint   `json:"sub_category_id"`
		SubjectName     string `json:"subject_name"`
		SubCategoryName string `json:"sub_category_name"`
		Description     string `json:"description"`
		Status          string `json:"status"`
	}

	var req Req

	if err := c.Bind().Body(&req); err != nil {
		return global.JSONResponseWithErrorV1(c, "400", "Invalid request", err, 400)
	}

	if req.SubCategoryID == 0 {
		return global.JSONResponseWithErrorV1(c, "400", "Sub category id is required", nil, 400)
	}

	trimmedSubCategoryName := strings.TrimSpace(req.SubCategoryName)
	trimmedSubjectName := strings.TrimSpace(req.SubjectName)

	if trimmedSubCategoryName == "" {
		return global.JSONResponseWithErrorV1(c, "400", "Sub category name is required", nil, 400)
	}
	if trimmedSubjectName == "" {
		return global.JSONResponseWithErrorV1(c, "400", "Subject name is required", nil, 400)
	}

	inst := c.Locals("institution_id")
	if inst == nil {
		return global.JSONResponseWithErrorV1(c, "401", "Unauthorized institution", nil, 401)
	}

	institutionID, ok := inst.(int)
	if !ok {
		return global.JSONResponseWithErrorV1(c, "500", "Invalid institution id type", nil, 500)
	}

	existing, err := script.GetSubCategoryByID(int(req.SubCategoryID))
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Fetch sub category failed", err, 500)
	}
	if existing == nil {
		return global.JSONResponseWithErrorV1(c, "404", "Sub category not found", nil, 404)
	}

	category, err := script.GetCategoryByID(int(existing.CategoryID))
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
	if uint(institutionID) != ticketType.InstitutionID {
		return global.JSONResponseWithErrorV1(c, "403", "Sub category does not belong to your institution", nil, 403)
	}

	// Check duplicates against siblings, excluding itself
	allSubCategories, err := script.GetSubCategoriesByCategoryID(int(existing.CategoryID))
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Fetch existing sub categories failed", err, 500)
	}

	for _, sc := range allSubCategories {
		if sc.SubCategoryID == req.SubCategoryID {
			continue
		}

		decName, err := encrypDecryptV1.DecryptV1(sc.SubCategoryName, config.SecretKey)
		if err != nil {
			return global.JSONResponseWithErrorV1(c, "500", "Decrypt sub category name failed", err, 500)
		}

		if strings.EqualFold(strings.TrimSpace(decName), trimmedSubCategoryName) {
			return global.JSONResponseWithErrorV1(c, "409", "Sub category name already exists", nil, 409)
		}
	}

	status := strings.TrimSpace(req.Status)
	if status == "" {
		status = existing.Status
	}

	encSubCategoryName, err := encrypDecryptV1.EncryptV1(req.SubCategoryName, config.SecretKey)
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Encrypt sub category name failed", err, 500)
	}

	err = script.EditSubCategory(
		int(req.SubCategoryID),
		encSubCategoryName,
		trimmedSubjectName,
		req.Description,
		status,
	)
	if err != nil {
		return global.JSONResponseWithErrorV1(c, "500", "Edit sub category failed", err, 500)
	}

	return global.JSONResponseV1(c, "200", "Sub Category updated successfully", 200)
}
