package cmd

import (
	"bytes"
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

const timeFormat = "Jan/02 @15:04"

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

func getOne(url string, id int) (item Item, err error) {
	fullURL := fmt.Sprintf("%s/todo/%d", url, id-1)
	items, err := getItems(fullURL)
	if err != nil {
		return item, err
	}
	if len(items) != 1 {
		return item, fmt.Errorf("invalid result: %w", ErrInvalidData)
	}
	return items[0], err
}

func sendRequest(method string, url string, contentType string, body io.Reader, expStatus int) error {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return err
	}
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	res, err := newClient().Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != expStatus {
		msg, err := io.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("unable to read response: %v", err)
		}
		err = ErrInvalidResponse
		if res.StatusCode == http.StatusNotFound {
			err = ErrNotFound
		}
		return fmt.Errorf("%w: %s, statusCode: %d", err, msg, res.StatusCode)
	}
	return nil
}

func addItem(apiRoot string, task string) error {
	url := fmt.Sprintf("%s/todo", apiRoot)
	item := struct {
		Task string `json:"task"`
	}{
		Task: task,
	}
	var body bytes.Buffer
	if err := json.NewEncoder(&body).Encode(item); err != nil {
		return err
	}
	return sendRequest(http.MethodPost, url, "application/json", &body, http.StatusCreated)
}

func completeItem(apiRoot string, id int) error {
	url := fmt.Sprintf("%s/todo/%d?complete", apiRoot, id-1)

	return sendRequest(http.MethodPatch, url, "", nil, http.StatusNoContent)
}

func deleteItem(apiRoot string, id int) error {
	url := fmt.Sprintf("%s/todo/%d", apiRoot, id-1)

	return sendRequest(http.MethodDelete, url, "", nil, http.StatusNoContent)
}
