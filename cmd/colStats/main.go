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

	fileCh := make(chan string)
	resCh := make(chan []float64)
	errCh := make(chan error)
	doneCh := make(chan bool)

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
				// Open file
				file, err := os.Open(fileName)
				if err != nil {
					errCh <- fmt.Errorf("cannot open file: %w", err)
					return
				}
				// Extract data of specific column from csv file
				columnData, err := csv2float(file, column)
				if err != nil {
					errCh <- err
				}
				// Close file
				if err = file.Close(); err != nil {
					errCh <- err
				}

				resCh <- columnData
			}
		}()
	}

	// Wait until all files have been processed
	go func() {
		wg.Wait()
		doneCh <- true
	}()

	// Final Consumer
	for {
		select {
		case err := <-errCh:
			return err
		case columnData := <-resCh:
			data = append(data, columnData...)
		case <-doneCh:
			if _, err := fmt.Fprintln(out, opFunc(data)); err != nil {
				return err
			}
			return nil
		}
	}
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
