package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

func run(fileNames []string, op string, column int, out io.Writer) error {
	var opFunc statsFunc

	// Validate file slice is not empty
	// Validate column number is larger than 1
	if len(fileNames) == 0 {
		return ErrNoFiles
	}
	if column < 1 {
		return fmt.Errorf("%w: %d", ErrInvalidColumn, column)
	}

	switch op {
	case "sum":
		opFunc = sum
	case "avg":
		opFunc = avg
	default:
		return fmt.Errorf("%w: %s", ErrInvalidOperation, op)
	}

	data := make([]float64, 0)
	for _, fileName := range fileNames {
		// Open file
		file, err := os.Open(fileName)
		if err != nil {
			return fmt.Errorf("cannot open file: %w", err)
		}

		// Extract data of specific column from csv file
		columnData, err := csv2float(file, column)
		if err != nil {
			return err
		}

		if err = file.Close(); err != nil {
			return err
		}

		data = append(data, columnData...)
	}

	if _, err := fmt.Fprintln(out, opFunc(data)); err != nil {
		return err
	}

	return nil
}

func main() {
	op := flag.String("op", "sum", "Operation to be executed")
	column := flag.Int("col", 1, "CSV column on which to execute operation")
	flag.Usage = func() {
		flag.PrintDefaults()
		_, _ = fmt.Fprintf(flag.CommandLine.Output(), "\ne.g., colStats -op [sum] -col [1] [file1.csv file2.csv] \n")
	}
	flag.Parse()

	if err := run(flag.Args(), *op, *column, os.Stdout); err != nil {
		log.Fatalln(err)
	}
}
