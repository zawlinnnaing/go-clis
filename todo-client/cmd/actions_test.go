package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"testing"
)

func TestListAc(t *testing.T) {
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
			replacer := strings.NewReplacer("\t", "", " ", "", "\n", "")
			if replacer.Replace(testCase.expOut) != replacer.Replace(out.String()) {
				t.Errorf("Expected out:\n %s, received:\n %s", testCase.expOut, out.String())
			}
		})
	}
}
