package main

import (
	"bytes"
	"errors"
	"html/template"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

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

	// Render the template with the html body and write it to the output file
	err = tmpl.ExecuteTemplate(outputFile, "base", &templateData{
		Title: "Markdown Preview Tool",
		Body:  body,
	})
	if err != nil {
		return err
	}

	return nil
}

func preview(fileName string) error {
	var cName string
	var cParams []string

	switch runtime.GOOS {
	case "linux":
		cName = "xdg-open"
	case "windows":
		cName = "cmd"
		cParams = []string{"/C", "start", "chrome"}
	case "darwin":
		cName = "open"
	default:
		return errors.New("OS not supported")
	}

	fileName, err := filepath.Abs(fileName)
	if err != nil {
		return err
	}

	cParams = append(cParams, fileName)
	cPath, err := exec.LookPath(cName)
	if err != nil {
		return err
	}

	return exec.Command(cPath, cParams...).Run()
}
