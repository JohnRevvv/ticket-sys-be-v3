package script

import (
	"fmt"
	"ideyanale-be/pkg/config"
	ticketModel "ideyanale-be/pkg/modules/tickets/model"
	// UserModel "ideyanale-be/pkg/modules/users/model"
	"sync"
)

var ticketIDMutex = &sync.Mutex{}

func GenerateTicketID() string {
	ticketIDMutex.Lock()
	defer ticketIDMutex.Unlock()

	var lastTicket ticketModel.Ticket

	err := config.DBConnList[0].Raw(`
		SELECT ticket_id
		FROM tickets
		ORDER BY ticket_id DESC
		LIMIT 1
	`).Scan(&lastTicket).Error

	if err != nil || lastTicket.TicketID == "" {
		return "SR000001"
	}

	var num int
	fmt.Sscanf(lastTicket.TicketID, "SR%06d", &num)

	num++

	return fmt.Sprintf("SR%06d", num)
}

func CreateTicket(ticket *ticketModel.Ticket) error {
	return config.DBConnList[0].Create(ticket).Error
}

func IsTicketTypeBelongsToInstitution(institutionID uint, ticketTypeID uint,) (bool, error) {
	var count int64

	err := config.DBConnList[0].Raw(`
		SELECT COUNT(*)
		FROM ticket_types
		WHERE institution_id = ?
		AND ticket_type_id = ?
	`,
		institutionID,
		ticketTypeID,
	).Scan(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func IsCategoryBelongsToTicketType(ticketTypeID uint,categoryID uint,) (bool, error) {
	var count int64

	err := config.DBConnList[0].Raw(`
		SELECT COUNT(*)
		FROM categories
		WHERE ticket_type_id = ?
		AND category_id = ?
	`,
		ticketTypeID,
		categoryID,
	).Scan(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func IsSubCategoryBelongsToCategory(categoryID uint,subCategoryID uint,) (bool, error) {
	var count int64

	err := config.DBConnList[0].Raw(`
		SELECT COUNT(*)
		FROM sub_categories
		WHERE category_id = ?
		AND sub_category_id = ?
	`,
		categoryID,
		subCategoryID,
	).Scan(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func IsValidEndorser(institutionID uint, userID uint) (bool, error) {
	var count int64

	err := config.DBConnList[0].Raw(`
		SELECT COUNT(*)
		FROM users u
		INNER JOIN roles r
			ON u.role = r.role_name
			AND u.institution_id = r.institution_id
		WHERE
			u.id = ?
			AND u.institution_id = ?
			AND r.can_endorse = TRUE
			AND r.deleted_at IS NULL
	`, userID, institutionID).Scan(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func CanUserCreateTicket(institutionID uint, userID uint,) (bool, error) {

	var count int64

	err := config.DBConnList[0].Raw(`
		SELECT COUNT(*)
		FROM users u
		INNER JOIN roles r
			ON u.role_id = r.role_id
		WHERE u.id = ?
		AND r.institution_id = ?
		AND r.can_create = true
	`,
		userID,
		institutionID,
	).Scan(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}