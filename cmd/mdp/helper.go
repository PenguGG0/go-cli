package main

import (
	"bytes"
	"html/template"
	"os"

	"github.com/yuin/goldmark"
)

func parseContent(input []byte) []byte {
	output := new(bytes.Buffer)

	if err := goldmark.Convert(input, output); err != nil {
		return nil
	}

	return output.Bytes()
}

func renderToFile(outputFile *os.File, body template.HTML) error {
	// Get the html template
	templateCache, err := newTemplateCache()
	if err != nil {
		return err
	}
	tmpl, ok := templateCache["base.gohtml"]
	if !ok {
		return err
	}

	// render the template with the html body and write it to the output file
	err = tmpl.ExecuteTemplate(outputFile, "base", &templateData{Body: body})
	if err != nil {
		return err
	}

	return nil
}
