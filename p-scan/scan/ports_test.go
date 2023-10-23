package scan_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/zawlinnnaing/go-clis/p-scan/scan"
)

func TestParsePorts(t *testing.T) {
	testCases := []struct {
		name      string
		input     string
		expect    []int
		expectErr error
	}{
		{
			name:   "portsSlice",
			input:  "22,80,443",
			expect: []int{22, 80, 443},
		},
		{
			name:   "portsRange",
			input:  "22-25",
			expect: []int{22, 23, 24, 25},
		},
		{
			name:      "invalidPorts",
			input:     "invalid",
			expect:    nil,
			expectErr: scan.ErrInvalidPortsFormat,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			out, err := scan.ParsePorts(testCase.input)
			if testCase.expectErr != nil {
				if err == nil {
					t.Errorf("Expected error, received none")
				}
				if !errors.Is(err, testCase.expectErr) {
					t.Errorf("Expected error: %s, received error: %s", testCase.expectErr, err)
				}
				return
			}
			if !reflect.DeepEqual(out, testCase.expect) {
				t.Errorf("Expected ports: %v, received ports: %v", out, testCase.expect)
			}
		})
	}
}
