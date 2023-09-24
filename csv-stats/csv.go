package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
)

type operation func(data []float64) float64

func sum(data []float64) float64 {
	sum := 0.0
	for _, val := range data {
		sum += val
	}
	return sum
}

func average(data []float64) float64 {
	if len(data) == 0 {
		return 0.0
	}
	sum := sum(data)
	return sum / float64(len(data))
}

func csv2Float(reader io.Reader, column int) ([]float64, error) {
	csvReader := csv.NewReader(reader)
	csvReader.ReuseRecord = true
	// From natural number index to slice index
	column -= 1

	var data []float64

	for i := 0; ; i++ {
		row, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("cannot read data from file: %w", err)
		}
		if i == 0 {
			continue
		}
		if column >= len(row) {
			return nil, fmt.Errorf("%w: %d", ErrInvalidColumn, len(row))
		}
		val, err := strconv.ParseFloat(row[column], 64)
		if err != nil {
			return nil, fmt.Errorf("%w: %s", ErrNotNumber, row[column])
		}
		data = append(data, val)
	}

	return data, nil
}
