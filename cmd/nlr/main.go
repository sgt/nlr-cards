package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"

	"github.com/alexflint/go-arg"
	"github.com/alitto/pond/v2"
	nlrcards "github.com/sgt/nlr-cards"
)

type dlCmd struct {
	MaxId          int    `default:"133781" help:"Max ID to download cards for"`
	PauseInterval  int    `default:"1000" help:"Pause after every such batch of files downloaded"`
	PauseDuration  int    `default:"3" help:"Pause duration in seconds"`
	MaxConcurrency int    `default:"5" help:"Max concurrent downloads"`
	JsonFile       string `default:"cards.json"`
}

type countCmd struct {
	MaxId          int    `default:"133781" help:"Max ID to download cards for"`
	MaxConcurrency int    `default:"2" help:"Max concurrent downloads"`
	JsonFile       string `default:"cards.json"`
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
		if err := countCards(args.Count); err != nil {
			panic(err)
		}
	}
}

func downloadCards(options *dlCmd) error {
	nlr := nlrcards.NewNLR()

	cardCounts, err := nlrcards.ReadCardsJsonFile(options.JsonFile)
	if err != nil {
		return fmt.Errorf("Failed to read '%s', error: %w\n", options.JsonFile, err)
	}

	pool := pond.NewPool(options.MaxConcurrency)

	for id, cardCount := range cardCounts {
		for cardNumber := 1; cardNumber <= cardCount; cardNumber++ {
			pool.Submit(func() {
				if err := nlr.FetchAndSave(id, cardNumber); err != nil {
					log.Printf("Failed to download card: %d/%d\n", id, cardNumber)
				}
			})
		}
	}

	pool.StopAndWait()
	return nil
}

func countCards(options *countCmd) error {

	nlr := nlrcards.NewNLR()

	log.Printf("Updating missing cards in %s...\n", options.JsonFile)
	cards, err := nlrcards.ReadCardsJsonFile(options.JsonFile)
	if os.IsNotExist(err) {
		cards = make(map[int]int)
	} else if err != nil {
		return err
	}

	type resultPair = struct{ Id, LastCardNumber int }
	resultsChan := make(chan resultPair, 1000)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	go func() {
		<-ctx.Done()
		log.Println("Ctrl+C pressed, wait for the workers to finish...")
	}()

	pool := pond.NewPool(options.MaxConcurrency, pond.WithContext(ctx))

	var cardsUpdateWG sync.WaitGroup
	cardsUpdateWG.Add(1)

	go func() {
		defer cardsUpdateWG.Done()
		for r := range resultsChan {
			cards[r.Id] = r.LastCardNumber
		}
	}()

	for id := 1; id <= options.MaxId; id++ {
		if _, ok := cards[id]; ok {
			continue
		}

		pool.Submit(func() {
			lastCardNumber, err := nlr.FindLastCardNumberInASmartWay(id)
			if err != nil {
				log.Printf("Failed to determine last card number for id %d\n", id)
				return
			}
			log.Printf("Id %d has %d cards.\n", id, lastCardNumber)
			resultsChan <- resultPair{id, lastCardNumber}
		})
	}

	pool.StopAndWait()
	close(resultsChan)
	cardsUpdateWG.Wait()

	log.Println("Writing json...")
	return nlrcards.WriteCardsJsonFile(options.JsonFile, cards)
}
