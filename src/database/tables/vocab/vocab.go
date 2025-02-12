package vocab

type Entry struct {
	EntSeq  int64    `json:"entSeq"`
	Kanji   []string `json:"kanji"`
	Reading []string `json:"kana"`
	Sense   []Sense  `json:"senses"`
}

type Sense struct {
	PartOfSpeech []string `json:"pos"`
	XRef         []string `json:"xref"`
	Glossary     []string `json:"gloss"`
	Misc         []string `json:"misc"`
	Field        []string `json:"field"`
}
