package main

import (
	"context"
	"log"
	"sync"
	"sync/atomic"
	"time"

	nlr_cards "github.com/sgt/nlr-cards"
)

const (
	maxId = 133781 // found by nlr-find-max

	// be courteous to the NLR server
	pauseInterval = 1000 // pause after every such batch of files downloaded
	pauseDuration = 3 * time.Second
	maxConcurrent = 20
)

var (
	// RWMutex acts as the pause gate
	// readers (downloaders) can pass freely until a writer (pauser) locks it
	gate       sync.RWMutex
	downloaded int64
)

func main() {

	nlr := nlr_cards.NewNLR()

	// semaphore channel to limit concurrency
	sem := make(chan struct{}, maxConcurrent)
	var wg sync.WaitGroup

	// context for graceful shutdown (optional but recommended)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for id := 1; id <= maxId; id++ {
		// acquire semaphore slot
		sem <- struct{}{}
		wg.Add(1)

		go func(id int) {
			defer wg.Done()
			defer func() { <-sem }() // release semaphore slot

			select {
			case <-ctx.Done():
				return
			default:
			}

			// try downloading cards until no card fetched
			cardNumber := 1
			for {
				// acquire read lock (blocks if pause is active)
				gate.RLock()
				ok, err := nlr.FetchAndSave(id, cardNumber)
				if err != nil {
					panic(err)
				}
				gate.RUnlock()

				count := atomic.AddInt64(&downloaded, 1)
				if count%pauseInterval == 0 {
					triggerPause()
				}

				if !ok {
					break
				}
				cardNumber += 1
			}
		}(id)
	}

	wg.Wait()
	log.Printf("Successfully downloaded cards for %d ids.\n", atomic.LoadInt64(&downloaded))
}

func triggerPause() {
	log.Printf("Downloaded %d files, pausing for %v... ", downloaded, pauseDuration)
	gate.Lock()
	time.Sleep(pauseDuration)
	gate.Unlock()
	log.Println("Resumed.")
}
