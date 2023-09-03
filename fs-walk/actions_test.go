package main

import (
	"os"
	"testing"
)

func TestFilterOut(t *testing.T) {
	testCases := []struct {
		name     string
		file     string
		ext      string
		minSize  int64
		expected bool
	}{
		{"FilterNoExtension", "testdata/dir.log", "", 0, false},
		{"FilterExtensionMatch", "testdata/dir.log", ".log", 0, false},
		{"FilterExtensionNoMatch", "testdata/dir.log", ".sh", 0, true},
		{"FilterExtensionSizeMatch", "testdata/dir.log", ".log", 10, false},
		{"FilterExtensionSizeNoMatch", "testdata/dir.log", ".log", 20, true},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			info, err := os.Stat(testCase.file)
			if err != nil {
				t.Fatal(err)
			}
			result := filterOut(testCase.file, config{
				ext:  testCase.ext,
				size: testCase.minSize,
			}, info)
			if result != testCase.expected {
				t.Fatalf("Expected %v, Received %v", testCase.expected, result)
			}
		})
	}
}
