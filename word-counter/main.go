package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
)

func Count(reader io.Reader, countLines bool, countBytes bool) int {
	scanner := bufio.NewScanner(reader)

	if countLines {
		scanner.Split(bufio.ScanLines)
	} else if countBytes {
		scanner.Split(bufio.ScanBytes)
	} else {
		scanner.Split(bufio.ScanWords)
	}

	wordCount := 0

	for scanner.Scan() {
		wordCount += 1
	}

	return wordCount
}

func main() {
	countLines := flag.Bool("l", false, "Count lines")
	countBytes := flag.Bool("b", false, "Count bytes")
	flag.Parse()
	fmt.Println(Count(os.Stdin, *countLines, *countBytes))
}
