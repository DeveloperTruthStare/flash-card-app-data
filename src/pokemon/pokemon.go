package pokemon

import (
	"fmt"
	"log"
	"os"

	"github.com/gocarina/gocsv"
)

const (
	ENGLISH  = 9
	JAPANESE = 11
)

func parseCsv[T any](filepath string) ([]T, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var results []T
	if err := gocsv.UnmarshalFile(file, &results); err != nil {
		return nil, err
	}

	return results, err
}

const CSV_FILE_DIR = "external/pokeapi/data/v2/csv"

func LoadData() ([]Pokemon, error) {
	// Load Types
	var typeDict = TypeDict{
		Types: []Type{},
	}

	typenames, err := parseCsv[TypeName](CSV_FILE_DIR + "/type_names.csv")
	if err != nil {
		log.Fatal(err)
	}

	for _, row := range typenames {
		switch row.LanguageId {
		case ENGLISH:
			typeDict.SetEnglish(row.TypeId, row.Name)
		case JAPANESE:
			typeDict.SetJapanese(row.TypeId, row.Name)
		}
	}

	for _, t := range typeDict.Types {
		fmt.Printf("Type Id: %d\t English: %s\t Japanese %s\n", t.TypeId, t.English, t.Japanese)
	}

	moves, err := parseCsv[Move](CSV_FILE_DIR + "/moves.csv")
	if err != nil {
		log.Fatal(err)
	}

	moveNames, err := parseCsv[MoveNames](CSV_FILE_DIR + "/move_names.csv")
	if err != nil {
		log.Fatal(err)
	}
	for _, moveName := range moveNames {
		for index, move := range moves {
			if moveName.MoveId == move.Id {
				switch moveName.LanguageId {
				case ENGLISH:
					moves[index].English = moveName.Name
				case JAPANESE:
					moves[index].Japanese = moveName.Name
				}
			}
		}
	}

	pokemon, err := parseCsv[Pokemon](CSV_FILE_DIR + "/pokemon.csv")
	if err != nil {
		return nil, err
	}

	// get pokemon names
	pokemonNames, err := parseCsv[PokemonSpeciesName](CSV_FILE_DIR + "/pokemon_species_names.csv")
	if err != nil {
		return nil, err
	}

	for _, name := range pokemonNames {
		for index, p := range pokemon {
			if p.Id == name.PokemonSpeciesId {
				switch name.LanguageId {
				case ENGLISH:
					pokemon[index].English = name.Name
				case JAPANESE:
					pokemon[index].Japanese = name.Name
				}
			}
		}
	}

	return pokemon, nil
}
