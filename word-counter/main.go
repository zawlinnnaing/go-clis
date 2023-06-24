package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
)

func Count(reader io.Reader, countLines bool) int {
	scanner := bufio.NewScanner(reader)

	if countLines {
		scanner.Split(bufio.ScanLines)
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
	flag.Parse()
	fmt.Println(Count(os.Stdin, *countLines))
}
