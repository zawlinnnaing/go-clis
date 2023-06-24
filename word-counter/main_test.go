package main

import (
	"bytes"
	"testing"
)

func TestCount(t *testing.T) {
	buffer := bytes.NewBufferString("hello world")
	words := Count(buffer, false, false)
	if words != 2 {
		t.Errorf("Expected %v, Received %v", 2, words)
	}
}

func TestCountLines(t *testing.T) {
	buffer := bytes.NewBufferString("hello world\n And this is another line\n With another line")
	lines := Count(buffer, true, false)
	if lines != 3 {
		t.Errorf("Expected 3 lines, Received %v", lines)
	}
}

func TestCountBytes(t *testing.T) {
	buffer := bytes.NewBufferString("hello world")
	bytes := Count(buffer, false, true)
	if bytes != 11 {
		t.Errorf("Expected 12 bytes, Received %v", bytes)
	}
}
