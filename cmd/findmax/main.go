package main

import (
	"errors"
	"fmt"
	"log"

	nlr_cards "github.com/sgt/nlr-cards"
)

type searchPair struct {
	nonZeroId, zeroId int
}

func findSearchPair(baseId, factor int) (searchPair, error) {
	nonZeroId := 1
	current := baseId
	for {
		_, err := nlr_cards.Fetch(current, 1)
		if err != nil && !errors.Is(err, nlr_cards.ErrEmptyContent) {
			return searchPair{0, 0}, err
		}

		if errors.Is(err, nlr_cards.ErrEmptyContent) {
			log.Printf("Max id is somewhere between %d and %d\n", nonZeroId, current)
			return searchPair{nonZeroId: nonZeroId, zeroId: current}, nil
		}
		nonZeroId = current
		current *= factor
	}
}

func binarySearchNonZero(sp searchPair) (int, error) {
	nonZeroId := sp.nonZeroId
	zeroId := sp.zeroId

	for nonZeroId+1 < zeroId {
		mid := (nonZeroId + zeroId) / 2
		_, err := nlr_cards.Fetch(mid, 1)
		if err != nil && !errors.Is(err, nlr_cards.ErrEmptyContent) {
			return 0, err
		}

		if errors.Is(err, nlr_cards.ErrEmptyContent) {
			zeroId = mid
		} else {
			nonZeroId = mid
		}
	}
	return nonZeroId, nil
}

func main() {

	sp, err := findSearchPair(1000, 10)
	if err != nil {
		panic(err)
	}

	lastNonZero, err := binarySearchNonZero(sp)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Max non-zero ID: %d\n", lastNonZero)
}
