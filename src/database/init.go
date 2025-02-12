package database

import (
	"flashcardmaker/src/database/tables/origins"
	"flashcardmaker/src/database/tables/vocab"
)

func Initialize(reinitialize bool) error {
	if reinitialize {
		err := origins.DropTable()
		if err != nil {
			return err
		}

		err = vocab.DropTables()
		if err != nil {
			return err
		}
	}
	err := origins.Initialize()
	if err != nil {
		return err
	}

	err = vocab.LoadDefaultDictionary()
	if err != nil {
		return err
	}

	return nil
}
