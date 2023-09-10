package main

import (
	"bytes"
	"errors"
	"os"
	"testing"
)

func TestRun(t *testing.T) {
	testCases := []struct {
		name   string
		col    int
		op     string
		exp    string
		files  []string
		expErr error
	}{
		{
			name:   "RunAvg1File",
			col:    3,
			op:     "avg",
			exp:    "227.6\n",
			expErr: nil,
			files:  []string{"./testdata/example.csv"},
		},
		{
			name:   "RunAvgMultipleFiles",
			col:    3,
			op:     "avg",
			exp:    "231.88235294117646\n",
			files:  []string{"./testdata/example.csv", "./testdata/example2.csv"},
			expErr: nil,
		},
		{
			name:   "RunFailRead",
			col:    2,
			op:     "avg",
			exp:    "",
			expErr: os.ErrNotExist,
			files:  []string{"./testdata/example.csv", "./testdata/fake.csv"},
		},
		{
			name:   "RunFailColumn",
			col:    0,
			op:     "avg",
			exp:    "",
			files:  []string{"./testdata/example.csv"},
			expErr: ErrInvalidColumn,
		},
		{
			name:   "RunFailNoFiles",
			col:    2,
			op:     "avg",
			exp:    "",
			files:  []string{},
			expErr: ErrNoFiles,
		},
		{
			name:   "RunFailOperation",
			col:    2,
			op:     "invalid",
			exp:    "",
			files:  []string{"./testdata/example.csv"},
			expErr: ErrInvalidOperation,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			var res bytes.Buffer
			err := run(testCase.files, testCase.op, testCase.col, &res)
			if testCase.expErr != nil {
				if !errors.Is(err, testCase.expErr) {
					t.Errorf("Expected error: %q, received %q", testCase.expErr, err)
				}
				return
			}
			if err != nil {
				t.Errorf("Unexpected error occurred: %q", err)
			}
			if res.String() != testCase.exp {
				t.Errorf("Expected %s, received %s", testCase.exp, res.String())
			}
		})
	}
}
