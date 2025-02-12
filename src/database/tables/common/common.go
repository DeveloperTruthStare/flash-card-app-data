package common

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

func TableHasRows(db *sql.DB, tableName string) (bool, error) {
	query := fmt.Sprintf(`
		SELECT EXISTS (SELECT 1 FROM %s LIMIT 1)
		`, tableName,
	)

	var exists bool
	err := db.QueryRow(query).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("error checking rows in table %s: %w", tableName, err)
	}
	return exists, nil
}
