package main

import (
	"io/ioutil"
	"os"
	"regexp"
	"testing"
)

const (
	inputFile  = "./testdata/test1.md"
	resultFile = "test1.md.html"
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
	result := parseContent(input)
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
	if err := run(inputFile); err != nil {
		t.Fatal(err)
	}
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
