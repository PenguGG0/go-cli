package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
)

type statsFunc func(data []float64) float64

func sum(data []float64) float64 {
	res := 0.0

	for _, v := range data {
		res += v
	}

	return res
}

func avg(data []float64) float64 {
	return sum(data) / float64(len(data))
}

// Convert specific column of csv file into slice of float64
func csv2float(r io.Reader, column int) ([]float64, error) {
	cr := csv.NewReader(r)
	column--

	allData, err := cr.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("cannot read data from file: %w", err)
	}

	var columnData []float64
	var cellData float64
	for i, row := range allData {
		// Skip the title row
		if i == 0 {
			continue
		}

		// Valid the number of column
		if len(row) <= column {
			return nil, fmt.Errorf("%w: File has only %d columns", ErrInvalidColumn, len(row))
		}

		// Get data from the specific column in every row
		cellData, err = strconv.ParseFloat(row[column], 64)
		if err != nil {
			return nil, fmt.Errorf("%w: %s", ErrNotNumber, err)
		}

		columnData = append(columnData, cellData)
	}

	return columnData, nil
}
