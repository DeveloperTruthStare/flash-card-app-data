package main

import (
	"bufio"
	"flag"
	"flashcardmaker/src/database"
	"flashcardmaker/src/database/tables/vocab"
	"flashcardmaker/src/logger"
	"flashcardmaker/src/pokemon"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	var reinitialize bool
	flag.BoolVar(&reinitialize, "reinit", false, "reinitialize the database on startup (defaults to false if not specified)")

	// Parse the flags
	flag.Parse()

	logger.SetLogLevel(logger.INFO | logger.WARN | logger.ERROR)

	err := database.Initialize(reinitialize)
	if err != nil {
		log.Fatal(err)
	}

	data, err := pokemon.LoadData()
	if err != nil {
		log.Fatal(err)
	}

	for _, pokemon := range data {
		fmt.Println(pokemon.Japanese)
	}

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Enter word to search for:")

		if !scanner.Scan() {
			fmt.Println("Error reading input. Exiting")
			break
		}

		userInput := strings.TrimSpace(scanner.Text())

		if strings.ToLower(userInput) == "exit" {
			fmt.Println("Exiting")
			break
		}

		entries, err := vocab.SearchVocab(userInput)
		if err != nil {
			log.Fatal(err)
		}

		for _, entry := range entries {
			fmt.Println(entry.EntSeq)
		}
	}
}
