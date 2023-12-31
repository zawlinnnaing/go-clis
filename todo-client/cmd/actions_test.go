package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
)

func transformString(input string) string {
	replacer := strings.NewReplacer("\t", "", " ", "", "\n", "")
	return replacer.Replace(input)
}

func TestListAction(t *testing.T) {
	testCases := []struct {
		name   string
		expErr error
		expOut string
		resp   struct {
			Status int
			Body   string
		}
		closeServer bool
	}{
		{
			name:   "Results",
			expOut: "-  1  Task 1  \n-  2  Task 2  \n",
			resp:   testResp["resultsMany"],
		},
		{
			name:   "NoResult",
			expErr: ErrNotFound,
			resp:   testResp["noResult"],
		},
		{
			name:        "InvalidURL",
			expErr:      ErrConnection,
			resp:        testResp["noResult"],
			closeServer: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			url, cleanup := mockServer(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(testCase.resp.Status)
				fmt.Fprintln(w, testCase.resp.Body)
			})
			defer cleanup()
			if testCase.closeServer {
				cleanup()
			}
			var out bytes.Buffer
			err := listAction(&out, url)
			if testCase.expErr != nil {
				if err == nil {
					t.Errorf("Expected error: %v, received nil", err)
				}
				if !errors.Is(err, testCase.expErr) {
					t.Errorf("Expected error: %v, received: %v", testCase.expErr, err)
				}
				return
			}
			if err != nil {
				t.Errorf("Unexpected error: %q", err)
			}
			if transformString(testCase.expOut) != transformString(out.String()) {
				t.Errorf("Expected out:\n %s, received:\n %s", testCase.expOut, out.String())
			}
		})
	}
}

func TestViewAction(t *testing.T) {
	testCases := []struct {
		name   string
		expErr error
		expOut string
		id     string
		resp   struct {
			Status int
			Body   string
		}
	}{
		{
			name:   "ResultOne",
			expErr: nil,
			expOut: `Task: Task1
			CreatedAt: Oct/28 @08:23
			Completed: No
			`,
			resp: testResp["resultOne"],
			id:   "1",
		},
		{
			name:   "NotFound",
			expErr: ErrNotFound,
			resp:   testResp["notFound"],
			id:     "1",
		},
		{
			name:   "InvalidNumber",
			expErr: ErrNotNumber,
			resp:   testResp["noResult"],
			id:     "a",
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			apiServer, cleanup := mockServer(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(testCase.resp.Status)
				fmt.Fprintln(w, testCase.resp.Body)
			})
			defer cleanup()
			var out bytes.Buffer
			err := viewAction(apiServer, testCase.id, &out)
			if testCase.expErr != nil {
				if !errors.Is(err, testCase.expErr) {
					t.Errorf("Expected error: %v, received: %v", testCase.expErr, err)
				}
				return
			}
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if transformString(testCase.expOut) != transformString(out.String()) {
				t.Errorf("Expected out: %s, received: %s", testCase.expOut, out.String())
			}
		})
	}
}

func TestAddAction(t *testing.T) {
	expURLPath := "/todo"
	expMethod := http.MethodPost
	expBody := "{\"task\":\"Task 1\"}"
	expOut := fmt.Sprintf("Added task %q to the list.\n", "Task 1")
	expContent := "application/json"
	args := []string{"Task", "1"}

	url, cleanup := mockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != expURLPath {
			t.Errorf("Expected url path: %s, received: %s", expURLPath, r.URL.Path)
		}
		if r.Method != expMethod {
			t.Errorf("Expected method: %s, received: %s", expMethod, r.Method)
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("Unexpected err: %s", err)
		}
		if transformString(string(body)) != transformString(expBody) {
			t.Errorf("Expected body: %s, received: %s", expBody, string(body))
		}
		r.Body.Close()
		contentType := r.Header.Get("Content-Type")
		if contentType != expContent {
			t.Errorf("Expected content type: %s, received: %s", expContent, contentType)
		}
		w.WriteHeader(testResp["created"].Status)
		fmt.Fprintln(w, testResp["created"].Body)
	})
	defer cleanup()
	var out bytes.Buffer
	err := addAction(url, args, &out)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if out.String() != expOut {
		t.Errorf("Expected output: %s, received: %s", expOut, out.String())
	}
}
