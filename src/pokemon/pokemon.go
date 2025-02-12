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

type TypeName struct {
	TypeId     int    `csv:"type_id"`
	LanguageId int    `csv:"local_language_id"`
	Name       string `csv:"name"`
}

type Type struct {
	TypeId   int
	English  string
	Japanese string
}

type Move struct {
	Id       int `csv:"id"`
	TypeId   int `csv:"type_id"`
	Power    int `csv:"power"`
	PP       int `csv:"pp"`
	Priority int `csv:"priority"`
	English  string
	Japanese string
}

type MoveNames struct {
	MoveId     int    `csv:"move_id"`
	LanguageId int    `csv:"local_language_id"`
	Name       string `csv:"name"`
}

type Pokemon struct {
	Id          int `csv:"id"`
	Height      int `csv:"height"`
	Weight      int `csv:"weight"`
	BaseExp     int `csv:"base_experience"`
	English     string
	Japanese    string
	MovesLearnt []PokemonMoves
}

type PokemonMoves struct {
	PokemonId   int `csv:"pokemon_id"`
	MoveId      int `csv:"move_id"`
	LevelLearnt int `csv:"level"`
}

type PokemonAbilities struct {
	PokemonId int  `csv:"pokemon_id"`
	AbilityId int  `csv:"ability_id"`
	IsHidden  bool `csv:"is_hidden"`
}

type PokemonSpeciesName struct {
	PokemonSpeciesId int    `csv:"pokemon_species_id"`
	LanguageId       int    `csv:"local_language_id"`
	Name             string `csv:"name"`
}

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

type TypeDict struct {
	Types []Type
}

func (td *TypeDict) SetEnglish(typeId int, name string) {
	searchedFor := -1
	for index, value := range td.Types {
		if value.TypeId == typeId {
			searchedFor = index
			break
		}
	}
	if searchedFor == -1 {
		td.Types = append(td.Types, Type{TypeId: typeId, English: name})
	} else {
		td.Types[searchedFor].English = name
	}
}

func (td *TypeDict) SetJapanese(typeId int, name string) {
	searchedFor := -1
	for index, value := range td.Types {
		if value.TypeId == typeId {
			searchedFor = index
			break
		}
	}
	if searchedFor == -1 {
		td.Types = append(td.Types, Type{TypeId: typeId, Japanese: name})
	} else {
		td.Types[searchedFor].Japanese = name
	}
}

func LoadData() {
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

	for _, move := range moves {
		fmt.Printf("Move Id: %d\t English: %s\t Japanese: %s\n", move.Id, move.English, move.Japanese)
	}
}
