package main

import (
	"errors"
	"testing"
)

func TestFileParser(t *testing.T) {
	testCases := []struct {
		file     string
		name     string
		expErr   error
		numSteps int
	}{
		{
			name:     "Valid file",
			file:     "./testdata/valid.yaml",
			expErr:   nil,
			numSteps: 4,
		},
		{
			name:   "Invalid file extension",
			file:   "./testdata/invalid.txt",
			expErr: ErrInvalidFile,
		},
		{
			name:   "Invalid step in file",
			file:   "./testdata/invalid-step.yaml",
			expErr: ErrInvalidStep,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			steps, err := parseFile(testCase.file)
			if testCase.expErr != nil {
				if err == nil {
					t.Errorf("Expected error %v. Received none.", testCase.expErr)
					return
				}
				if !errors.Is(err, testCase.expErr) {
					t.Errorf("Expected error %v, Received %v", testCase.expErr, err)
					return
				}
				return
			}
			if err != nil {
				t.Errorf("Unexpected error %v", err)
				return
			}
			if len(steps) != testCase.numSteps {
				t.Errorf("Expected steps to %d, received: %d", testCase.numSteps, len(steps))
			}
		})
	}

}
