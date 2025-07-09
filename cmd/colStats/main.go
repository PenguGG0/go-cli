package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"sync"
)

type consumerRes struct {
	data []float64
	err  error
}

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
	case "min":
		opFunc = dataMin
	case "max":
		opFunc = dataMax
	default:
		return fmt.Errorf("%w: %s", ErrInvalidOperation, op)
	}

	data := make([]float64, 0)

	fileCh := make(chan string)
	resCh := make(chan consumerRes)

	// Producer
	go func() {
		defer close(fileCh)
		for _, fileName := range fileNames {
			fileCh <- fileName
		}
	}()

	// Consumer
	wg := sync.WaitGroup{}
	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for fileName := range fileCh {
				var res consumerRes

				// Open file
				file, err := os.Open(fileName)
				if err != nil {
					res.err = fmt.Errorf("cannot open file %s: %w", fileName, err)
					resCh <- res
					continue
				}
				// Extract data of specific column from csv file
				res.data, err = csv2float(file, column)
				if err != nil {
					res.err = err
				}
				// Close file
				if err = file.Close(); err != nil {
					res.err = err
				}

				resCh <- res
			}
		}()
	}

	// Wait until all files have been processed
	go func() {
		wg.Wait()
		close(resCh)
	}()

	// Final consumer
	for res := range resCh {
		if res.err != nil {
			return res.err
		}
		data = append(data, res.data...)
	}

	// Perform operation on data
	if len(data) == 0 {
		return ErrNoData
	}
	if _, err := fmt.Fprintln(out, opFunc(data)); err != nil {
		return err
	}
	return nil
}

func main() {
	op := flag.String("op", "sum", "Operation to be executed\nValid options are: sum, avg, min, max")
	column := flag.Int("col", 1, "CSV column on which to execute operation")
	flag.Usage = func() {
		flag.PrintDefaults()
		_, _ = fmt.Fprintf(flag.CommandLine.Output(), "\ne.g., colStats -op [str] -col [int] [file1.csv file2.csv] \n")
	}
	flag.Parse()

	if err := run(flag.Args(), *op, *column, os.Stdout); err != nil {
		log.Fatalln(err)
	}
}
