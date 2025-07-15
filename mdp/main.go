package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"os"
)

func run(inputFileName string) (string, error) {
	tempFile, err := os.CreateTemp("", "mdp*.html")
	if err != nil {
		return "", err
	}
	defer func() {
		if err = tempFile.Close(); err != nil {
			log.Fatalln(err)
		}
	}()
	outName := tempFile.Name()
	fmt.Println("out: ", outName)

	// Read markdown data from input file and transform it into html body
	input, err := os.ReadFile(inputFileName)
	if err != nil {
		return "", err
	}
	body := parseContent(input)

	// Render the body to output file
	err = renderToFile(tempFile, template.HTML(body))
	if err != nil {
		return "", err
	}

	return outName, nil
}

func main() {
	fileName := flag.String("file", "", "Markdown file to preview")
	skipPreview := flag.Bool("s", false, "Skip auto-preview")
	flag.Parse()

	if *fileName == "" {
		flag.Usage()
		os.Exit(1)
	}

	outName, err := run(*fileName)
	if err != nil {
		log.Fatalln(err)
	}

	// Skip preview if skipPreview is true
	if !*skipPreview {
		err = preview(outName)
		if err != nil {
			return
		}
	}
}
