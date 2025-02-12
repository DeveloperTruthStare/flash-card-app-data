package main

import (
	"flag"
	"flashcardmaker/src/database"
	"flashcardmaker/src/logger"
	"flashcardmaker/src/pokemon"
	"log"
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

	pokemon.LoadData()
}
