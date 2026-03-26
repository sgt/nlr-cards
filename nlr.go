package nlr_cards

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-retryablehttp"
)

const (
	baseId               = 1000
	idMultFactor         = 10
	baseCardNumber       = 1
	cardNumberMultFactor = 3
)

var lastCardNumberRegexp = regexp.MustCompile(`<meta name="Description" content=".*\( \d+/(\d+) \)">`)

type NLR struct {
	OutputDir string

	basePngUrl  string
	baseHtmlUrl string

	client *http.Client
}

func NewNLR() NLR {
	retryClient := retryablehttp.NewClient()
	retryClient.Logger = nil
	retryClient.RetryWaitMax = 5 * time.Second
	retryClient.RetryMax = 2

	return NLR{
		OutputDir: "downloads",

		basePngUrl:  "https://nlr.ru/e-case3/sc2.php/web_gak/gc",
		baseHtmlUrl: "https://nlr.ru/e-case3/sc2.php/web_gak/lc",

		client: retryClient.StandardClient(),
	}
}

// Fetching

func (nlr *NLR) fetchPng(id, cardNumber int) ([]byte, error) {
	url := fmt.Sprintf("%s/%d/%d", nlr.basePngUrl, id, cardNumber)
	if resp, err := nlr.client.Get(url); err != nil {
		return nil, err
	} else {
		defer func() {
			if err := resp.Body.Close(); err != nil {
				panic(err)
			}
		}()
		return io.ReadAll(resp.Body)
	}
}

func (nlr *NLR) pngExists(id, cardNumber int) (bool, error) {
	if data, err := nlr.fetchPng(id, cardNumber); err != nil {
		return false, err
	} else {
		return len(data) != 0, nil
	}
}

func (nlr *NLR) fetchHtml(id int) (html []byte, err error) {
	url := fmt.Sprintf("%s/%d/1", nlr.baseHtmlUrl, id)
	if resp, err := nlr.client.Get(url); err != nil {
		return []byte{}, err
	} else {
		defer func() {
			err = resp.Body.Close()
		}()
		return io.ReadAll(resp.Body)
	}
}

func (nlr *NLR) htmlExists(id int) (bool, error) {
	if data, err := nlr.fetchHtml(id); err != nil {
		return false, err
	} else {
		html := string(data)
		return strings.HasPrefix(html, "\nError"), nil
	}
}

// generates save dir for id 1234 in downloads/1/2
func (nlr *NLR) saveDir(id int) string {
	thousands := id / 1000
	hundreds := id / 100 % 10
	return filepath.Join(nlr.OutputDir, strconv.Itoa(thousands), strconv.Itoa(hundreds))
}

func (nlr *NLR) save(id, cardNumber int, data []byte) error {
	dir := nlr.saveDir(id)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	filename := fmt.Sprintf("%d-%d.png", id, cardNumber)
	fullPath := filepath.Join(dir, filename)
	if err := os.WriteFile(fullPath, data, 0644); err != nil {
		return err
	}
	return nil
}

func (nlr *NLR) FetchAndSave(id, cardNumber int) error {
	data, err := nlr.fetchPng(id, cardNumber)
	if err != nil {
		return err
	}
	if len(data) == 0 {
		return errors.New("Card does not exist on the server.")
	}
	if err := nlr.save(id, cardNumber, data); err != nil {
		return err
	}
	return nil
}

func findPairForMaxIdSearch(baseId, multFactor int, existsFn func(int) (bool, error)) (left int, right int, err error) {
	left = 1
	right = baseId
	for {
		var ok bool
		ok, err = existsFn(right)
		if err != nil || !ok {
			return
		}
		left = right
		right *= multFactor
	}
}

func binarySearchLastId(left, right int, existsFn func(int) (bool, error)) (int, error) {
	nonZeroId := left
	zeroId := right

	for nonZeroId+1 < zeroId {
		mid := (nonZeroId + zeroId) / 2
		ok, err := existsFn(mid)
		if err != nil {
			return -1, err
		}

		if !ok {
			zeroId = mid
		} else {
			nonZeroId = mid
		}
	}
	return nonZeroId, nil
}

func (nlr *NLR) FindLastId() (int, error) {
	existsFn := func(id int) (bool, error) { return nlr.pngExists(id, 1) }

	left, right, err := findPairForMaxIdSearch(baseId, idMultFactor, existsFn)
	if err != nil {
		panic(err)
	}

	return binarySearchLastId(left, right, existsFn)
}

func (nlr *NLR) FindLastCardNumberInASillyWay(id int) (int, error) {
	existsFn := func(cardNumber int) (bool, error) { return nlr.pngExists(id, cardNumber) }

	left, right, err := findPairForMaxIdSearch(baseCardNumber, cardNumberMultFactor, existsFn)
	if err != nil {
		return -1, nil
	}

	return binarySearchLastId(left, right, existsFn)
}

func (nlr *NLR) FindLastCardNumberInASmartWay(id int) (int, error) {
	data, err := nlr.fetchHtml(id)
	if err != nil {
		return -1, err
	}

	matches := lastCardNumberRegexp.FindStringSubmatch(string(data))
	if len(matches) < 2 {
		return -1, errors.New("match not found")
	}
	return strconv.Atoi(matches[1])
}

func ReadCardsJsonFile(filename string) (map[int]int, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var cards map[int]int
	err = json.Unmarshal(data, &cards)
	if err != nil {
		return nil, err
	}

	return cards, nil
}

func WriteCardsJsonFile(filename string, cards map[int]int) error {
	jsonStr, err := json.Marshal(cards)
	if err != nil {
		return err
	}

	return os.WriteFile(filename, jsonStr, 0644)
}
