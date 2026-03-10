package nlr_cards

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type NLR struct {
	outputDir string
	baseUrl   string

	client *http.Client
}

var ErrEmptyContent = errors.New("content is empty")

func NewNLR() NLR {
	return NLR{
		outputDir: "downloads",
		baseUrl:   "https://nlr.ru/e-case3/sc2.php/web_gak/gc",

		client: &http.Client{Timeout: 30 * time.Second},
	}
}
func (nlr *NLR) Fetch(id, cardNumber int) ([]byte, error) {
	url := fmt.Sprintf("%s/%d/%d", nlr.baseUrl, id, cardNumber)
	resp, err := nlr.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if len(content) == 0 {
		return nil, ErrEmptyContent
	}
	return content, nil
}

// generates save dir for id 1234 in downloads/1/2
func (nlr *NLR) saveDir(id int) string {
	thousands := id / 1000
	hundreds := id / 100 % 10
	return filepath.Join(nlr.outputDir, strconv.Itoa(thousands), strconv.Itoa(hundreds))
}

// returns false if there is no such file on the server
func (nlr *NLR) save(id, cardNumber int, data []byte) (bool, error) {
	dir := nlr.saveDir(id)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return false, err
	}
	filename := fmt.Sprintf("%d-%d.png", id, cardNumber)
	fullPath := filepath.Join(dir, filename)
	os.WriteFile(fullPath, data, 0644)
	return true, nil
}

func (nlr *NLR) FetchAndSave(id, cardNumber int) (bool, error) {
	data, err := nlr.Fetch(id, cardNumber)
	if err != nil && !errors.Is(err, ErrEmptyContent) {
		return false, err
	}
	if errors.Is(err, ErrEmptyContent) {
		return false, nil
	}
	return nlr.save(id, cardNumber, data)
}
