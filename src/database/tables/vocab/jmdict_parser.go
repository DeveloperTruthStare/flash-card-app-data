package vocab

import (
	"encoding/json"
	"flashcardmaker/src/logger"
	"fmt"
	"log"
	"os"
	"strconv"

	"golang.org/x/net/html"
)

const DEFAULT_DICTIONARY = "JMdict_e"

var unknownTags = map[string]int{}

func getEntryId(node *html.Node) int64 {
	i, err := strconv.ParseInt(node.FirstChild.Data, 10, 64)
	if err != nil {
		log.Fatal(err)
	}

	return i
}

func getKanji(node *html.Node) string {
	kanji := ""
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && c.Data == "keb" {
			if kanji != "" {
				log.Fatal(fmt.Errorf("additional kanji was found: %s", c.Data))
			}
			kanji = c.FirstChild.Data
		} else {
			val, ok := unknownTags[c.Data]
			if ok {
				unknownTags[c.Data] = val + 1
			} else {
				unknownTags[c.Data] = 1
			}
		}
	}
	return kanji
}

func getReading(node *html.Node) string {
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && c.Data == "reb" {
			return c.FirstChild.Data
		} else {
			val, ok := unknownTags[c.Data]
			if ok {
				unknownTags[c.Data] = val + 1
			} else {
				unknownTags[c.Data] = 1
			}
		}
	}

	return ""
}

func getSense(node *html.Node) Sense {
	sense := Sense{
		PartOfSpeech: []string{},
		XRef:         []string{},
		Glossary:     []string{},
		Field:        []string{},
		Misc:         []string{},
	}

	for c := node.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode {
			if c.Data == "pos" {
				sense.PartOfSpeech = append(sense.PartOfSpeech, c.FirstChild.Data)
			} else if c.Data == "xref" {
				sense.XRef = append(sense.XRef, c.FirstChild.Data)
			} else if c.Data == "gloss" {
				sense.Glossary = append(sense.Glossary, c.FirstChild.Data)
			} else if c.Data == "field" {
				sense.Field = append(sense.Field, c.FirstChild.Data)
			} else if c.Data == "misc" {
				sense.Misc = append(sense.Misc, c.FirstChild.Data)
			} else {
				val, ok := unknownTags[c.Data]
				if ok {
					unknownTags[c.Data] = val + 1
				} else {
					unknownTags[c.Data] = 1
				}
			}
		}
	}
	return sense
}

func getEntry(node *html.Node) Entry {
	entry := Entry{
		EntSeq:  -1,
		Kanji:   []string{},
		Reading: []string{},
		Sense:   []Sense{},
	}

	for c := node.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.TextNode {
			// ignore em[ty text nodes
		} else if c.Type == html.ElementNode {
			// get thing based pn data
			if c.Data == "ent_seq" {
				entry.EntSeq = getEntryId(c)
			} else if c.Data == "k_ele" {
				entry.Kanji = append(entry.Kanji, getKanji(c))
			} else if c.Data == "r_ele" {
				entry.Reading = append(entry.Reading, getReading(c))
			} else if c.Data == "sense" {
				entry.Sense = append(entry.Sense, getSense(c))
			}
		}
	}

	return entry
}

func processEntry(node *html.Node) (*Entry, error) {
	if node.Type != html.ElementNode || node.Data != "entry" {
		//logger.Error("Could not parse *html.Node type: %s, Data: %s", node.Type, node.Data)
		return nil, nil
	}
	entry := getEntry(node)
	return &entry, nil
}

func processDictionary(node *html.Node) ([]Entry, error) {
	entries := []Entry{}
	if node.Type == html.ElementNode && node.Data == "jmdict" {
		logger.Debug("Parsing jmdict")
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			entry, err := processEntry(c)
			if err != nil {
				return nil, err
			}
			if entry != nil {
				entries = append(entries, *entry)
			}
		}
	} else {
		logger.Debug("Parsing some other xml: %s", node.Data)
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			newEntries, err := processDictionary(c)
			if err != nil {
				return nil, err
			}
			entries = append(entries, newEntries...)
		}
	}
	return entries, nil
}

func ParseJMDict(filepath string) ([]Entry, error) {

	xmlFile, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer xmlFile.Close()

	doc, err := html.Parse(xmlFile)
	if err != nil {
		return nil, err
	}

	return processDictionary(doc)
}

func ToJson(entries []Entry, filepath string) error {
	jsonData, err := json.MarshalIndent(entries, "", " ")
	if err != nil {
		return err
	}

	// Write the json to a file
	err = os.WriteFile(filepath, jsonData, 0644)
	return err
}
