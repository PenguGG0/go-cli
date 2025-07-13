package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
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
	pipeline = append(pipeline, timeoutStep{
		name:    "git push",
		exe:     "git",
		message: "Git Push: SUCCESS",
		proj:    proj,
		args:    []string{"push", "origin", "master"},
		timeout: 5 * time.Second,
	})

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	errCh := make(chan error)
	done := make(chan struct{})

	go func() {
		for _, s := range pipeline {
			msg, err := s.execute()
			if err != nil {
				errCh <- err
				return
			}

			if _, err = fmt.Fprintln(out, msg); err != nil {
				errCh <- fmt.Errorf("can't print: %w", err)
				return
			}
		}
		close(done)
	}()

	for {
		select {
		case recSignal := <-sig:
			signal.Stop(sig)
			return fmt.Errorf("%s Exiting: %w", recSignal, ErrSignal)
		case err := <-errCh:
			return err
		case <-done:
			return nil
		}
	}
}

func main() {
	proj := flag.String("p", ".", "Project directory")
	flag.Parse()

	if err := run(*proj, os.Stdout); err != nil {
		log.Fatalln(err)
	}
}
