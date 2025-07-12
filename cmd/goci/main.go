package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

type executer interface {
	execute() (string, error)
}

func run(proj string, out io.Writer) error {
	pipeline := make([]executer, 0)

	pipeline = append(pipeline, step{
		name:    "go build",
		exe:     "go",
		message: "Go Build: SUCCESS",
		proj:    proj,
		args:    []string{"build", ".", "errors"},
	})

	pipeline = append(pipeline, step{
		name:    "go test",
		exe:     "go",
		message: "Go Test: SUCCESS",
		proj:    proj,
		args:    []string{"test", "-v"},
	})
	pipeline = append(pipeline, exceptionStep{
		name:    "go fmt",
		exe:     "gofmt",
		message: "Gofmt: SUCCESS",
		proj:    proj,
		args:    []string{"-l", "."},
	})

	for _, s := range pipeline {
		msg, err := s.execute()
		if err != nil {
			return err
		}

		if _, err = fmt.Fprintln(out, msg); err != nil {
			return fmt.Errorf("can't print: %w", err)
		}
	}

	return nil
}

func main() {
	proj := flag.String("p", ".", "Project directory")
	flag.Parse()

	if err := run(*proj, os.Stdout); err != nil {
		log.Fatalln(err)
	}
}
