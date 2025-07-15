package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	countLines := flag.Bool("l", false, "Count lines")
	countBytes := flag.Bool("b", false, "Count bytes")
	fileName := flag.String("f", "", "File to count")
	flag.Parse()

	var r io.Reader
	var err error

	if *fileName == "" {
		r = os.Stdin
	} else {
		r, err = os.Open(*fileName)
		if err != nil {
			log.Fatalln(err)
		}
	}

	c, b, err := count(r, *countLines, *countBytes)
	if err != nil {
		log.Fatalln(err)
	}

	if *countLines {
		fmt.Println("lines: ", c)
	} else {
		fmt.Println("words: ", c)
	}

	if *countBytes {
		fmt.Println("bytes: ", b)
	}
}

func count(reader io.Reader, countLines, countBytes bool) (int, int, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return 0, 0, err
	}
	newReader := bytes.NewReader(data)

	b := 0
	c := 0

	scannerCount := bufio.NewScanner(newReader)
	// If the countLines flag is not set, we'll count words
	// Otherwise, we'll use the default line-delimited config
	if !countLines {
		scannerCount.Split(bufio.ScanWords)
	}

	for scannerCount.Scan() {
		c += 1
	}
	if err = scannerCount.Err(); err != nil {
		return c, b, err
	}

	// If the countBytes flag is set, we'll count bytes
	if countBytes {
		_, err = newReader.Seek(0, io.SeekStart)
		if err != nil {
			return 0, 0, err
		}

		scannerByte := bufio.NewScanner(newReader)
		scannerByte.Split(bufio.ScanBytes)
		for scannerByte.Scan() {
			b += 1
		}
		if err = scannerByte.Err(); err != nil {
			return c, b, err
		}
	}

	return c, b, nil
}
