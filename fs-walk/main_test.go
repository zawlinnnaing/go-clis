package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRun(t *testing.T) {
	testCases := []struct {
		name     string
		root     string
		cfg      config
		expected string
	}{
		{name: "NoFilter", root: "testdata", cfg: config{ext: "", size: 0, list: true}, expected: "testdata/dir.log\ntestdata/dir2/script.sh\n"},
		{name: "FilterExtensionMatch", root: "testdata", cfg: config{ext: ".log", size: 0, list: true}, expected: "testdata/dir.log\n"},
		{name: "FilterExtensionSizeMatch", root: "testdata", cfg: config{ext: ".log", size: 10, list: true}, expected: "testdata/dir.log\n"},
		{name: "FilterExtensionSizeNoMatch", root: "testdata", cfg: config{ext: ".log", size: 20, list: true}, expected: ""},
		{name: "FilterExtensionNoMatch", root: "testdata", cfg: config{ext: ".gz", size: 0, list: true}, expected: ""},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			var buffer bytes.Buffer

			if err := run(testCase.root, &buffer, testCase.cfg); err != nil {
				t.Fatal(err)
			}
			result := buffer.String()

			if result != testCase.expected {
				t.Errorf("Expected: %v, received: %v", testCase.expected, result)
			}
		})
	}
}

func TestRunDel(t *testing.T) {
	testCases := []struct {
		name        string
		cfg         config
		extNoDelete string
		nDelete     int
		nNoDelete   int
		expected    string
	}{
		{name: "DeleteExtensionNoMatch", cfg: config{ext: ".log", del: true}, extNoDelete: ".gz", nDelete: 0, nNoDelete: 10, expected: ""},
		{name: "DeleteExtensionMatch", cfg: config{ext: ".log", del: true}, extNoDelete: "", nDelete: 10, nNoDelete: 0, expected: ""},
		{name: "DeleteExtensionMixed", cfg: config{ext: ".log", del: true}, extNoDelete: ".gz", nDelete: 5, nNoDelete: 5, expected: ""},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			var (
				buffer        bytes.Buffer
				logFileBuffer bytes.Buffer
			)

			testCase.cfg.delWrite = &logFileBuffer

			tempDir, cleanup := createTempDir(t, map[string]int{
				testCase.cfg.ext:     testCase.nDelete,
				testCase.extNoDelete: testCase.nNoDelete,
			})
			defer cleanup()
			if err := run(tempDir, &buffer, testCase.cfg); err != nil {
				t.Fatal(err)
			}
			result := buffer.String()
			if result != testCase.expected {
				t.Errorf("Expected %v, received %v", testCase.expected, result)
			}
			filesLeft, err := ioutil.ReadDir(tempDir)
			if err != nil {
				t.Fatal(err)
			}
			if len(filesLeft) != testCase.nNoDelete {
				t.Errorf("Expected files to left: %d, actual left %d", testCase.nNoDelete, len(filesLeft))
			}

			expectedLogLines := testCase.nDelete + 1
			logLines := bytes.Split(logFileBuffer.Bytes(), []byte("\n"))
			if len(logLines) != expectedLogLines {
				t.Errorf("Expected %d log lines, received %d log lines.", expectedLogLines, len(logLines))
			}
		})
	}
}

func TestRunArchive(t *testing.T) {
	testCases := []struct {
		name         string
		cfg          config
		extNoArchive string
		nArchive     int
		nNoArchive   int
	}{
		{name: "ArchiveExtensionNoMatch", cfg: config{ext: ".log"}, extNoArchive: ".gz", nArchive: 0, nNoArchive: 10},
		{name: "ArchiveExtensionMatch", cfg: config{ext: ".log"}, extNoArchive: "", nArchive: 10, nNoArchive: 0},
		{name: "ArchiveExtensionMixed", cfg: config{ext: ".log"}, extNoArchive: ".gz", nArchive: 5, nNoArchive: 5},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			var outBuffer bytes.Buffer
			tempDir, cleanup := createTempDir(t, map[string]int{
				testCase.extNoArchive: testCase.nNoArchive,
				testCase.cfg.ext:      testCase.nArchive,
			})
			defer cleanup()
			archiveDir, cleanupArchive := createTempDir(t, nil)
			defer cleanupArchive()
			testCase.cfg.archive = archiveDir
			if err := run(tempDir, &outBuffer, testCase.cfg); err != nil {
				t.Fatal(err)
			}

			pattern := filepath.Join(tempDir, fmt.Sprintf("*%s", testCase.cfg.ext))
			expectedFiles, err := filepath.Glob(pattern)
			if err != nil {
				t.Fatal(err)
			}
			expectOutput := strings.Join(expectedFiles, "\n")
			result := strings.TrimSpace(outBuffer.String())

			if result != expectOutput {
				t.Errorf("Expected %v, received %v", expectOutput, result)
			}

			archivedFiles, err := ioutil.ReadDir(archiveDir)
			if err != nil {
				t.Fatal(err)
			}

			if len(archivedFiles) != testCase.nArchive {
				t.Errorf("Expected %d files to be archived, got %d", testCase.nArchive, len(archivedFiles))
			}

		})
	}
}

func createTempDir(t *testing.T, filesMap map[string]int) (dirname string, cleanup func()) {
	t.Helper()
	tempDir, err := ioutil.TempDir("", "fs-walk-test")
	if err != nil {
		t.Fatal(err)
	}
	for ext, num := range filesMap {
		for i := 0; i < num; i++ {
			fileName := fmt.Sprintf("file-%d%s", i+1, ext)
			filePath := filepath.Join(tempDir, fileName)
			if err := ioutil.WriteFile(filePath, []byte("Dummy"), 0644); err != nil {
				t.Fatal(err)
			}
		}
	}
	return tempDir, func() { os.RemoveAll(tempDir) }
}
