package authenticate

import (
	"database/sql"
	"fmt"

	"github.com/BillyPurvis/boommessaging-go/database"
)

// TokenCheck Authenticate x-api-token for protected routes
func TokenCheck(token string) error {
	// Test Test
	db := database.DBCon

	// We're being lazy as QueryRow requires 3 round trips
	// to do prepared statements
	var customerID string
	stmt := fmt.Sprintf("select customer_id from api_keys where `key` = '%v'", token)
	err := db.QueryRow(stmt).Scan(&customerID)

	// Determine Error
	switch {
	case err == sql.ErrNoRows:
		return err
	case err != nil:
		return err
	default:
		return nil
	}
}
