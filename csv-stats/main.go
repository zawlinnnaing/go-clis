package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"runtime"
	"sync"
)

var AVAILABLE_OPERATIONS = []string{
	"sum",
	"avg",
}

func main() {
	op := flag.String("op", "sum", "Operation to be executed. Available operations: sum, avg")
	column := flag.Int("col", 1, "CSV column on which operation will be executed. (Starting from 1)")

	flag.Parse()

	if err := run(flag.Args(), *op, *column, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(fileNames []string, op string, column int, out io.Writer) error {
	var operation operation

	if len(fileNames) == 0 {
		return ErrNoFiles
	}

	if column < 1 {
		return fmt.Errorf("%w: %d", ErrInvalidColumn, column)
	}

	switch op {
	case "sum":
		operation = sum
	case "avg":
		operation = average
	default:
		return fmt.Errorf("%w: %s", ErrInvalidOperation, op)
	}

	allValues := []float64{}
	resCh := make(chan []float64)
	errCh := make(chan error)
	doneCh := make(chan struct{})
	filesCh := make(chan string)
	go func() {
		defer close(filesCh)
		for _, fileName := range fileNames {
			filesCh <- fileName
		}
	}()
	wg := sync.WaitGroup{}
	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for fileName := range filesCh {
				fileReader, err := os.Open(fileName)
				if err != nil {
					errCh <- fmt.Errorf("cannot open file: %w", err)
					continue
				}
				data, err := csv2Float(fileReader, column)
				if err != nil {
					errCh <- err
				}
				if err = fileReader.Close(); err != nil {
					errCh <- fmt.Errorf("failed to close file: %s", fileName)
				}
				resCh <- data
			}
		}()
	}

	go func() {
		wg.Wait()
		close(doneCh)
	}()

	for {
		select {
		case err := <-errCh:
			return err
		case data := <-resCh:
			allValues = append(allValues, data...)
		case <-doneCh:
			_, err := fmt.Fprintln(out, operation(allValues))
			return err
		}
	}
}
