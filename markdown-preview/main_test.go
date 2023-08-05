package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"testing"
)

const (
	inputFile  = "./testdata/test1.md"
	goldenFile = "./testdata/test1.md.html"
)

func normalizeHTML(data []byte) string {
	dataStr := string(data)
	regex := regexp.MustCompile(" |\n")
	return regex.ReplaceAllLiteralString(dataStr, "")
}

func TestParseContent(t *testing.T) {
	input, err := ioutil.ReadFile(inputFile)
	if err != nil {
		t.Fatal(err)
	}
	result, err := parseContent(input, "")
	if err != nil {
		t.Fatal(err)
	}
	expected, err := ioutil.ReadFile(goldenFile)
	if err != nil {
		t.Fatal(err)
	}
	if normalizeHTML(result) != normalizeHTML(expected) {
		t.Logf("golden:\n%s\n", normalizeHTML(expected))
		t.Logf("result:\n%s\n", normalizeHTML(result))
		t.Error("Result content does not match golden file")
	}
}

func TestRun(t *testing.T) {
	var mockStdOut bytes.Buffer
	if err := run(inputFile, &mockStdOut, true, ""); err != nil {
		t.Fatal(err)
	}

	resultFile := strings.TrimSpace(mockStdOut.String())

	result, err := ioutil.ReadFile(resultFile)
	if err != nil {
		t.Fatal(err)
	}
	expected, err := ioutil.ReadFile(goldenFile)
	if err != nil {
		t.Fatal(err)
	}

	if normalizeHTML(result) != normalizeHTML(expected) {
		t.Logf("golden:\n%s\n", normalizeHTML(expected))
		t.Logf("result:\n%s\n", normalizeHTML(result))
		t.Error("Result content does not match golden file")
	}

	os.Remove(resultFile)
}
