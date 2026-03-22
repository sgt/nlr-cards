package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/alexflint/go-arg"
	"github.com/alitto/pond/v2"
	nlrcards "github.com/sgt/nlr-cards"
)

type dlCmd struct {
	MaxId          int    `default:"133781" help:"Max ID to download cards for"`
	PauseInterval  int    `default:"1000" help:"Pause after every such batch of files downloaded"`
	PauseDuration  int    `default:"3" help:"Pause duration in seconds"`
	MaxConcurrency int    `default:"20" help:"Max concurrent downloads"`
	JsonFile       string `default:"cards.json"`
}

type countCmd struct {
	JsonFile string `default:"cards.json"`
}

func main() {
	var args struct {
		Dl    *dlCmd    `arg:"subcommand:dl"`
		Count *countCmd `arg:"subcommand:count"`
	}
	p := arg.MustParse(&args)
	if p.Subcommand() == nil {
		p.Fail("missing subcommand")
	}

	switch {
	case args.Dl != nil:
		if err := downloadCards(args.Dl); err != nil {
			panic(err)
		}
	case args.Count != nil:
		if err := countCards(args.Count.JsonFile); err != nil {
			panic(err)
		}
	}
}

func downloadCards(options *dlCmd) error {
	nlr := nlrcards.NewNLR()

	cardCounts, err := nlrcards.ReadCardsJson(options.JsonFile)
	if err != nil {
		return fmt.Errorf("Failed to read '%s', error: %w\n", options.JsonFile, err)
	}

	pool := pond.NewPool(options.MaxConcurrency)

	for id, cardCount := range cardCounts {
		pool.Submit(func() {
		})
	}

	pool.StopAndWait()
	return nil
}

func countCards(jsonFilename string) error {
	nlr := nlrcards.NewNLR()

	log.Println("Finding max non-zero ID...")

	lastId, err := nlr.FindLastId()
	if err != nil {
		return err
	}

	log.Printf("Max non-zero ID: %d\n", lastId)

	log.Printf("Updating missing cards in %s...\n", jsonFilename)

	return nil
}
