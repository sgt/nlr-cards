package nlr_cards

import (
	"errors"
	"fmt"
	"io"
	"net/http"
)

const baseURL = "https://nlr.ru/e-case3/sc2.php/web_gak/gc"

var ErrEmptyContent = errors.New("content is empty")

func Fetch(id, cardNumber int) ([]byte, error) {
	url := fmt.Sprintf("%s/%d/%d", baseURL, id, cardNumber)
	resp, err := http.Get(url)
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
