package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
)

func run(inputFileName string, skipPreview bool) error {
	// Create an empty output html file
	outName := fmt.Sprintf("./tmp/%s.html", filepath.Base(inputFileName))
	fmt.Println("out: ", outName)
	outputFile, err := os.Create(outName)
	if err != nil {
		return err
	}
	defer func() { _ = outputFile.Close() }()

	// Read markdown data from input file and transform it into html body
	input, err := os.ReadFile(inputFileName)
	if err != nil {
		return err
	}
	body := parseContent(input)

	// Render the body to output file
	err = renderToFile(outputFile, template.HTML(body))
	if err != nil {
		return err
	}

	// Don't preview if skipPreview is true
	if skipPreview {
		return nil
	}

	return preview(outName)
}

func main() {
	fileName := flag.String("file", "", "Markdown file to preview")
	skipPreview := flag.Bool("s", false, "Skip auto-preview")
	flag.Parse()

	if *fileName == "" {
		flag.Usage()
		os.Exit(1)
	}

	if err := run(*fileName, *skipPreview); err != nil {
		log.Fatalln(err)
	}
}
