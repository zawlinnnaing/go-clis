package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func setUpAPI(t *testing.T) (string, func()) {
	t.Helper()

	testServer := httptest.NewServer(newMux(""))

	return testServer.URL, func() {
		testServer.Close()
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
			path:    "/todo/not-found",
			expCode: http.StatusNotFound,
		},
	}

	testUrl, cleanup := setUpAPI(t)
	defer cleanup()
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			var (
				body []byte
				err  error
			)
			resp, err := http.Get(testUrl + testCase.path)
			if err != nil {
				t.Error(err)
			}
			defer resp.Body.Close()
			if resp.StatusCode != testCase.expCode {
				t.Errorf("Expected code: %d, received %d", testCase.expCode, resp.ProtoMinor)
			}
			switch {
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
