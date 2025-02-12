package vocab

import (
	"database/sql"
	"encoding/json"
	"flashcardmaker/src/database/tables/common"
	"flashcardmaker/src/database/tables/origins"
	"flashcardmaker/src/logger"
	"fmt"
	"sync"

	_ "modernc.org/sqlite"
)

const ( // Database Info
	SQL_PRODUCT   = "sqlite"
	DATABASE_NAME = "dictionary.db"
)

const ( // Vocab Table Info
	VOCAB_TABLE_NAME = "Vocab"
	COL_VOCAB_ID     = "VocabId"
	COL_ORIGIN_ID    = origins.COL_ORIGIN_ID
)

const ( // Kanji Reading Table  Info
	VIS_TABLE_NAME = "Visualization"
	COL_VISUAL     = "Visual"
)

const ( // Kana Reading Table Info
	READ_TABLE_NAME = "Readings"
	COL_READING     = "Reading"
)

const ( // English Definitions Table Info
	SENSE_TABLE_NAME = "Senses"
	COL_JSON         = "SenseJson"
)

var (
	mu sync.Mutex
)

var tableNames = []string{
	VOCAB_TABLE_NAME,
	VIS_TABLE_NAME,
	READ_TABLE_NAME,
	SENSE_TABLE_NAME,
}

func DropTables() error {
	mu.Lock()
	defer mu.Unlock()

	db, err := sql.Open(SQL_PRODUCT, DATABASE_NAME)
	if err != nil {
		return err
	}
	defer db.Close()

	for _, tableName := range tableNames {
		dropTableQuery := fmt.Sprintf(
			`DROP TABLE IF EXISTS %s;`,
			tableName,
		)

		_, err = db.Exec(dropTableQuery)
		if err != nil {
			return err
		}
		logger.Debug("Successfully dropped %s.", tableName)
	}
	logger.Info("Dropped all Vocab Tables")
	return nil
}

func LoadDefaultDictionary() error {
	logger.Debug("Loading default Dictionary")
	return Initialize(DEFAULT_DICTIONARY)
}

func Initialize(dictionaryPath string) error {
	mu.Lock()
	defer mu.Unlock()

	db, err := sql.Open(SQL_PRODUCT, DATABASE_NAME)
	if err != nil {
		return err
	}
	defer db.Close()

	err = createTables(db)
	if err != nil {
		return err
	}

	initialized, err := common.TableHasRows(db, VOCAB_TABLE_NAME)
	if err != nil {
		return err
	}

	if initialized {
		logger.Info("Vocab Tables are already initialized, skipping.")
		return nil
	}

	data, err := ParseJMDict(dictionaryPath)
	if err != nil {
		return err
	}

	logger.Error("%d entries found", len(data))

	err = populateTables(db, data)
	if err != nil {
		return err
	}

	logger.Info("%s Tables Successfully initialized!", VOCAB_TABLE_NAME)

	return nil
}

func createTables(db *sql.DB) error {
	err := createVocabTable(db)
	if err != nil {
		return err
	}

	err = createVisualizationTable(db)
	if err != nil {
		return err
	}

	err = createReadingTable(db)
	if err != nil {
		return err
	}

	err = createSensesTable(db)
	if err != nil {
		return err
	}

	logger.Debug("Finished Creating Vocab Tables")
	return nil
}

func createVocabTable(db *sql.DB) error {
	createTableSQL := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
		"%s" INTEGER PRIMARY KEY AUTOINCREMENT,
		"%s" TEXT NOT NULL	
		);
	`, VOCAB_TABLE_NAME, COL_VOCAB_ID, COL_ORIGIN_ID)

	_, err := db.Exec(createTableSQL)
	if err != nil {
		return err
	}

	logger.Debug("Created Table: %s", VOCAB_TABLE_NAME)
	return nil
}

func createVisualizationTable(db *sql.DB) error {
	createTableSQL := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
		"%s" INTEGER NOT NULL,
		"%s" TEXT NOT NULL	
		);
	`, VIS_TABLE_NAME, COL_VOCAB_ID, COL_VISUAL)

	_, err := db.Exec(createTableSQL)
	if err != nil {
		return err
	}

	logger.Debug("Created Table: %s", VIS_TABLE_NAME)
	return nil
}

