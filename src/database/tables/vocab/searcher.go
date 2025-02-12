package vocab

import (
	"database/sql"
	"fmt"
)

func SearchVocab(vocab string) ([]Entry, error) {
	query := fmt.Sprintf("SELECT VocabId FROM %s WHERE %s = ?", READ_TABLE_NAME, COL_READING)
	db, err := sql.Open(SQL_PRODUCT, DATABASE_NAME)
	if err != nil {
		return nil, err
	}

	rows, err := db.Query(query, vocab)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []Entry
	for rows.Next() {
		var entry Entry
		if err := rows.Scan(&entry.EntSeq); err != nil {
			return nil, err
		}

		entries = append(entries, entry)
	}

	return entries, nil
}
