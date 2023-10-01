package main

import (
	"bytes"
	"errors"
	"testing"
)

func TestRun(t *testing.T) {
	testCases := []struct {
		proj   string
		output string
		expErr error
		name   string
	}{
		{proj: "./testdata/tool", name: "success", expErr: nil, output: "Go build: success\nGo test: success\nGo format: success\nGit push: success\n"},
		{proj: "./testdata/toolErr", name: "fail", expErr: &StepError{step: "go build"}, output: ""},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			out := bytes.Buffer{}
			err := run(testCase.proj, &out)
			if testCase.expErr != nil {
				if err == nil {
					t.Errorf("Expected error; received none")
					return
				}
				if !errors.Is(err, testCase.expErr) {
					t.Errorf("Expected error %v; received %v", testCase.expErr, err)
					return
				}
				return
			}
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			if out.String() != testCase.output {
				t.Errorf("Expected out: %s, received %s", testCase.output, out.String())
			}
		})
	}
}