func createReadingTable(db *sql.DB) error {
	createTableSQL := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
		"%s" INTEGER NOT NULL,
		"%s" TEXT NOT NULL	
		);
	`, READ_TABLE_NAME, COL_VOCAB_ID, COL_READING)

	_, err := db.Exec(createTableSQL)
	if err != nil {
		return err
	}

	logger.Debug("Created Table: %s", READ_TABLE_NAME)
	return nil
}

func createSensesTable(db *sql.DB) error {
	createTableSQL := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
		"%s" INTEGER NOT NULL,
		"%s" TEXT NOT NULL	
		);
	`, SENSE_TABLE_NAME, COL_VOCAB_ID, COL_JSON)

	_, err := db.Exec(createTableSQL)
	if err != nil {
		return err
	}

	logger.Debug("Created Table: %s", SENSE_TABLE_NAME)
	return nil
}

func populateTables(db *sql.DB, data []Entry) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	vocabStmt, err := tx.Prepare(fmt.Sprintf(`
    	INSERT INTO %s (%s, %s)
		    SELECT ?, ?
    	WHERE NOT EXISTS (
        	SELECT 1 FROM %s WHERE %s = ?
		);
	`, VOCAB_TABLE_NAME, COL_VOCAB_ID, COL_ORIGIN_ID, VOCAB_TABLE_NAME, COL_VOCAB_ID))
	if err != nil {
		return err
	}
	defer vocabStmt.Close()

	kanjiStmt, err := tx.Prepare(fmt.Sprintf(`
		INSERT INTO %s (%s, %s) VALUES (?, ?);
	`, VIS_TABLE_NAME, COL_VOCAB_ID, COL_VISUAL))
	if err != nil {
		return err
	}
	defer kanjiStmt.Close()

	readingStmt, err := tx.Prepare(fmt.Sprintf(`
		INSERT INTO %s (%s, %s) VALUES (?, ?);
	`, READ_TABLE_NAME, COL_VOCAB_ID, COL_READING))
	if err != nil {
		return err
	}
	defer readingStmt.Close()

	senseStmt, err := tx.Prepare(fmt.Sprintf(`
		INSERT INTO %s (%s, %s) VALUES (?, ?);
	`, SENSE_TABLE_NAME, COL_VOCAB_ID, COL_JSON))
	if err != nil {
		return err
	}
	defer senseStmt.Close()

	var prevProgress uint8 = 100

	for i, entry := range data {
		res, err := vocabStmt.Exec(entry.EntSeq, origins.JMDICT_ORIGIN, entry.EntSeq)
		if err != nil {
			return nil
		}
		rows, err := res.RowsAffected()
		if err != nil {
			return err
		}
		if rows == 0 {
			continue
		}

		for _, kanji := range entry.Kanji {
			_, err = kanjiStmt.Exec(entry.EntSeq, kanji)
			if err != nil {
				return err
			}
		}

		for _, reading := range entry.Reading {
			_, err = readingStmt.Exec(entry.EntSeq, reading)
			if err != nil {
				return err
			}
		}

		for _, sense := range entry.Sense {
			jsonData, err := json.MarshalIndent(sense, "", " ")
			if err != nil {
				return err
			}

			_, err = senseStmt.Exec(entry.EntSeq, jsonData)
			if err != nil {
				return err
			}
		}

		var newProgress uint8 = (uint8)(100 * i / len(data))
		if newProgress != prevProgress {
			logger.Progress("Loading Vocab", newProgress)
		}
		prevProgress = newProgress
	}

	return tx.Commit()
}
