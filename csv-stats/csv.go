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
	allData, err := csvReader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("cannot read data from file: %w", err)
	}
	// From natural number index to slice index
	column -= 1

	var data []float64

	for idx, row := range allData {
		if idx == 0 {
			continue
		}
		if column > len(row) {
			return nil, fmt.Errorf("%w: file only has %d columns", err, len(row))
		}
		val, err := strconv.ParseFloat(row[column], 64)
		if err != nil {
			return nil, fmt.Errorf("%w: %s", ErrNotNumber, row[column])
		}
		data = append(data, val)
	}

	return data, nil
}
