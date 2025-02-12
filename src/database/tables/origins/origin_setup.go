package origins

import (
	"database/sql"
	"flashcardmaker/src/database/tables/common"
	"flashcardmaker/src/logger"
	"fmt"
	"sync"

	_ "modernc.org/sqlite"
)

const (
	SQL_PRODUCT     = "sqlite"
	DATABASE_NAME   = "dictionary.db"
	TABLE_NAME      = "VocabOrigin"
	COL_ORIGIN_ID   = "OriginId"
	COL_ORIGIN_TEXT = "Origin"

	JMDICT_ORIGIN  = 1
	POKEMON_ORIGIN = 2
	OOKAMI_ORIGIN  = 3
)

var EXTERNAL_PRODUCTS = []string{
	"JMdict",
	"Pokemon",
	"Ookami",
}

var (
	mu sync.Mutex
)

func DropTable() error {
	mu.Lock()
	defer mu.Unlock()

	db, err := sql.Open(SQL_PRODUCT, DATABASE_NAME)
	if err != nil {
		return err
	}
	defer db.Close()

	dropTableSql := fmt.Sprintf(`
		DROP TABLE IF EXISTS %s;
	`, TABLE_NAME)

	_, err = db.Exec(dropTableSql)
	if err != nil {
		return err
	}

	logger.Info("Dropped Origin Table")
	return nil
}

func Initialize() error {
	mu.Lock()
	defer mu.Unlock()

	db, err := sql.Open(SQL_PRODUCT, DATABASE_NAME)
	if err != nil {
		return err
	}
	defer db.Close()

	err = createOriginTable(db)
	if err != nil {
		return err
	}

	initialized, err := common.TableHasRows(db, TABLE_NAME)
	if err != nil {
		return err
	}

	if initialized {
		logger.Info("%s has been initialized", TABLE_NAME)
		return nil
	}

	err = populateOriginTable(db)
	if err != nil {
		return err
	}

	logger.Info("Successfully initialized %s.", TABLE_NAME)
	return nil
}

func createOriginTable(db *sql.DB) error {
	createTableSQL := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
		"%s" INTEGER PRIMARY KEY AUTOINCREMENT,
		"%s" TEXT NOT NULL	
		);
	`, TABLE_NAME, COL_ORIGIN_ID, COL_ORIGIN_TEXT)

	_, err := db.Exec(createTableSQL)
	if err != nil {
		return err
	}

	logger.Debug("Created Table: %s", TABLE_NAME)
	return nil
}

func populateOriginTable(db *sql.DB) error {
	for _, product := range EXTERNAL_PRODUCTS {
		insertSQL := fmt.Sprintf(
			`INSERT INTO %s (%s) VALUES (?)`,
			TABLE_NAME, COL_ORIGIN_TEXT,
		)
		_, err := db.Exec(insertSQL, product)
		if err != nil {
			return err
		}
		logger.Debug("Inserted %s into %s", product, TABLE_NAME)
	}

	return nil
}
