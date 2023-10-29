package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

var (
	ErrConnection      = errors.New("connection error")
	ErrNotFound        = errors.New("not found")
	ErrInvalidResponse = errors.New("invalid server response")
	ErrInvalidData     = errors.New("invalid data")
	ErrNotNumber       = errors.New("not a number")
)

type Item struct {
	Task        string
	Done        bool
	CreatedAt   time.Time
	CompletedAt time.Time
}

type Response struct {
	Data         []Item `json:"data"`
	Date         int64  `json:"date"`
	TotalResults int    `json:"total_results"`
}

func newClient() *http.Client {
	client := http.Client{
		Timeout: 10 * time.Second,
	}
	return &client
}

func getAll(apiRoot string) ([]Item, error) {
	url := fmt.Sprintf("%s/todo", apiRoot)

	return getItems(url)
}

func getItems(url string) ([]Item, error) {
	res, err := newClient().Get(url)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrConnection, err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		msg, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, fmt.Errorf("cannot read body: %v", err)
		}
		err = ErrInvalidResponse
		if res.StatusCode == http.StatusNotFound {
			err = ErrNotFound
		}
		return nil, fmt.Errorf("%w: %s", err, msg)
	}
	var resp Response
	if err = json.NewDecoder(res.Body).Decode(&resp); err != nil {
		return nil, fmt.Errorf("failed to decode: %v", err)
	}
	if resp.TotalResults == 0 {
		return nil, fmt.Errorf("no results found. %w", ErrNotFound)
	}
	return resp.Data, nil
}
