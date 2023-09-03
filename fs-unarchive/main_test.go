package main

import (
	"bytes"
	"compress/gzip"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

func TestRun(t *testing.T) {
	testFiles := []string{
		"hello-world.txt.gzip",
		"dir1/test.txt.gzip",
		"dir1/dir2/test.txt.gzip",
	}
	tempDir, cleanup := createTempDir(t, testFiles)
	defer cleanup()
	targetDir, cleanupTarget := createTempDir(t, []string{})
	defer cleanupTarget()
	var outBuffer bytes.Buffer
	err := run(tempDir, targetDir, &outBuffer)
	if err != nil {
		t.Fatal(err)
	}
	unzippedFiles := []string{}

	err = filepath.Walk(targetDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			t.Fatal(err)
		}
		ext := filepath.Ext(path)
		if ext != ".txt" {
			return nil
		}
		relDir, err := filepath.Rel(targetDir, path)
		if err != nil {
			return err
		}
		unzippedFiles = append(unzippedFiles, relDir)
		return err
	})
	if err != nil {
		t.Fatal(err)
	}

	if err != nil {
		t.Fatal(err)
	}
	if len(unzippedFiles) != len(testFiles) {
		t.Errorf("Expected unzipped files to be %d, received %d", len(testFiles), len(unzippedFiles))
	}
	sort.Strings(testFiles)
	sort.Strings(unzippedFiles)
	testFileStr := strings.ReplaceAll(strings.Join(testFiles, "\n"), ".gzip", "")
	resultStr := strings.Join(unzippedFiles, "\n")
	if testFileStr != resultStr {
		t.Errorf("Expected %s, received %s", testFileStr, resultStr)
	}
}

func createTempDir(t *testing.T, files []string) (string, func()) {
	t.Helper()
	tempDir, err := os.MkdirTemp("", "fs-unarchive")
	if err != nil {
		t.Fatal(err)
	}

	for _, file := range files {
		filePath := filepath.Join(tempDir, file)
		err = os.MkdirAll(filepath.Dir(filePath), 0755)
		if err != nil {
			t.Fatal(err)
		}
		f, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0755)
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()
		gzipWriter := gzip.NewWriter(f)
		defer gzipWriter.Close()
		gzipWriter.Write([]byte("Hello world"))
	}

	return tempDir, func() {
		os.RemoveAll(tempDir)
	}
}
