package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/zawlinnnaing/go-clis/to-do/todo"
)

func setUpAPI(t *testing.T) (string, func()) {
	t.Helper()

	tempFile, err := os.CreateTemp("", "todo-test-")
	if err != nil {
		t.Fatal(err)
	}
	todoServer := httptest.NewServer(newMux(tempFile.Name()))

	for i := 1; i <= 3; i++ {
		var body bytes.Buffer
		taskName := fmt.Sprintf("Task number %d", i)
		item := struct {
			TaskName string `json:"task"`
		}{
			TaskName: taskName,
		}
		if err := json.NewEncoder(&body).Encode(&item); err != nil {
			t.Fatal(err)
		}
		resp, err := http.Post(todoServer.URL+"/todo", "application/json", &body)
		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != http.StatusCreated {
			t.Fatalf("Failed to seed items: Status %d, item number %d", resp.StatusCode, i)
		}
	}

	return todoServer.URL, func() {
		todoServer.Close()
		os.Remove(tempFile.Name())
	}
}

func TestGet(t *testing.T) {
	testCases := []struct {
		name       string
		path       string
		expCode    int
		expItems   int
		expContent string
	}{
		{
			name:       "Root",
			path:       "/",
			expCode:    http.StatusOK,
			expContent: "Hello from root route",
		},
		{
			name:    "NotFound",
			path:    "/todo/500",
			expCode: http.StatusNotFound,
		},
		{
			name:       "GetAll",
			path:       "/todo",
			expItems:   3,
			expCode:    http.StatusOK,
			expContent: "Task number 1",
		}, {
			name:       "GetItem",
			path:       "/todo/0",
			expItems:   1,
			expCode:    http.StatusOK,
			expContent: "Task number 1",
		},
	}

	testUrl, cleanup := setUpAPI(t)
	defer cleanup()
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			var (
				r struct {
					Data         todo.TaskList `json:"data"`
					Date         int64         `json:"date"`
					TotalResults int           `json:"total_results"`
				}
				body []byte
				err  error
			)
			resp, err := http.Get(testUrl + testCase.path)
			if err != nil {
				t.Error(err)
			}
			defer resp.Body.Close()
			if resp.StatusCode != testCase.expCode {
				t.Errorf("Expected code: %d, received %d", testCase.expCode, resp.StatusCode)
			}
			switch {
			case strings.Contains(resp.Header.Get("Content-Type"), "application/json"):
				if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
					t.Error(err)
				}
				if r.TotalResults != testCase.expItems {
					t.Errorf("Expected items: %d, received items: %d", testCase.expItems, r.TotalResults)
				}
				if r.Data[0].Task != testCase.expContent {
					t.Errorf("Expected content: %s, received content: %s", testCase.expContent, r.Data[0].Task)
				}
			case strings.Contains(resp.Header.Get("Content-Type"), "text/plain"):
				if body, err = io.ReadAll(resp.Body); err != nil {
					t.Error(err)
				}
				if !strings.Contains(string(body), testCase.expContent) {
					t.Errorf("Expected content: %s, received: %s", testCase.expContent, string(body))
				}
			default:
				t.Errorf("Unexpected content type: %s", resp.Header.Get("Content-Type"))
			}
		})
	}
}

func TestAdd(t *testing.T) {
	url, cleanup := setUpAPI(t)
	defer cleanup()
	taskToAdd := "Task number 4"
	t.Run("Add", func(t *testing.T) {
		var body bytes.Buffer
		item := struct {
			TaskName string `json:"task"`
		}{
			TaskName: taskToAdd,
		}
		if err := json.NewEncoder(&body).Encode(&item); err != nil {
			t.Fatal(err)
		}
		resp, err := http.Post(url+"/todo", "application/json", &body)
		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != http.StatusCreated {
			t.Errorf("Expected status code: %d, received status code: %d", http.StatusCreated, resp.StatusCode)
		}
	})
	t.Run("CheckAdd", func(t *testing.T) {
		resp, err := http.Get(url + "/todo/3")
		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status code: %d, received: %d", http.StatusOK, resp.StatusCode)
		}
		var item TodoResponse
		if err := json.NewDecoder(resp.Body).Decode(&item); err != nil {
			t.Fatal(err)
		}
		resp.Body.Close()
		if item.Data[0].Task != taskToAdd {
			t.Errorf("Expected task name: %s, received: %s", taskToAdd, item.Data[0].Task)
		}
	})
}

func TestDelete(t *testing.T) {
	url, cleanup := setUpAPI(t)
	defer cleanup()
	t.Run("Delete", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodDelete, url+"/todo/0", nil)
		if err != nil {
			t.Fatal(err)
		}
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		if res.StatusCode != http.StatusNoContent {
			t.Errorf("Expected status code: %d, received: %d", http.StatusNoContent, res.StatusCode)
		}
	})
	t.Run("CheckDelete", func(t *testing.T) {
		res, err := http.Get(url + "/todo/0")
		if err != nil {
			t.Fatal(err)
		}
		if res.StatusCode != http.StatusOK {
			t.Errorf("Expected status code: %d, received: %d", http.StatusOK, res.StatusCode)
		}
		var item TodoResponse
		if err := json.NewDecoder(res.Body).Decode(&item); err != nil {
			t.Fatal(err)
		}
		res.Body.Close()
		expTask := "Task number 2"
		if item.Data[0].Task != expTask {
			t.Errorf("Expected task: %s, received: %s", expTask, item.Data[0].Task)
		}
	})
}

func TestComplete(t *testing.T) {
	url, cleanup := setUpAPI(t)
	defer cleanup()
	t.Run("Complete", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPatch, url+"/todo/1?complete=true", nil)
		if err != nil {
			t.Fatal(err)
		}
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		if res.StatusCode != http.StatusNoContent {
			t.Errorf("Expected status code: %d, received: %d", http.StatusNoContent, res.StatusCode)
		}
	})

	t.Run("CheckComplete", func(t *testing.T) {
		res, err := http.Get(url + "/todo/1")
		if err != nil {
			t.Fatal(err)
		}
		var item TodoResponse
		if err := json.NewDecoder(res.Body).Decode(&item); err != nil {
			t.Fatal(err)
		}
		if !item.Data[0].Done {
			t.Errorf("Expect task to be completed")
		}
	})
}

func TestMain(m *testing.M) {
	log.SetOutput(io.Discard)
	os.Exit(m.Run())
}
