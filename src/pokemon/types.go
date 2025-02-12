package pokemon

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