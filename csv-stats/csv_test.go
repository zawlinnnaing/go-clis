package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"testing"
	"testing/iotest"
)

func TestOperations(t *testing.T) {
	data := [][]float64{
		{10, 20, 15, 30, 45, 50, 100, 30},
		{5.5, 8, 2.2, 9.75, 8.45, 3, 2.5, 10.25, 4.75, 6.1, 7.67, 12.287, 5.47},
		{-10, -20},
		{102, 37, 44, 57, 67, 129},
	}

	testCases := []struct {
		name string
		op   operation
		exp  []float64
	}{
		{"Sum", sum, []float64{300, 85.927, -30, 436}},
		{"Avg", average, []float64{37.5, 6.609769230769231, -15, 72.666666666666666}},
	}

	for _, testCase := range testCases {
		for idx, exp := range testCase.exp {
			name := fmt.Sprintf("%sData%d", testCase.name, idx)
			t.Run(name, func(t *testing.T) {
				result := testCase.op(data[idx])
				if result != exp {
					t.Errorf("Expected %g, received %g instead", exp, result)
				}
			})
		}
	}
}

func TestCSV2Float(t *testing.T) {
	csvData := `IP Address,Requests,Response Time
	192.168.0.199,2056,236
	192.168.0.88,899,220
	192.168.0.199,3054,226
	192.168.0.100,4133,218
	192.168.0.199,950,238`
	testCases := []struct {
		name   string
		col    int
		exp    []float64
		expErr error
		r      io.Reader
	}{
		{
			name:   "Column2",
			col:    2,
			expErr: nil,
			exp:    []float64{2056, 899, 3054, 4133, 950},
			r:      bytes.NewBufferString(csvData),
		},
		{
			name:   "Column3",
			col:    3,
			expErr: nil,
			exp:    []float64{236, 220, 226, 218, 238},
			r:      bytes.NewBufferString(csvData),
		},
		{
			name:   "FailedRead",
			col:    1,
			exp:    nil,
			expErr: iotest.ErrTimeout,
			r:      iotest.TimeoutReader(bytes.NewBuffer([]byte{0})),
		},
		{
			name:   "FailedNotNumber",
			col:    1,
			exp:    nil,
			expErr: ErrNotNumber,
			r:      bytes.NewBufferString(csvData),
		},
		{
			name:   "FailedInvalidColumn",
			col:    4,
			exp:    nil,
			expErr: ErrInvalidColumn,
			r:      bytes.NewBufferString(csvData),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			res, err := csv2Float(testCase.r, testCase.col)
			if testCase.expErr != nil {
				if err == nil {
					t.Errorf("Expected error: %d, got nil instead", testCase.expErr)
				}
				if !errors.Is(err, testCase.expErr) {
					t.Errorf("Expected error %q, received %q", testCase.expErr, err)
				}
				return
			}
			if err != nil {
				t.Errorf("Unexpected error %q:", err)
			}
			for i, exp := range testCase.exp {
				if res[i] != exp {
					t.Errorf("Expected %g, received %g", exp, res[i])
				}
			}
		})
	}
}
