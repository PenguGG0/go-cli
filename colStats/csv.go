package main

import (
	"encoding/csv"
	"errors"
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

func dataMin(data []float64) float64 {
	if len(data) == 0 {
		return 0
	}

	res := data[0]
	for _, v := range data[1:] {
		res = min(res, v)
	}

	return res
}

func dataMax(data []float64) float64 {
	if len(data) == 0 {
		return 0
	}

	res := data[0]
	for _, v := range data[1:] {
		res = max(res, v)
	}

	return res
}

// Convert specific column of csv file into slice of float64.
func csv2float(r io.Reader, column int) ([]float64, error) {
	cr := csv.NewReader(r)
	cr.ReuseRecord = true
	column--

	columnData := make([]float64, 0, 2500)
	var cellData float64
	for i := 0; ; i++ {
		// Read a record
		row, err := cr.Read()
		if err != nil {
			// break the loop when reach the end of file
			if errors.Is(err, io.EOF) {
				break
			} else {
				return nil, fmt.Errorf("cannot read data from file: %w", err)
			}
		}
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
			return nil, fmt.Errorf("%w: %w", ErrNotNumber, err)
		}

		columnData = append(columnData, cellData)
	}

	return columnData, nil
}
