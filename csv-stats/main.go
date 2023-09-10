package main

import (
	"flag"
	"fmt"
	"io"
	"os"
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

	for _, fileName := range fileNames {
		fileReader, err := os.Open(fileName)

		if err != nil {
			return err
		}
		data, err := csv2Float(fileReader, column)
		if err != nil {
			return err
		}
		if err = fileReader.Close(); err != nil {
			return fmt.Errorf("failed to close file: %s", fileName)
		}
		allValues = append(allValues, data...)
	}

	_, err := fmt.Fprintln(out, operation(allValues))

	return err
}
